package dto

import (
	"github.com/fansqz/fancode-backend/constants"
	"github.com/fansqz/fancode-backend/service/visual_debug_servcie/debug_core"
)

// StructVisualRequest 可视化查询的参数
type StructVisualRequest struct {
	DebugID string                       `json:"debugID"`
	Query   debug_core.StructVisualQuery `json:"query"`
}

// ArrayVisualRequest
// 数组可视化请求
type ArrayVisualRequest struct {
	DebugID string                      `json:"debugID"`
	Query   debug_core.ArrayVisualQuery `json:"query"`
}

// Array2DVisualRequest
// 二维数组可视化请求
type Array2DVisualRequest struct {
	DebugID string                        `json:"debugID"`
	Query   debug_core.Array2DVisualQuery `json:"query"`
}

// VisualDescription 可视化描述
type VisualDescription struct {
	VisualType  constants.VisualType `json:"visualType"`
	Description interface{}          `json:"description"`
}

// ArrayDescription 数组可视化描述
type ArrayDescription struct {
	ArrayName  string   `json:"arrayName"`
	PointNames []string `json:"pointNames"`
	// 数组可视化展示类型，array普通展示，array-bar柱状图展示
	DisplayType string `json:"displayType"`
}

// Array2DDescription 二维数组可视化描述
type Array2DDescription struct {
	ArrayName     string   `json:"arrayName"`
	RowPointNames []string `json:"rowPointNames"`
	ColPointNames []string `json:"colPointNames"`
}

// BinaryTreeDescription 二叉树可视化描述
type BinaryTreeDescription struct {
	TreeNode string `json:"treeNode"` // 二叉树节点结构体名称
	Data     string `json:"data"`     // 数据域
	Left     string `json:"left"`     // 左子树属性名称
	Right    string `json:"right"`    // 右子树属性名称
}

// GraphDescription 图的可视化描述
type GraphDescription struct {
	GraphNode string   `json:"graphNode"` // 图节点结构体名称
	Data      string   `json:"data"`      // 数据域
	Nexts     []string `json:"nexts"`     // 邻接节点属性名
}

// LinkListDescription 链表的可视化描述
type LinkListDescription struct {
	LinkNode string `json:"linkNode"`       // 链表节点结构体名称
	Data     string `json:"data"`           // 数据域
	Next     string `json:"next"`           // next属性名
	Prev     string `json:"prev,omitempty"` // prev属性名，可选
}
