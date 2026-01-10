package debug_core

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	mapset "github.com/deckarep/golang-set"
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/constants"
	"github.com/fansqz/fancode-backend/utils"
	"github.com/sirupsen/logrus"
)

func (d *debugger) StructVisual(ctx context.Context, query *StructVisualQuery) (*StructVisualData, error) {
	logger.WithCtx(ctx).Infof("[GoDebugger] StructVisual")
	valueQuerySet := utils.List2set(query.Values)
	pointQuerySet := utils.List2set(query.Points)

	// 返回的 指针和可视化结构体
	pointVariables := make([]*VisualVariable, 0, 10)
	// VisualNodes 可视化结构的节点列表
	visualNodeSet := make(map[string]*StructVisualNode)
	visualNodes := make([]*StructVisualNode, 0, len(visualNodeSet))

	// 读取所有栈帧
	stackTrace, err := d.GetStackTrace(ctx)
	if err != nil {
		return nil, err
	}

	// 收集指针
	variables, err := d.GetFrameVariables(ctx, stackTrace[0].ID)
	if err != nil {
		return nil, err
	}
	for _, variable := range variables {
		// 如果变量的类型是结构体的指针，那么收集可视化变量以及指针列表
		if d.isTargetStruct(variable.Type, query.Struct) {
			// 获取指向目标结构体的指针
			visualVariable := NewVisualVariable(variable)
			pointVariables = append(pointVariables, visualVariable)
		}
	}

	// 读取目标结构体
	var allVariables []*Variable
	if allVariables, err = d.getAllFrameVariables(ctx); err != nil {
		logrus.Errorf("[StructVisual] getAllFrameVariables fail, err = %v", err)
		return nil, err
	}
	variables = make([]*Variable, 0, len(allVariables))
	for _, v := range allVariables {
		if d.isTargetStruct(v.Type, query.Struct) {
			variables = append(variables, v)
		}
	}

	// 广度遍历所有目标节点
	for len(variables) != 0 {
		newVariables := make([]*Variable, 0, len(variables))
		for _, variable := range variables {
			visualNode := &StructVisualNode{
				ID:   variable.Value,
				Type: variable.Type,
			}

			if variable.Reference == 0 || visualNode.ID == "" {
				continue
			}
			if _, ok := visualNodeSet[visualNode.ID]; ok {
				continue
			}

			// 读取结构体
			vars, err := d.GetVariables(ctx, variable.Reference)
			if err != nil {
				return nil, err
			}
			// 指针指向的结构体，查询出数据域和指针域
			for _, v := range vars {
				// 处理数据域
				if valueQuerySet.Contains(v.Name) {
					visualNode.Values = append(visualNode.Values, NewVisualVariable(v))
				}
				if pointQuerySet.Contains(v.Name) {
					// 处理指针域
					visualVariable := NewVisualVariable(v)
					visualVariable.Value = v.Value
					visualNode.Points = append(visualNode.Points, visualVariable)
					newVariables = append(newVariables, v)
				}
			}
			visualNodeSet[visualNode.ID] = visualNode
			visualNodes = append(visualNodes, visualNode)

		}
		variables = newVariables
	}

	return &StructVisualData{
		Points: pointVariables,
		Nodes:  visualNodes,
	}, nil
}

// 校验是否是目标结构体
func (d *debugger) isTargetStruct(structType string, targetType string) bool {
	var regexpStr string
	switch d.option.Language {
	case constants.LanguageGo:
		regexpStr = fmt.Sprintf(`^\*.*\.%s.*$`, targetType)
	case constants.LanguageC:
		regexpStr = fmt.Sprintf(`\b(%s)\b`, targetType)
	}
	re := regexp.MustCompile(regexpStr)
	return re.MatchString(structType)
}

