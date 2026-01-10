package dto

import (
	"github.com/fansqz/fancode-backend/models/po"
	"github.com/fansqz/fancode-backend/utils"
)

type VisualDocumentBankDto struct {
	ID          uint                     `json:"id"`
	Name        string                   `gorm:"column:name" json:"name"`
	Description string                   `gorm:"column:description" json:"description"`
	CreatedAt   utils.Time               `json:"createdAt"`
	UpdatedAt   utils.Time               `json:"updatedAt"`
	CreatorID   uint                     `json:"creatorID"`
	CreatorName string                   `json:"creatorName"`
	Enable      bool                     `json:"enable"`
	CodeList    []*VisualDocumentCodeDto `json:"codeList"`
}

func NewVisualDocumentBankDto(bank *po.VisualDocumentBank) *VisualDocumentBankDto {
	response := &VisualDocumentBankDto{
		ID:          bank.ID,
		Name:        bank.Name,
		Description: bank.Description,
		CreatedAt:   utils.Time(bank.CreatedAt),
		UpdatedAt:   utils.Time(bank.UpdatedAt),
		CreatorID:   bank.CreatorID,
		Enable:      bank.Enable,
	}
	return response
}
