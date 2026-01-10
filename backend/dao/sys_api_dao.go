package dao

import (
	"github.com/fansqz/fancode-backend/models/po"
	"gorm.io/gorm"
)

type SysApiDao interface {
	// InsertApi 创建角色
	InsertApi(db *gorm.DB, api *po.SysApi) error
	// UpdateApi 修改api
	UpdateApi(db *gorm.DB, api *po.SysApi) error
	// DeleteApiByID 根据api的id进行删除
	DeleteApiByID(db *gorm.DB, id uint) error
	// GetApiByID 通过api的id获取api
	GetApiByID(db *gorm.DB, id uint) (*po.SysApi, error)
	// GetApiListByParentID 通过父id找到其所有子api
	GetApiListByParentID(db *gorm.DB, parentID int32) ([]*po.SysApi, error)
	// GetApiCount 获取api总数
	GetApiCount(db *gorm.DB) (int64, error)
	// GetApiListByPathKeyword 模糊查询api
	GetApiListByPathKeyword(db *gorm.DB, keyword string, page int, pageSize int) ([]*po.SysApi, error)
	// GetChildApisByParentID 根据父API的ID获取所有子API
	GetChildApisByParentID(db *gorm.DB, parentID uint) ([]*po.SysApi, error)
	// GetAllApi 获取所有api
	GetAllApi(db *gorm.DB) ([]*po.SysApi, error)
}

type sysApiDao struct {
}

func NewSysApiDao() SysApiDao {
	return &sysApiDao{}
}

func (s *sysApiDao) GetApiByID(db *gorm.DB, id uint) (*po.SysApi, error) {
	var api po.SysApi
	err := db.First(&api, id).Error
	return &api, err
}

func (s *sysApiDao) GetApiListByParentID(db *gorm.DB, parentID int32) ([]*po.SysApi, error) {
	var sysApis []*po.SysApi
	err := db.Where("parent_api_id = ?", parentID).Find(&sysApis).Error
	if err != nil {
		return nil, err
	}
	return sysApis, nil
}

func (s *sysApiDao) GetApiCount(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&po.SysApi{}).Count(&count).Error
	return count, err
}

func (s *sysApiDao) GetApiListByPathKeyword(db *gorm.DB, keyword string, page int, pageSize int) ([]*po.SysApi, error) {
	var sysApis []*po.SysApi
	err := db.Where("path LIKE ?", "%"+keyword+"%").Offset((page - 1) * pageSize).Limit(pageSize).Find(&sysApis).Error
	if err != nil {
		return nil, err
	}
	return sysApis, nil
}

func (s *sysApiDao) GetChildApisByParentID(db *gorm.DB, parentID uint) ([]*po.SysApi, error) {
	var childApis []*po.SysApi
	if err := db.Where("parent_api_id = ?", parentID).Find(&childApis).Error; err != nil {
		return nil, err
	}
	return childApis, nil
}

func (s *sysApiDao) GetAllApi(db *gorm.DB) ([]*po.SysApi, error) {
	var apiList []*po.SysApi
	err := db.Find(&apiList).Error
	return apiList, err
}

func (s *sysApiDao) InsertApi(db *gorm.DB, api *po.SysApi) error {
	return db.Create(api).Error
}

func (s *sysApiDao) UpdateApi(db *gorm.DB, api *po.SysApi) error {
	updMap := map[string]interface{}{
		"path":        api.Path,
		"method":      api.Method,
		"name":        api.Name,
		"description": api.Description,
	}
	return db.Model(api).Updates(updMap).Error
}

func (s *sysApiDao) DeleteApiByID(db *gorm.DB, id uint) error {
	return db.Delete(&po.SysApi{}, id).Error
}