func (d *debugger) ArrayVisual(ctx context.Context, query *ArrayVisualQuery) (*ArrayVisualData, error) {
	logger.WithCtx(ctx).Infof("[GoDebugger] VariableVisual")
	stackTrace, err := d.GetStackTrace(ctx)
	if err != nil {
		return nil, err
	}

	variables, err := d.GetFrameVariables(ctx, stackTrace[0].ID)
	if err != nil {
		return nil, err
	}
	pointQuerySet := mapset.NewSet()
	for _, point := range query.PointNames {
		pointQuerySet.Add(point)
	}

	// 返回的 指针和可视化结构体
	pointVariables := make([]*VisualVariable, 0, 10)
	// 目标结构体变量
	var structVariable *Variable

	for _, variable := range variables {
		if pointQuerySet.Contains(variable.Name) {
			visualVariable := NewVisualVariable(variable)
			pointVariables = append(pointVariables, visualVariable)
		}
		if query.ArrayName == variable.Name {
			structVariable = variable
		}
	}

	// c语言做特殊处理
	if d.option.Language == constants.LanguageC {
		structVariable, err = d.getStructVariableForC(ctx, stackTrace, structVariable)
		if err != nil {
			logger.WithCtx(ctx).Errorf("[VariableVisual] getStructVariableForC fail, err=%v", err)
			return nil, err
		}
	}

	// 读取结构体里面的内容
	var structs []*VisualVariable
	switch d.option.Language {
	case constants.LanguageGo:
		structs = d.getArrayNodesForGo(ctx, structVariable)
	case constants.LanguageC:
		structs = d.getArrayNodesForC(ctx, structVariable)
	default:
		structs = d.getArrayNodesNormal(ctx, structVariable)
	}

	return &ArrayVisualData{
		Points: pointVariables,
		Array:  structs,
	}, nil
}

// getStringVisualNodesForGo 将收集到的目标遍历转换为可视化节点 Go语言字符串要特殊处理
func (d *debugger) getArrayNodesForGo(ctx context.Context, structVariable *Variable) []*VisualVariable {
	if structVariable == nil {
		return []*VisualVariable{}
	}
	if structVariable.Type != "string" {
		return d.getArrayNodesNormal(ctx, structVariable)
	} else {
		trimmedStr := strings.Trim(structVariable.Value, "\"")
		vars := []*VisualVariable{}
		index := 0
		for i := 0; i < len(trimmedStr); i++ {
			char := trimmedStr[i]
			value := fmt.Sprintf("'%c'", char)
			// 特殊处理转义字符
			if char == '\\' && i < len(trimmedStr)-1 {
				value = fmt.Sprintf("'%c%c'", trimmedStr[i], trimmedStr[i+1])
				i++
			}
			vars = append(vars, &VisualVariable{
				Name:  fmt.Sprintf("[%d]", index),
				Type:  "byte",
				Value: value,
			})
			index++
		}
		return vars
	}
}

// getStringVisualNodesForC 将收集到的目标遍历转换为可视化节点 C语言字符串要特殊处理
func (d *debugger) getArrayNodesForC(ctx context.Context, structVariable *Variable) []*VisualVariable {
	if structVariable == nil {
		return []*VisualVariable{}
	}
	structs := d.getArrayNodesNormal(ctx, structVariable)
	// 将的char类型值进行处理，去除前面的数字
	re := regexp.MustCompile(`^\d+\s*`)
	for _, v := range structs {
		if v.Type == "char" {
			v.Value = re.ReplaceAllString(v.Value, "")
		}
	}
	return structs
}

func (d *debugger) getArrayNodesNormal(ctx context.Context, structVariable *Variable) []*VisualVariable {
	if structVariable == nil {
		return []*VisualVariable{}
	}
	vars, err := d.GetVariables(ctx, structVariable.Reference)
	if err != nil {
		return []*VisualVariable{}
	}
	newVars := make([]*VisualVariable, len(vars))
	for i, va := range vars {
		newVars[i] = NewVisualVariable(va)
	}
	return newVars
}

