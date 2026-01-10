package ai_analyze_core

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/fansqz/fancode-backend/common/ai_provider"
	"github.com/fansqz/fancode-backend/common/config"
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/constants"
	"github.com/fansqz/fancode-backend/models/dto"
)

// VisualDescriptionAnalyzer AI分析代码中数据结构的分析器
// 先触发StartAnalyzeCode，然后通过GetVisualDescription获取结果
// 只有startAnalyzeCode返回成功以后，GetVisualDescription才能返回结果，否则被阻塞
type VisualDescriptionAnalyzer struct {
	aiProvider        ai_provider.AIProvider
	visualDescription *dto.VisualDescription
	mutex             sync.Mutex
	cond              *sync.Cond
	isAnalyzing       bool
	analysisError     error
}

// NewVisualDescriptionAnalyzer 创建新的可视化描述分析器
func NewVisualDescriptionAnalyzer(aiConfig *config.AIConfig) *VisualDescriptionAnalyzer {
	analyzer := &VisualDescriptionAnalyzer{
		aiProvider: ai_provider.NewAIProvider(aiConfig),
	}
	analyzer.cond = sync.NewCond(&analyzer.mutex)
	return analyzer
}

// GetVisualDescription 获取用户代码的分析结果
// 如果解析未完成，会阻塞等待直到解析完成
func (a *VisualDescriptionAnalyzer) GetVisualDescription(ctx context.Context) (*dto.VisualDescription, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// 等待解析完成
	for a.isAnalyzing {
		a.cond.Wait()
	}

	// 检查是否有解析错误
	if a.analysisError != nil {
		return nil, a.analysisError
	}

	if a.visualDescription == nil {
		return nil, fmt.Errorf("visual description not found")
	}
	return a.visualDescription, nil
}

// StartAnalyzeCode 开始分析用户代码,异步解析，解析完成以后存储结果
func (a *VisualDescriptionAnalyzer) StartAnalyzeCode(ctx context.Context, code string, language constants.LanguageType) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// 如果已经在解析中，返回错误
	if a.isAnalyzing {
		return fmt.Errorf("analysis is already in progress")
	}

	// 重置状态
	a.isAnalyzing = true
	a.analysisError = nil
	a.visualDescription = nil

	go func() {
		visualDescription, err := a.analyzeCode(ctx, code, language)

		a.mutex.Lock()
		defer a.mutex.Unlock()

		if err != nil {
			a.analysisError = err
			logger.WithCtx(ctx).Errorf("[VisualDescriptionAnalyzer] Failed to analyze code: %v", err)
		} else {
			a.visualDescription = visualDescription
		}

		// 标记解析完成并通知等待的goroutine
		a.isAnalyzing = false
		a.cond.Broadcast()
	}()
	return nil
}

// IsAnalyzing 检查是否正在解析中
func (a *VisualDescriptionAnalyzer) IsAnalyzing() bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.isAnalyzing
}

// GetAnalysisStatus 获取解析状态
func (a *VisualDescriptionAnalyzer) GetAnalysisStatus() (bool, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.isAnalyzing, a.analysisError
}

// AnalyzeCode 分析用户代码中的数据结构
func (a *VisualDescriptionAnalyzer) analyzeCode(ctx context.Context, code string, language constants.LanguageType) (*dto.VisualDescription, error) {
	logger.WithCtx(ctx).Infof("[VisualDescriptionAnalyzer] Analyzing code for language: %s", language)

	// 构建AI提示词
	prompt := a.buildAnalysisPrompt(code, language)

	// 调用AI服务
	response, err := a.aiProvider.Chat(ctx, prompt)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[VisualDescriptionAnalyzer] AI analysis failed: %v", err)
		return nil, fmt.Errorf("AI analysis failed: %w", err)
	}

	// 解析AI响应
	visualDescription, err := a.parseAIResponse(response)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[VisualDescriptionAnalyzer] Failed to parse AI response: %v", err)
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	logger.WithCtx(ctx).Infof("[VisualDescriptionAnalyzer] Analysis completed, found type: %s", visualDescription.VisualType)
	return visualDescription, nil
}

