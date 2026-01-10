package gdb_debugger

import (
	"context"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/cpp"
)

// VariableInfo 存储变量信息
type VariableInfo struct {
	Name     string
	Type     string
	Location struct {
		Line   int
		Column int
	}
}

// FunctionInfo 存储函数信息及其包含的变量
type FunctionInfo struct {
	Name     string
	Location struct {
		Line   int
		Column int
	}
	Variables []VariableInfo
}

// extractVariableName 从声明中提取变量名
func extractVariableName(node *sitter.Node, sourceCode []byte) string {
	if node == nil {
		return ""
	}

	// 获取完整的声明文本
	declText := node.Content(sourceCode)

	// 处理数组声明
	if strings.Contains(declText, "[") {
		// 找到第一个 [ 的位置
		idx := strings.Index(declText, "[")
		if idx > 0 {
			// 提取 [ 之前的部分作为变量名
			return strings.TrimSpace(declText[:idx])
		}
	}

	// 处理指针声明
	if strings.Contains(declText, "*") {
		// 找到最后一个 * 的位置
		idx := strings.LastIndex(declText, "*")
		if idx >= 0 && idx < len(declText)-1 {
			// 提取 * 之后的部分作为变量名
			return strings.TrimSpace(declText[idx+1:])
		}
	}

	// 处理引用声明
	if strings.Contains(declText, "&") {
		// 找到最后一个 & 的位置
		idx := strings.LastIndex(declText, "&")
		if idx >= 0 && idx < len(declText)-1 {
			// 提取 & 之后的部分作为变量名
			return strings.TrimSpace(declText[idx+1:])
		}
	}

	// 如果没有特殊字符，直接返回清理后的文本
	return strings.TrimSpace(declText)
}

// createVariableInfo 创建变量信息
func createVariableInfo(declarator *sitter.Node, typeNode *sitter.Node, sourceCode []byte) *VariableInfo {
	if declarator == nil {
		return nil
	}

	varName := extractVariableName(declarator, []byte(sourceCode))
	if varName == "" {
		return nil
	}

	typeStr := ""
	if typeNode != nil {
		typeStr = typeNode.Content([]byte(sourceCode))
	}

	return &VariableInfo{
		Name: varName,
		Type: typeStr,
		Location: struct {
			Line   int
			Column int
		}{
			Line:   int(declarator.StartPoint().Row + 1),
			Column: int(declarator.StartPoint().Column + 1),
		},
	}
}

// collectVariables 收集变量声明
func collectVariables(node *sitter.Node, sourceCode []byte, existingVars map[string]bool) []VariableInfo {
	var variables []VariableInfo

	// 处理变量声明
	if node.Type() == "declaration" {
		declarator := node.ChildByFieldName("declarator")
		if declarator != nil {
			identifier := declarator.ChildByFieldName("declarator")
			if identifier != nil {
				if varInfo := createVariableInfo(identifier, node.ChildByFieldName("type"), sourceCode); varInfo != nil {
					if !existingVars[varInfo.Name] {
						variables = append(variables, *varInfo)
						existingVars[varInfo.Name] = true
					}
				}
			}
		}
	}

	// 处理初始化声明
	if node.Type() == "init_declarator" {
		declarator := node.ChildByFieldName("declarator")
		if declarator != nil {
			parent := node.Parent()
			if parent != nil {
				if varInfo := createVariableInfo(declarator, parent.ChildByFieldName("type"), sourceCode); varInfo != nil {
					if !existingVars[varInfo.Name] {
						variables = append(variables, *varInfo)
						existingVars[varInfo.Name] = true
					}
				}
			}
		}
	}

	// 处理简单的变量声明（如 Item localItem; Value localValue;）
	if node.Type() == "declaration" {
		// 遍历所有子节点寻找变量声明
		cursor := sitter.NewTreeCursor(node)
		defer cursor.Close()

		if cursor.GoToFirstChild() {
			for {
				childNode := cursor.CurrentNode()

				// 检查是否是标识符（变量名）
				if childNode.Type() == "identifier" {
					varName := childNode.Content(sourceCode)
					if varName != "" && !existingVars[varName] {
						// 获取类型信息
						typeNode := node.ChildByFieldName("type")
						typeStr := ""
						if typeNode != nil {
							typeStr = typeNode.Content(sourceCode)
						}

						varInfo := &VariableInfo{
							Name: varName,
							Type: typeStr,
							Location: struct {
								Line   int
								Column int
							}{
								Line:   int(childNode.StartPoint().Row + 1),
								Column: int(childNode.StartPoint().Column + 1),
							},
						}
						variables = append(variables, *varInfo)
						existingVars[varName] = true
					}
				}

				if !cursor.GoToNextSibling() {
					break
				}
			}
		}
	}

	return variables
}

