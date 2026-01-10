package system_service

import (
	"context"

	"github.com/fansqz/fancode-backend/common"
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/dao"
	"github.com/fansqz/fancode-backend/models/dto"
	"github.com/fansqz/fancode-backend/models/po"
	"gorm.io/gorm"
)

type SysApiService interface {

	// GetApiCount 获取api数目
	GetApiCount(ctx context.Context) (int64, error)
	// DeleteApiByID 删除api
	DeleteApiByID(ctx context.Context, id uint) error
	// UpdateApi 更新api
	UpdateApi(ctx context.Context, api *po.SysApi) error
	// GetApiByID 根据id获取api
	GetApiByID(ctx context.Context, id uint) (*po.SysApi, error)
	// GetApiTree 获取api树
	GetApiTree(ctx context.Context) ([]*dto.SysApiTreeDto, error)
	// InsertApi 添加api
	InsertApi(ctx context.Context, api *po.SysApi) (uint, error)
}

type sysApiService struct {
	sysApiDao dao.SysApiDao
}

func NewSysApiService(apiDao dao.SysApiDao) SysApiService {
	return &sysApiService{
		sysApiDao: apiDao,
	}
}

func (s *sysApiService) GetApiCount(ctx context.Context) (int64, error) {
	count, err := s.sysApiDao.GetApiCount(common.Mysql)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetApiCount] get api count error, err = %v", err)
		return 0, e.ErrMysql
	}
	return count, nil
}

// DeleteApiByID 根据api的id进行删除
func (s *sysApiService) DeleteApiByID(ctx context.Context, id uint) error {
	err := common.Mysql.Transaction(func(tx *gorm.DB) error {
		// 递归删除API
		err := s.deleteApisRecursive(ctx, tx, id)
		return err
	})
	if err != nil {
		logger.WithCtx(ctx).Errorf("[DeleteApiByID] delete api error, err = %v", err)
		return e.ErrMysql
	}
	return nil
}

// deleteApisRecursive 递归删除API
func (s *sysApiService) deleteApisRecursive(ctx context.Context, db *gorm.DB, parentID uint) error {
	childApis, err := s.sysApiDao.GetChildApisByParentID(db, parentID)
	if err != nil {
		return err
	}
	for _, childAPI := range childApis {
		// 删除子api的子api
		if err = s.deleteApisRecursive(ctx, db, childAPI.ID); err != nil {
			return err
		}
	}
	// 当前api
	if err = s.sysApiDao.DeleteApiByID(db, parentID); err != nil {
		logger.WithCtx(ctx).Errorf("[deleteApisRecursive] delete api error, err = %v", err)
		return err
	}
	return nil
}

func (s *sysApiService) UpdateApi(ctx context.Context, api *po.SysApi) error {
	err := s.sysApiDao.UpdateApi(common.Mysql, api)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[UpdateApi] update api error, err = %v", err)
		return e.ErrMysql
	}
	return nil
}

func (s *sysApiService) GetApiByID(ctx context.Context, id uint) (*po.SysApi, error) {
	api, err := s.sysApiDao.GetApiByID(common.Mysql, id)
	if err == gorm.ErrRecordNotFound {
		return nil, e.NewCustomMsg("The api is not exist")
	}
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetApiByID] update api error, err = %v", err)
		return nil, e.ErrMysql
	}
	return api, nil
}

func (s *sysApiService) GetApiTree(ctx context.Context) ([]*dto.SysApiTreeDto, error) {
	var apiList []*po.SysApi
	var err error
	if apiList, err = s.sysApiDao.GetAllApi(common.Mysql); err != nil {
		logger.WithCtx(ctx).Errorf("[GetApiTree] get all api error, err = %v", err)
		return nil, e.ErrMysql
	}

	apiMap := make(map[uint]*dto.SysApiTreeDto)
	var rootApis []*dto.SysApiTreeDto

	// 添加到map中保存
	for _, api := range apiList {
		apiMap[api.ID] = dto.NewSysApiTreeDto(api)
	}

	// 遍历并添加到父节点中
	for _, api := range apiList {
		if api.ParentApiID == 0 {
			rootApis = append(rootApis, apiMap[api.ID])
		} else {
			parentApi, exists := apiMap[api.ParentApiID]
			if !exists {
				return nil, e.ErrUnknown
			}
			parentApi.Children = append(parentApi.Children, apiMap[api.ID])
		}
	}

	return rootApis, nil
}

func (s *sysApiService) InsertApi(ctx context.Context, api *po.SysApi) (uint, error) {
	err := s.sysApiDao.InsertApi(common.Mysql, api)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[InsertApi] insert api error, err = %v", err)
		return 0, e.ErrMysql
	}
	return api.ID, nil
}