// buildAnalysisPrompt 构建AI分析提示词
func (a *VisualDescriptionAnalyzer) buildAnalysisPrompt(code string, language constants.LanguageType) string {
	return fmt.Sprintf(`请分析以下代码中的数据结构，并通过如下步骤处理代码分析任务：
1. 首先分析提供的代码，明确其核心功能和意图
2. 识别代码中使用的主要数据结构
3. 根据识别结果，返回指定格式的 JSON 数据

以下为%s代码：
"%s"

请识别代码中的数据结构类型，可能包括：
1. 数组 (array) - 一维数组
2. 链表 (linkList) - 单向或双向链表
3. 二叉树 (binaryTree) - 二叉树结构
4. 图 (graph) - 图结构

请以JSON格式返回结果，格式如下：
{
  "visualType": "array|linkList|binaryTree|graph",
  "description": {
    // 根据类型返回相应的描述对象
  }
}

对于二维数组类型，返回：
{
  "visualType": "array2d", // array2d表示二维数组，不要改动
  "description": {
    "arrayName": "数组变量名",
    "rowPointNames": ["用于数组行取值的变量名1", "用于数组行取值的变量名2"],  // 举例：array[i][j]，这里填i
    "colPointNames": ["用于数组列取值的变量名1", "用于数组列取值的变量名2"],  // 举例：array[i][j]，这里填j
  }
}

对于一维数组 or 字符串类型，返回：
{
  "visualType": "array", // 如果是数组则填数组名称，如果是字符串则填字符串名称
  "description": {
    "arrayName": "数组变量名",
    "pointNames": ["用于数组取值的变量名1", "用于数组取值的变量名2"], // 所有用于目标数组取值的变量都需要配置，比如'数组名称[i]'，你需要识别出i，我指所有取值的都需要识别。哪些未用于取值的，但是代表了数组下标的也需要进行识别
    "displayType": "array" // 数组展示的类型，可选为array（普通数组展示）和array-bar（将数组展示为柱状图），array-bar用于排序算法，或者其他需要展示高度的算法中，其余请使用array
  }
}

对于链表类型，返回：
{
  "visualType": "linkList", 
  "description": {
    "linkNode": "链表节点结构体名",
    "data": "数据域属性名",
    "next": "next指针属性名",
    "prev": "prev指针属性名" // 双向链表才有
  }
}

对于二叉树类型，返回：
{
  "visualType": "binaryTree",
  "description": {
    "treeNode": "树节点结构体名",
    "data": "数据域属性名", 
    "left": "左子树指针属性名",
    "right": "右子树指针属性名"
  }
}

对于图类型，返回：
{
  "visualType": "graph",
  "description": {
    "graphNode": "图节点结构体名",
    "data": "数据域属性名",
    "nexts": ["邻接节点属性名1", "邻接节点属性名2"]
  }
}

如果代码中没有明显的数据结构，返回：
{
  "visualType": "array",
  "description": {
    "arrayName": "array",
    "pointNames": [],
    "displayType": "array"
  }
}

请确保分析准确，识别结果应基于代码中明确出现的结构，避免主观推测。`, language, code)
}

