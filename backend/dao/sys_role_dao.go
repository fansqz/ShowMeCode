package dao

import (
	"github.com/fansqz/fancode-backend/models/dto"
	"github.com/fansqz/fancode-backend/models/po"
	"gorm.io/gorm"
)

type SysRoleDao interface {
	// InsertRole 创建角色
	InsertRole(db *gorm.DB, role *po.SysRole) error
	// UpdateRole 更新角色
	UpdateRole(db *gorm.DB, role *po.SysRole) error
	// DeleteRoleByID 删除角色
	DeleteRoleByID(db *gorm.DB, id uint) error
	// GetRoleByID 通过角色id获取角色
	GetRoleByID(db *gorm.DB, roleID uint) (*po.SysRole, error)
	// GetRoleList 获取角色列表
	GetRoleList(db *gorm.DB, pageQuery *dto.PageQuery) ([]*po.SysRole, error)
	// GetRoleCount 获取所有角色数量
	GetRoleCount(db *gorm.DB, role *po.SysRole) (int64, error)
	// InsertMenusToRole 给角色添加menu
	InsertMenusToRole(db *gorm.DB, roleID uint, menuIDs []uint) error
	// DeleteRoleMenusByRoleID 清除所有与roleID关联的roleID-menuID数据
	DeleteRoleMenusByRoleID(db *gorm.DB, roleID uint) error
	// GetMenuIDsByRoleID 获取用户关联的所有menu的id
	GetMenuIDsByRoleID(db *gorm.DB, roleID uint) ([]uint, error)
	// GetMenusByRoleID 通过用户角色获取菜单列表
	GetMenusByRoleID(db *gorm.DB, roleID uint) ([]*po.SysMenu, error)
	// InsertApisToRole 给角色添加api
	InsertApisToRole(db *gorm.DB, roleID uint, apiIDs []uint) error
	// DeleteRoleAPIsByRoleID 清除所有与roleID关联的roleID-apiID数据
	DeleteRoleAPIsByRoleID(db *gorm.DB, roleID uint) error
	// GetApiIDsByRoleID 获取用户关联的所有api的id
	GetApiIDsByRoleID(db *gorm.DB, roleID uint) ([]uint, error)
	// GetApisByRoleID 获取用户关联的所有api
	GetApisByRoleID(db *gorm.DB, roleID uint) ([]*po.SysApi, error)
	// GetAllSimpleRoleList 获取所有角色列表，只含有id和name
	GetAllSimpleRoleList(db *gorm.DB) ([]*po.SysRole, error)
}

type sysRoleDao struct {
}

func NewSysRoleDao() SysRoleDao {
	return &sysRoleDao{}
}

func (r *sysRoleDao) InsertRole(db *gorm.DB, role *po.SysRole) error {
	return db.Create(role).Error
}

func (r *sysRoleDao) UpdateRole(db *gorm.DB, role *po.SysRole) error {
	return db.Model(role).Updates(role).Error
}

func (r *sysRoleDao) DeleteRoleByID(db *gorm.DB, id uint) error {
	return db.Delete(&po.SysRole{}, id).Error
}

func (r *sysRoleDao) GetRoleByID(db *gorm.DB, roleID uint) (*po.SysRole, error) {
	var role po.SysRole
	err := db.First(&role, roleID).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *sysRoleDao) GetRoleList(db *gorm.DB, pageQuery *dto.PageQuery) ([]*po.SysRole, error) {
	var role *po.SysRole
	if pageQuery.Query != nil {
		role = pageQuery.Query.(*po.SysRole)
	}
	offset := (pageQuery.Page - 1) * pageQuery.PageSize
	var roles []*po.SysRole
	if role != nil && role.Name != "" {
		db = db.Where("name LIKE ?", "%"+role.Name+"%")
	}
	if role != nil && role.Description != "" {
		db = db.Where("description LIKE ?", "%"+role.Description+"%")
	}
	db = db.Limit(pageQuery.PageSize).Offset(offset)
	if pageQuery.SortProperty != "" && pageQuery.SortRule != "" {
		order := pageQuery.SortProperty + " " + pageQuery.SortRule
		db = db.Order(order)
	}
	err := db.Find(&roles).Error
	return roles, err
}

