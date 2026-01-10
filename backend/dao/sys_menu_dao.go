package dao

import (
	"github.com/fansqz/fancode-backend/models/po"
	"gorm.io/gorm"
)

type SysMenuDao interface {
	// InsertMenu 创建角色
	InsertMenu(db *gorm.DB, menu *po.SysMenu) error
	// GetMenuByID 通过menu的id获取menu
	GetMenuByID(db *gorm.DB, id uint) (*po.SysMenu, error)
	// GetMenuListByParentID 通过父id找到其所有子menu
	GetMenuListByParentID(db *gorm.DB, parentID int32) ([]*po.SysMenu, error)
	// GetMenuCount 获取menu总数
	GetMenuCount(db *gorm.DB) (int64, error)
	// GetMenuListByPathKeyword 模糊查询menu
	GetMenuListByPathKeyword(db *gorm.DB, keyword string, page int, pageSize int) ([]*po.SysMenu, error)
	// DeleteMenuByID 根据menu的id进行删除
	DeleteMenuByID(db *gorm.DB, id uint) error
	// UpdateMenu 修改menu
	UpdateMenu(db *gorm.DB, menu *po.SysMenu) error
	// GetChildMenusByParentID 根据父API的ID获取所有子API
	GetChildMenusByParentID(db *gorm.DB, parentID uint) ([]*po.SysMenu, error)
	// GetAllMenu 获取所有menu
	GetAllMenu(db *gorm.DB) ([]*po.SysMenu, error)
}

type sysMenuDao struct {
}

func NewSysMenuDao() SysMenuDao {
	return &sysMenuDao{}
}

func (s *sysMenuDao) InsertMenu(db *gorm.DB, menu *po.SysMenu) error {
	return db.Create(menu).Error
}

func (s *sysMenuDao) GetMenuByID(db *gorm.DB, id uint) (*po.SysMenu, error) {
	var menu po.SysMenu
	err := db.First(&menu, id).Error
	return &menu, err
}

func (s *sysMenuDao) GetMenuListByParentID(db *gorm.DB, parentID int32) ([]*po.SysMenu, error) {
	var sysMenus []*po.SysMenu
	err := db.Where("parent_menu_id = ?", parentID).Find(&sysMenus).Error
	if err != nil {
		return nil, err
	}
	return sysMenus, nil
}

func (s *sysMenuDao) GetMenuCount(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&po.SysMenu{}).Count(&count).Error
	return count, err
}

func (s *sysMenuDao) GetMenuListByPathKeyword(db *gorm.DB, keyword string, page int, pageSize int) ([]*po.SysMenu, error) {
	var sysMenus []*po.SysMenu
	err := db.Where("path LIKE ?", "%"+keyword+"%").Offset((page - 1) * pageSize).Limit(pageSize).Find(&sysMenus).Error
	if err != nil {
		return nil, err
	}
	return sysMenus, nil
}

func (s *sysMenuDao) DeleteMenuByID(db *gorm.DB, id uint) error {
	return db.Delete(&po.SysMenu{}, id).Error
}

func (s *sysMenuDao) UpdateMenu(db *gorm.DB, menu *po.SysMenu) error {
	updMap := map[string]interface{}{
		"name":        menu.Name,
		"code":        menu.Code,
		"description": menu.Description,
	}
	return db.Model(menu).Updates(updMap).Error
}

func (s *sysMenuDao) GetChildMenusByParentID(db *gorm.DB, parentID uint) ([]*po.SysMenu, error) {
	var childMenus []*po.SysMenu
	if err := db.Where("parent_menu_id = ?", parentID).Find(&childMenus).Error; err != nil {
		return nil, err
	}
	return childMenus, nil
}

func (s *sysMenuDao) GetAllMenu(db *gorm.DB) ([]*po.SysMenu, error) {
	var menuList []*po.SysMenu
	err := db.Find(&menuList).Error
	return menuList, err
}