// parseAIResponse 解析AI响应
func (a *VisualDescriptionAnalyzer) parseAIResponse(response string) (*dto.VisualDescription, error) {
	// 尝试从响应中提取JSON
	jsonStr := a.extractJSONFromResponse(response)

	var result struct {
		VisualType  string      `json:"visualType"`
		Description interface{} `json:"description"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal AI response: %w", err)
	}

	// 验证可视化类型
	visualType := constants.VisualType(result.VisualType)
	if !a.isValidVisualType(visualType) {
		// 如果类型无效，默认为数组
		visualType = constants.ArrayType
	}

	// 根据类型构建具体的描述对象
	description, err := a.buildDescriptionObject(visualType, result.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to build description object: %w", err)
	}

	return &dto.VisualDescription{
		VisualType:  visualType,
		Description: description,
	}, nil
}

// extractJSONFromResponse 从AI响应中提取JSON字符串
func (a *VisualDescriptionAnalyzer) extractJSONFromResponse(response string) string {
	// 尝试找到JSON开始和结束的位置
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")

	if start == -1 || end == -1 || start >= end {
		// 如果没有找到JSON，返回原始响应
		return response
	}

	return response[start : end+1]
}

// isValidVisualType 验证可视化类型是否有效
func (a *VisualDescriptionAnalyzer) isValidVisualType(visualType constants.VisualType) bool {
	for _, validType := range constants.VisualTypeList {
		if visualType == validType {
			return true
		}
	}
	return false
}

// buildDescriptionObject 根据类型构建描述对象
func (a *VisualDescriptionAnalyzer) buildDescriptionObject(visualType constants.VisualType, rawDescription interface{}) (interface{}, error) {
	descriptionBytes, err := json.Marshal(rawDescription)
	if err != nil {
		return nil, err
	}

	switch visualType {
	case constants.ArrayType:
		var desc dto.ArrayDescription
		if err := json.Unmarshal(descriptionBytes, &desc); err != nil {
			return nil, err
		}
		return desc, nil

	case constants.Array2DType:
		var desc dto.Array2DDescription
		if err := json.Unmarshal(descriptionBytes, &desc); err != nil {
			return nil, err
		}
		return desc, nil

	case constants.LinkListType:
		var desc dto.LinkListDescription
		if err := json.Unmarshal(descriptionBytes, &desc); err != nil {
			return nil, err
		}
		return desc, nil

	case constants.BinaryTreeType:
		var desc dto.BinaryTreeDescription
		if err := json.Unmarshal(descriptionBytes, &desc); err != nil {
			return nil, err
		}
		return desc, nil

	case constants.GraphType:
		var desc dto.GraphDescription
		if err := json.Unmarshal(descriptionBytes, &desc); err != nil {
			return nil, err
		}
		return desc, nil

	default:
		// 默认返回数组描述
		return dto.ArrayDescription{
			ArrayName:   "array",
			PointNames:  []string{},
			DisplayType: "array",
		}, nil
	}
}

// AnalyzeCodeWithFallback 带降级策略的代码分析
func (a *VisualDescriptionAnalyzer) AnalyzeCodeWithFallback(ctx context.Context, code string, language constants.LanguageType) (*dto.VisualDescription, error) {
	// 首先尝试AI分析
	visualDescription, err := a.analyzeCode(ctx, code, language)
	if err == nil {
		return visualDescription, nil
	}

	logger.WithCtx(ctx).Warnf("[VisualDescriptionAnalyzer] AI analysis failed, falling back to rule-based analysis: %v", err)

	// AI分析失败时，使用基于规则的分析作为降级策略
	return a.ruleBasedAnalysis(code, language), nil
}

// ruleBasedAnalysis 基于规则的代码分析（降级策略）
func (a *VisualDescriptionAnalyzer) ruleBasedAnalysis(code string, language constants.LanguageType) *dto.VisualDescription {
	// 简单的规则匹配
	if a.containsLinkedListPattern(code) {
		return &dto.VisualDescription{
			VisualType: constants.LinkListType,
			Description: dto.LinkListDescription{
				LinkNode: "Node",
				Data:     "data",
				Next:     "next",
			},
		}
	}

	if a.containsBinaryTreePattern(code) {
		return &dto.VisualDescription{
			VisualType: constants.BinaryTreeType,
			Description: dto.BinaryTreeDescription{
				TreeNode: "TreeNode",
				Data:     "val",
				Left:     "left",
				Right:    "right",
			},
		}
	}

	if a.containsGraphPattern(code) {
		return &dto.VisualDescription{
			VisualType: constants.GraphType,
			Description: dto.GraphDescription{
				GraphNode: "Node",
				Data:      "val",
				Nexts:     []string{"neighbors"},
			},
		}
	}

	// 默认返回数组
	return &dto.VisualDescription{
		VisualType: constants.ArrayType,
		Description: dto.ArrayDescription{
			ArrayName:   "array",
			PointNames:  []string{},
			DisplayType: "array",
		},
	}
}

// containsLinkedListPattern 检查是否包含链表模式
func (a *VisualDescriptionAnalyzer) containsLinkedListPattern(code string) bool {
	patterns := []string{
		`type.*struct.*\{.*\w+\s+\w+.*\w+\s+\*.*\}`,
		`\.next\s*=`,
		`\.prev\s*=`,
		`ListNode`,
		`Node.*next`,
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, code); matched {
			return true
		}
	}
	return false
}

// containsBinaryTreePattern 检查是否包含二叉树模式
func (a *VisualDescriptionAnalyzer) containsBinaryTreePattern(code string) bool {
	patterns := []string{
		`type.*struct.*\{.*\w+\s+\w+.*\w+\s+\*.*\w+\s+\*.*\}`,
		`\.left\s*=`,
		`\.right\s*=`,
		`TreeNode`,
		`left.*right`,
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, code); matched {
			return true
		}
	}
	return false
}

// containsGraphPattern 检查是否包含图模式
func (a *VisualDescriptionAnalyzer) containsGraphPattern(code string) bool {
	patterns := []string{
		`neighbors`,
		`adjacency`,
		`graph`,
		`\[.*\]\*.*\w+`,
		`map.*\[\]`,
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, code); matched {
			return true
		}
	}
	return false
}
