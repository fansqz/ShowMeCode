package visual_document_service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/fansqz/fancode-backend/common"
	conf "github.com/fansqz/fancode-backend/common/config"
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/dao"
	"github.com/fansqz/fancode-backend/models/dto"
	"github.com/fansqz/fancode-backend/models/po"
	"github.com/fansqz/fancode-backend/utils"
	"gorm.io/gorm"
	"sort"
)

type VisualDocumentService interface {
	// GetVisualDocumentDirectory 获取可视化文档目录
	GetVisualDocumentDirectory(ctx context.Context, bankID uint, enable *bool) ([]*dto.VisualDocumentDirectoryDto, error)
	// GetVisualDocumentByID 读取可视化文档
	GetVisualDocumentByID(ctx context.Context, id uint, enable *bool) (*dto.VisualDocumentDto, error)

	// InsertVisualDocument 添加可视化文档
	InsertVisualDocument(ctx context.Context, document *po.VisualDocument) (uint, error)
	// UpdateVisualDocument 更新可视化文档
	UpdateVisualDocument(ctx context.Context, document *dto.VisualDocumentDto) error
	// UpdateVisualDocumentDirectory 更新可视化文档的目录结构
	UpdateVisualDocumentDirectory(ctx context.Context, req *dto.UpdateVisualDocumentReq) error
	// DeleteVisualDocumentByID 删除可视化文档
	DeleteVisualDocumentByID(ctx context.Context, id uint) error
}

func NewVisualDocumentService(vdd dao.VisualDocumentDao, config *conf.AppConfig) VisualDocumentService {
	return &visualDocumentService{
		visualDocumentDao: vdd,
		config:            config,
	}
}

type visualDocumentService struct {
	visualDocumentDao dao.VisualDocumentDao
	config            *conf.AppConfig
}

func (v *visualDocumentService) GetVisualDocumentDirectory(ctx context.Context, bankID uint, enable *bool) ([]*dto.VisualDocumentDirectoryDto, error) {
	var documentList []*po.VisualDocument
	var err error
	if documentList, err = v.visualDocumentDao.GetAllSimpleDocument(common.Mysql, bankID); err != nil {
		logger.WithCtx(ctx).Errorf("[GetMenuTree] get all menu error, err = %v", err)
		return nil, e.ErrMysql
	}
	// 根据enable过滤
	if enable != nil {
		newDocumentList := make([]*po.VisualDocument, 0, len(documentList))
		for _, d := range documentList {
			if d.Enable == *enable {
				newDocumentList = append(newDocumentList, d)
			}
		}
		documentList = newDocumentList
	}
	// by order优先级排序
	sort.Slice(documentList, func(i, j int) bool {
		return documentList[i].Order > documentList[j].Order
	})

	documentMap := make(map[uint]*dto.VisualDocumentDirectoryDto)
	roots := make([]*dto.VisualDocumentDirectoryDto, 0, 10)

	// 添加到map中保存
	for _, document := range documentList {
		documentMap[document.ID] = dto.NewVisualDocumentDirectoryDto(document)
	}

	// 遍历并添加到父节点中
	for _, document := range documentList {
		if document.ParentID == 0 {
			roots = append(roots, documentMap[document.ID])
		} else {
			parent, exists := documentMap[document.ParentID]
			if !exists {
				continue
			}
			parent.Children = append(parent.Children, documentMap[document.ID])
		}
	}
	return roots, nil
}

func (v *visualDocumentService) GetVisualDocumentByID(ctx context.Context, id uint, enable *bool) (*dto.VisualDocumentDto, error) {
	document, err := v.visualDocumentDao.GetVisualDocumentByID(common.Mysql, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, e.NewRecordNotFoundErr("document not exist")
	}
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetVisualDocumentByID] GetVisualDocumentByID fail, err = %v", err)
		return nil, err
	}
	// 根据enable过滤
	if enable != nil && document.Enable != *enable {
		return nil, e.NewRecordNotFoundErr("document not exist")
	}
	// 获取可视化文档支持的所有语言
	codeList, err := v.visualDocumentDao.GetCodeListByDocumentID(common.Mysql, id)
	vcodes := make([]*dto.VisualDocumentCodeDto, 0, len(codeList))
	for _, code := range codeList {
		vcodes = append(vcodes, dto.NewVisualDocumentCodeDto(code))
	}
	answer := dto.NewVisualDocumentDto(document)
	answer.CodeList = vcodes
	return answer, nil
}

func (v *visualDocumentService) InsertVisualDocument(ctx context.Context, document *po.VisualDocument) (uint, error) {
	userID := utils.GetUserIDWithCtx(ctx)
	document.CreatorID = userID
	documentList, err := v.visualDocumentDao.GetSimpleDocumentByParentID(common.Mysql, document.ParentID)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[InsertVisualDocument] GetAllSimpleDocument fail, err = %v", err)
		return 0, e.ErrUnknown
	}
	// 设置文章的排序，新增的文章放在末尾
	sort.Slice(documentList, func(i, j int) bool {
		return documentList[i].Order > documentList[j].Order
	})
	if len(documentList) != 0 {
		document.Order = documentList[0].Order + 1
	} else {
		document.Order = 1
	}
	if err := v.visualDocumentDao.InsertVisualDocument(common.Mysql, document); err != nil {
		logger.WithCtx(ctx).Errorf("[InsertVisualDocument] InsertVisualDocument fail, err = %v", err)
		return 0, err
	}
	return document.ID, nil
}

