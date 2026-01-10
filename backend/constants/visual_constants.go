package constants

// VisualType 可视化类型
type VisualType string

const (
	// ArrayType 数组
	ArrayType VisualType = "array"
	// Array2DType 二维数组
	Array2DType VisualType = "array2d"
	// BinaryTreeType 二叉树
	BinaryTreeType VisualType = "binaryTree"
	// GraphType 图
	GraphType VisualType = "graph"
	// LinkListType 链表
	LinkListType VisualType = "linkList"
)

var VisualTypeList = [5]VisualType{
	ArrayType,
	Array2DType,
	BinaryTreeType,
	GraphType,
	LinkListType,
}
