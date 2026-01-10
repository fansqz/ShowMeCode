package gdb_debugger

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
)

const (
	globalScopeReference = 1001
)

type ReferenceType string

const (
	StructType ReferenceType = "v"
	PointType  ReferenceType = "p"
)

// ReferenceUtil 引用工具类
type ReferenceUtil struct {
	nextRef       int
	mutex         sync.RWMutex
	refInt2Struct map[int]string
	refStruct2Int map[string]int
}

// ReferenceStruct 定义的引用结构体
type ReferenceStruct struct {
	Type         ReferenceType
	FrameId      string
	VariableName string
	VariableType string
	Address      string
	FieldPath    string
}

// NewPointReferenceStruct 创建指针引用结构体
func NewPointReferenceStruct(variableType string, address string) *ReferenceStruct {
	return &ReferenceStruct{
		Type:         PointType,
		VariableType: variableType,
		Address:      address,
	}
}

// NewStructReferenceStruct 创建结构体引用结构体
func NewStructReferenceStruct(frameId string, variableName string, variableType string) *ReferenceStruct {
	return &ReferenceStruct{
		Type:         StructType,
		FrameId:      frameId,
		VariableName: variableName,
		VariableType: variableType,
	}
}

func NewReferenceUtil() *ReferenceUtil {
	return &ReferenceUtil{
		nextRef:       1100,
		refInt2Struct: map[int]string{},
		refStruct2Int: map[string]int{},
	}
}

// GetScopesReference 根据栈帧获取ScopeId
func (r *ReferenceUtil) GetScopesReference(frameId int) int {
	return 1002 + frameId
}

// CheckIsScopeReference 判断是否是Scope引用
func (r *ReferenceUtil) CheckIsScopeReference(reference int) bool {
	return reference < 1100
}

// CheckIsGlobalScope  判断是否是GlobalScope引用
func (r *ReferenceUtil) CheckIsGlobalScope(reference int) bool {
	return reference == 1001
}

// CheckIsLocalScope  判断是否是LocalScope引用
func (r *ReferenceUtil) CheckIsLocalScope(reference int) bool {
	return reference < 1100 && reference > 1001
}

// GetFrameIDByLocalReference 获取栈帧id
func (r *ReferenceUtil) GetFrameIDByLocalReference(reference int) int {
	return reference - 1002
}

// ParseVariableReference 解析引用
func (r *ReferenceUtil) ParseVariableReference(reference int) (*ReferenceStruct, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	refStr, ok := r.refInt2Struct[reference]
	if !ok {
		return nil, fmt.Errorf("reference not found")
	}
	return r.parseReference(refStr)
}

// CreateVariableReference 创建引用
func (r *ReferenceUtil) CreateVariableReference(refStruct *ReferenceStruct) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	strRef, err := r.convertReference(refStruct)
	if err != nil {
		log.Printf("Error converting reference %v to struct: %v", refStruct, err)
		return 0, err
	}
	// 如果引用已经存在，直接返回
	if _, ok := r.refStruct2Int[strRef]; ok {
		return r.refStruct2Int[strRef], nil
	}
	// 创建引用
	intRef := r.nextRef
	r.nextRef++
	r.refStruct2Int[strRef] = intRef
	r.refInt2Struct[intRef] = strRef
	return intRef, nil
}

// parseReference
// 目前一共有两种类型的引用，分别是：
// 1. 普通类型：v-<栈帧id>-<变量名>-<属性（如果有）>;
// 2.指针：p-<指针类型>-<地址>-<地址的变量名称>-<属性（如果有）>
func (r *ReferenceUtil) parseReference(reference string) (*ReferenceStruct, error) {
	answer := &ReferenceStruct{}
	err := json.Unmarshal([]byte(reference), answer)
	if err != nil {
		return nil, err
	}
	return answer, nil
}

// convertReference 把refStruct结构体转成引用字符串
func (r *ReferenceUtil) convertReference(refStruct *ReferenceStruct) (string, error) {
	answer, err := json.Marshal(refStruct)
	return string(answer), err
}

func GetFieldReferenceStruct(refStruct *ReferenceStruct, fieldName string) *ReferenceStruct {
	newRef := &ReferenceStruct{
		Type:         refStruct.Type,
		FrameId:      refStruct.FrameId,
		VariableName: refStruct.VariableName,
		VariableType: refStruct.VariableType,
		Address:      refStruct.Address,
	}
	if refStruct.FieldPath == "" {
		if checkIsNumber(fieldName) {
			newRef.FieldPath = fmt.Sprintf("[%s]", fieldName)
		} else {
			newRef.FieldPath = fmt.Sprintf(".%s", fieldName)
		}
	} else {
		if checkIsNumber(fieldName) {
			newRef.FieldPath = fmt.Sprintf("%s[%s]", refStruct.FieldPath, fieldName)
		} else {
			newRef.FieldPath = fmt.Sprintf("%s.%s", refStruct.FieldPath, fieldName)
		}
	}
	return newRef
}

// 校验string是否是数值类型
func checkIsNumber(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}