func (v *visualDocumentService) UpdateVisualDocument(ctx context.Context, document *dto.VisualDocumentDto) error {
	err := common.Mysql.Transaction(func(tx *gorm.DB) error {
		// 更新可视化文档
		if err := v.visualDocumentDao.UpdateVisualDocument(common.Mysql, &po.VisualDocument{
			Model: gorm.Model{
				ID: document.ID,
			},
			ParentID: document.ParentID,
			Title:    document.Title,
			Content:  document.Content,
			Enable:   document.Enable,
		}); err != nil {
			logger.WithCtx(ctx).Errorf("[UpdateVisualDocument] UpdateVisualDocument fail, err = %v", err)
			return err
		}
		// 删除code
		if err := v.visualDocumentDao.DeleteVisualDocumentCodeByDocumentID(tx, document.ID); err != nil {
			return err
		}
		// 插入code
		for _, code := range document.CodeList {
			bps, err := json.Marshal(code.Breakpoints)
			if err != nil {
				return err
			}
			err = v.visualDocumentDao.InsertVisualDocumentCode(tx, &po.VisualDocumentCode{
				DocumentID:  document.ID,
				Code:        code.Code,
				Language:    code.Language,
				Breakpoints: string(bps),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.WithCtx(ctx).Errorf("[DeleteApiByID] delete api error, err = %v", err)
		return e.ErrMysql
	}
	return nil

}

func (v *visualDocumentService) UpdateVisualDocumentDirectory(ctx context.Context, req *dto.UpdateVisualDocumentReq) error {
	if req.EventType == "inner" {
		return v.innerDocument(ctx, req.DraggingDocumentID, req.DragDocumentID)
	} else if req.EventType == "before" {
		return v.afterOrBeforeDocument(ctx, req.DraggingDocumentID, req.DragDocumentID, false)
	} else {
		return v.afterOrBeforeDocument(ctx, req.DraggingDocumentID, req.DragDocumentID, true)
	}
}

func (v *visualDocumentService) innerDocument(ctx context.Context, draggingDocumentID uint, dragDocumentID uint) error {
	err := common.Mysql.Transaction(func(tx *gorm.DB) error {
		documentList, err := v.visualDocumentDao.GetSimpleDocumentByParentID(tx, dragDocumentID)
		if err != nil {
			return err
		}
		var order uint = 1000000
		sort.Slice(documentList, func(i, j int) bool {
			return documentList[i].Order > documentList[j].Order
		})
		if len(documentList) != 0 {
			order = documentList[0].Order - 1
		}
		if err = v.visualDocumentDao.UpdateVisualDocumentParentAndOrder(tx, draggingDocumentID, dragDocumentID, order); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.WithCtx(ctx).Errorf("[UpdateVisualDocument] innerDocument fail, err = %v", err)
	}
	return nil
}

func (v *visualDocumentService) afterOrBeforeDocument(ctx context.Context, draggingDocumentID uint, dragDocumentID uint, after bool) error {
	err := common.Mysql.Transaction(func(tx *gorm.DB) error {
		document, err := v.visualDocumentDao.GetVisualDocumentByID(tx, dragDocumentID)
		if err != nil {
			return err
		}
		parentID := document.ParentID
		documentList, err := v.visualDocumentDao.GetSimpleDocumentByParentID(tx, parentID)
		if err != nil {
			return err
		}
		sort.Slice(documentList, func(i, j int) bool {
			return documentList[i].Order > documentList[j].Order
		})
		documentIDList := make([]uint, 0, len(documentList)+1)
		for _, d := range documentList {
			if d.ID == draggingDocumentID {
				continue
			}
			if d.ID != dragDocumentID {
				documentIDList = append(documentIDList, d.ID)
				continue
			}
			if after {
				documentIDList = append(documentIDList, d.ID, draggingDocumentID)
			} else {
				documentIDList = append(documentIDList, draggingDocumentID, d.ID)
			}
		}
		// 更新order和parent
		for i, id := range documentIDList {
			if err = v.visualDocumentDao.UpdateVisualDocumentParentAndOrder(tx, id, parentID, 1000000-uint(i)); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.WithCtx(ctx).Errorf("[UpdateVisualDocument] afterOrBeforeDocument fail, err = %v", err)
	}
	return nil
}

func (v *visualDocumentService) DeleteVisualDocumentByID(ctx context.Context, id uint) error {
	err := common.Mysql.Transaction(func(tx *gorm.DB) error {
		// 删除code
		if err := v.visualDocumentDao.DeleteVisualDocumentCodeByDocumentID(tx, id); err != nil {
			return err
		}
		if err := v.visualDocumentDao.DeleteVisualDocumentByID(tx, id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.WithCtx(ctx).Errorf("[DeleteApiByID] delete api error, err = %v", err)
		return e.ErrMysql
	}
	return nil
}