// collectParameters 收集函数参数
func collectParameters(parameters *sitter.Node, sourceCode []byte, existingVars map[string]bool) []VariableInfo {
	var params []VariableInfo
	if parameters == nil {
		return params
	}

	paramCursor := sitter.NewTreeCursor(parameters)
	defer paramCursor.Close()

	for paramCursor.GoToFirstChild(); paramCursor.GoToNextSibling(); {
		paramNode := paramCursor.CurrentNode()
		if paramNode.Type() == "parameter_declaration" {
			paramDeclarator := paramNode.ChildByFieldName("declarator")
			if paramDeclarator != nil {
				if varInfo := createVariableInfo(paramDeclarator, paramNode.ChildByFieldName("type"), sourceCode); varInfo != nil {
					if !existingVars[varInfo.Name] {
						params = append(params, *varInfo)
						existingVars[varInfo.Name] = true
					}
				}
			}
		}
	}

	return params
}

// ParseSourceFile 解析C/C++文件并返回函数及其包含的变量信息
func ParseSourceFile(sourceCode string) ([]FunctionInfo, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(cpp.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, []byte(sourceCode))
	if err != nil {
		return nil, fmt.Errorf("解析失败: %v", err)
	}
	defer tree.Close()

	var functions []FunctionInfo
	var currentFunction *FunctionInfo

	cursor := sitter.NewTreeCursor(tree.RootNode())
	defer cursor.Close()

	var traverse func(*sitter.TreeCursor, map[string]bool)
	traverse = func(cursor *sitter.TreeCursor, existingVars map[string]bool) {
		node := cursor.CurrentNode()

		// 检查函数定义
		if node.Type() == "function_definition" {
			declarator := node.ChildByFieldName("declarator")
			if declarator != nil {
				identifier := declarator.ChildByFieldName("declarator")
				if identifier != nil {
					funcInfo := FunctionInfo{
						Name: identifier.Content([]byte(sourceCode)),
						Location: struct {
							Line   int
							Column int
						}{
							Line:   int(identifier.StartPoint().Row + 1),
							Column: int(identifier.StartPoint().Column + 1),
						},
					}
					functions = append(functions, funcInfo)
					currentFunction = &functions[len(functions)-1]

					// 创建新的变量名映射用于去重
					funcVars := make(map[string]bool)

					// 收集函数参数
					parameters := declarator.ChildByFieldName("parameters")
					currentFunction.Variables = append(currentFunction.Variables, collectParameters(parameters, []byte(sourceCode), funcVars)...)

					// 检查函数体内的变量声明
					if cursor.GoToFirstChild() {
						traverse(cursor, funcVars)
						for cursor.GoToNextSibling() {
							traverse(cursor, funcVars)
						}
						cursor.GoToParent()
					}
				}
			}
		} else if currentFunction != nil {
			// 收集当前节点的变量
			variables := collectVariables(node, []byte(sourceCode), existingVars)
			currentFunction.Variables = append(currentFunction.Variables, variables...)

			// 递归遍历子节点
			if cursor.GoToFirstChild() {
				traverse(cursor, existingVars)
				for cursor.GoToNextSibling() {
					traverse(cursor, existingVars)
				}
				cursor.GoToParent()
			}
		} else {
			// 递归遍历子节点
			if cursor.GoToFirstChild() {
				traverse(cursor, existingVars)
				for cursor.GoToNextSibling() {
					traverse(cursor, existingVars)
				}
				cursor.GoToParent()
			}
		}
	}

	// 创建全局变量名映射
	globalVars := make(map[string]bool)
	traverse(cursor, globalVars)
	return functions, nil
}