func (r *sysRoleDao) GetRoleCount(db *gorm.DB, role *po.SysRole) (int64, error) {
	var count int64
	if role != nil && role.Name != "" {
		db = db.Where("name LIKE ?", "%"+role.Name+"%")
	}
	if role != nil && role.Description != "" {
		db = db.Where("description LIKE ?", "%"+role.Description+"%")
	}
	err := db.Model(&po.SysRole{}).Count(&count).Error
	return count, err
}

func (r *sysRoleDao) InsertMenusToRole(db *gorm.DB, roleID uint, menuIDs []uint) error {
	role := &po.SysRole{}
	role.ID = roleID
	var menus []po.SysMenu
	for _, menuID := range menuIDs {
		menu := po.SysMenu{}
		menu.ID = menuID
		menus = append(menus, menu)
	}
	err := db.Model(role).Association("Menus").Append(menus)
	return err
}

func (r *sysRoleDao) DeleteRoleMenusByRoleID(db *gorm.DB, roleID uint) error {
	role := po.SysRole{}
	role.ID = roleID
	if err := db.Model(&role).Association("Menus").Clear(); err != nil {
		return err
	}
	return nil
}

func (r *sysRoleDao) GetMenuIDsByRoleID(db *gorm.DB, roleID uint) ([]uint, error) {
	var role po.SysRole
	role.ID = roleID
	if err := db.Model(&role).Association("Menus").Find(&role.Menus); err != nil {
		return nil, err
	}
	menuIDs := make([]uint, len(role.Menus))
	for i, menu := range role.Menus {
		menuIDs[i] = menu.ID
	}
	return menuIDs, nil
}

func (r *sysRoleDao) GetMenusByRoleID(db *gorm.DB, roleID uint) ([]*po.SysMenu, error) {
	var role po.SysRole
	role.ID = roleID
	if err := db.Model(&role).Association("Menus").Find(&role.Menus); err != nil {
		return nil, err
	}
	return role.Menus, nil
}

func (r *sysRoleDao) InsertApisToRole(db *gorm.DB, roleID uint, apiIDs []uint) error {
	role := &po.SysRole{}
	role.ID = roleID
	var apis []po.SysApi
	for _, apiID := range apiIDs {
		api := po.SysApi{}
		api.ID = apiID
		apis = append(apis, api)
	}
	err := db.Model(role).Association("Apis").Append(apis)
	return err
}

func (r *sysRoleDao) DeleteRoleAPIsByRoleID(db *gorm.DB, roleID uint) error {
	role := po.SysRole{}
	role.ID = roleID
	if err := db.Model(&role).Association("Apis").Clear(); err != nil {
		return err
	}

	return nil
}

func (r *sysRoleDao) GetApiIDsByRoleID(db *gorm.DB, roleID uint) ([]uint, error) {
	var role po.SysRole
	role.ID = roleID
	if err := db.Model(&role).Association("Apis").Find(&role.Apis); err != nil {
		return nil, err
	}
	apiIDs := make([]uint, len(role.Apis))
	for i, api := range role.Apis {
		apiIDs[i] = api.ID
	}
	return apiIDs, nil
}

func (r *sysRoleDao) GetApisByRoleID(db *gorm.DB, roleID uint) ([]*po.SysApi, error) {
	var role po.SysRole
	role.ID = roleID
	if err := db.Model(&role).Association("Apis").Find(&role.Apis); err != nil {
		return nil, err
	}
	return role.Apis, nil
}

func (r *sysRoleDao) GetAllSimpleRoleList(db *gorm.DB) ([]*po.SysRole, error) {
	var roles []*po.SysRole
	err := db.Select([]string{"id", "name"}).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
