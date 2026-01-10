package dto

import (
	"encoding/json"
	"github.com/fansqz/fancode-backend/models/po"
	"log"
)

type VisualDocumentDto struct {
	ID       uint                     `json:"id"`
	BankID   uint                     `json:"bankID"`
	ParentID uint                     `json:"parentID"`
	Title    string                   `json:"title"`
	Content  string                   `json:"content"`
	Enable   bool                     `json:"enable"`
	CodeList []*VisualDocumentCodeDto `json:"codeList"`
}

func NewVisualDocumentDto(document *po.VisualDocument) *VisualDocumentDto {
	return &VisualDocumentDto{
		BankID:   document.BankID,
		ID:       document.ID,
		ParentID: document.ParentID,
		Title:    document.Title,
		Content:  document.Content,
		Enable:   document.Enable,
	}
}

type VisualDocumentCodeDto struct {
	Code     string `gorm:"column:code" json:"code"`
	Language string `gorm:"column:language" json:"language"`
	// Breakpoints 初始断点
	Breakpoints []int `gorm:"column:breakpoints" json:"breakpoints"`
}

func NewVisualDocumentCodeDto(code *po.VisualDocumentCode) *VisualDocumentCodeDto {
	bps := []int{}
	if err := json.Unmarshal([]byte(code.Breakpoints), &bps); err != nil {
		log.Println(err)
	}
	return &VisualDocumentCodeDto{
		Code:        code.Code,
		Language:    code.Language,
		Breakpoints: bps,
	}
}

type VisualDocumentDirectoryDto struct {
	ID       uint                          `json:"id"`
	ParentID uint                          `json:"parentID"`
	Title    string                        `json:"title"`
	Enable   bool                          `json:"enable"`
	Order    uint                          `json:"order"`
	Children []*VisualDocumentDirectoryDto `json:"children"`
}

func NewVisualDocumentDirectoryDto(document *po.VisualDocument) *VisualDocumentDirectoryDto {
	return &VisualDocumentDirectoryDto{
		ID:       document.ID,
		ParentID: document.ParentID,
		Order:    document.Order,
		Title:    document.Title,
		Enable:   document.Enable,
	}
}

type UpdateVisualDocumentReq struct {
	// inner代表添加到DraggingDocument添加到DragDocument作为子节点
	// after代表DraggingDocument添加到DragDocument后
	// before代表DraggingDocument添加到DragDocument前
	EventType          string `json:"eventType"`
	DraggingDocumentID uint   `json:"draggingDocumentID"`
	DragDocumentID     uint   `json:"dragDocumentID"`
}
