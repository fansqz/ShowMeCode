package debugger

// StructVisualQuery 可视化查询的参数
// 结构体导向，会根据结构体去遍历这个结构体的所有数据和指针
// 1. 该方法会找出所有指向该结构体的指针
// 2. 深度遍历所有指定type的结构体数据并返回
type StructVisualQuery struct {
	// 需要查询的结构体名称
	Struct string `json:"struct"`
	// Values 数据域
	Values []string `json:"values"`
	// Points 指针域，指针域如果是一个数组，那么这个数组下所有元素都将作为指针
	Points []string `json:"points"`
}

// StructVisualData 可视化拆查询返回的数据
type StructVisualData struct {
	// Nodes 可视化结构的节点列表
	Nodes []*VisualNode `json:"nodes"`
	// Points 变量列表
	Points []*VisualVariable `json:"points"`
}

// VariableVisualQuery
// 变量为导向，查询变量的值作为可视化数据，数组类型会使用这种
type VariableVisualQuery struct {
	// StructVars 作为结构体的变量
	StructVars []string `json:"structVars"`
	// PointVars 作为指针的变量
	PointVars []string `json:"pointVars"`
}

// VariableVisualData
// 变量为导向，查询变量的值作为可视化数据
type VariableVisualData struct {
	Structs []*VisualNode     `json:"structs"`
	Points  []*VisualVariable `json:"points"`
}

// VisualVariable
type VisualVariable struct {
	// 变量名称
	Name string `json:"name"`
	// 变量类型
	Type string `json:"type"`
	// 变量的值，在可视化中大部分作为指针使用
	Value string `json:"value"`
}

// VisualNode 可视化的一个节点
// 包含所有的数据域和指针域
type VisualNode struct {
	Name string `json:"name"`
	// 可以理解为地址
	ID     string            `json:"id"`
	Type   string            `json:"type"`
	Values []*VisualVariable `json:"values"`
	Points []*VisualVariable `json:"points"`
}

func NewVisualVariable(variable *Variable) *VisualVariable {
	return &VisualVariable{
		Name:  variable.Name,
		Type:  variable.Type,
		Value: *variable.Value,
	}
}