// getStructVariableForC c语言获取到的是结构体的指针，需要通过栈中获取到该数组定义的时候的结构名称
func (d *debugger) getStructVariableForC(ctx context.Context, stackTrace []*StackFrame, structVariable *Variable) (*Variable, error) {
	if structVariable == nil {
		return nil, nil
	}
	// 收集所有变量
	variables := []*Variable{}
	for _, s := range stackTrace[1:] {
		vs, err := d.GetFrameVariables(ctx, s.ID)
		if err != nil {
			return nil, err
		}
		variables = append(variables, vs...)
	}
	for _, v := range variables {
		if structVariable.Value == v.Value && v.Reference != 0 {
			return v, nil
		}
	}
	return structVariable, nil
}

// 读取所有栈帧中的变量列表
func (d *debugger) getAllFrameVariables(ctx context.Context) ([]*Variable, error) {
	variables := make([]*Variable, 0, 10)
	frames, err := d.GetStackTrace(ctx)
	if err != nil {
		logger.WithCtx(ctx).Errorf("getAllFrameVariables fail, err = %v", err)
		return nil, err
	}
	for _, frame := range frames {
		vs, err2 := d.GetFrameVariables(ctx, frame.ID)
		if err2 != nil {
			logger.WithCtx(ctx).Errorf("getAllFrameVariables fail, err = %v", err)
			return nil, err2
		}
		variables = append(variables, vs...)
	}
	return variables, nil
}

func (d *debugger) Array2DVisual(ctx context.Context, query *Array2DVisualQuery) (*Array2DVisualData, error) {
	logger.WithCtx(ctx).Infof("[GoDebugger] VariableVisual")
	stackTrace, err := d.GetStackTrace(ctx)
	if err != nil {
		return nil, err
	}

	variables, err := d.GetFrameVariables(ctx, stackTrace[0].ID)
	if err != nil {
		return nil, err
	}
	rowPointQuerySet := mapset.NewSet()
	columnQuerySet := mapset.NewSet()
	for _, point := range query.RowPointNames {
		rowPointQuerySet.Add(point)
	}
	for _, point := range query.ColPointNames {
		columnQuerySet.Add(point)
	}

	// 返回的 指针和可视化结构体
	rowPointVariables := make([]*VisualVariable, 0, 10)
	columnPointVariables := make([]*VisualVariable, 0, 10)
	// 目标结构体变量
	var structVariable *Variable

	for _, variable := range variables {
		if rowPointQuerySet.Contains(variable.Name) {
			visualVariable := NewVisualVariable(variable)
			rowPointVariables = append(rowPointVariables, visualVariable)
		}
		if columnQuerySet.Contains(variable.Name) {
			visualVariable := NewVisualVariable(variable)
			columnPointVariables = append(columnPointVariables, visualVariable)
		}
		if query.ArrayName == variable.Name {
			structVariable = variable
		}
	}

	// c语言做特殊处理
	if d.option.Language == constants.LanguageC {
		structVariable, err = d.getStructVariableForC(ctx, stackTrace, structVariable)
		if err != nil {
			logger.WithCtx(ctx).Errorf("[VariableVisual] getStructVariableForC fail, err=%v", err)
			return nil, err
		}
	}

	if structVariable == nil {
		return &Array2DVisualData{
			RowPoints: rowPointVariables,
			ColPoints: columnPointVariables,
			Array:     [][]*VisualVariable{},
		}, nil
	}
	arrases, err := d.GetVariables(ctx, structVariable.Reference)
	var array2d [][]*VisualVariable
	for _, array := range arrases {
		var structs []*VisualVariable
		switch d.option.Language {
		case constants.LanguageGo:
			structs = d.getArrayNodesForGo(ctx, array)
		case constants.LanguageC:
			structs = d.getArrayNodesForC(ctx, array)
		default:
			structs = d.getArrayNodesNormal(ctx, array)
		}
		array2d = append(array2d, structs)
	}
	return &Array2DVisualData{
		RowPoints: rowPointVariables,
		ColPoints: columnPointVariables,
		Array:     array2d,
	}, nil

}
