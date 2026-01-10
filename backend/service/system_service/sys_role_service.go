package system_service

import (
	"context"
	"errors"
	"github.com/fansqz/fancode-backend/common"
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/dao"
	"github.com/fansqz/fancode-backend/models/dto"
	"github.com/fansqz/fancode-backend/models/po"
	"gorm.io/gorm"
)

type SysRoleService interface {

	// GetRoleByID 根角色户id获取角色信息
	GetRoleByID(ctx context.Context, roleID uint) (*po.SysRole, error)
	// InsertSysRole 添加角色
	InsertSysRole(ctx context.Context, sysSysRole *po.SysRole) (uint, error)
	// UpdateSysRole 更新角色
	UpdateSysRole(ctx context.Context, SysRole *po.SysRole) error
	// DeleteSysRole 删除角色
	DeleteSysRole(ctx context.Context, id uint) error
	// GetSysRoleList 获取角色列表
	GetSysRoleList(ctx context.Context, pageQuery *dto.PageQuery) (*dto.PageInfo, error)
	// UpdateRoleApis 更新角色apis
	UpdateRoleApis(ctx context.Context, roleID uint, apiIDs []uint) error
	// UpdateRoleMenus 更新角色menu
	UpdateRoleMenus(ctx context.Context, roleID uint, menuIDs []uint) error
	// GetApiIDsByRoleID 通过角色id获取该角色拥有的apiID
	GetApiIDsByRoleID(ctx context.Context, roleID uint) ([]uint, error)
	// GetMenuIDsByRoleID 通过角色id获取该角色拥有的menuID
	GetMenuIDsByRoleID(ctx context.Context, roleID uint) ([]uint, error)
	// GetApisByRoleID 通过角色id获取该角色的所有api
	GetApisByRoleID(ctx context.Context, roleID uint) ([]*po.SysApi, error)
}

type sysRoleService struct {
	sysRoleDao dao.SysRoleDao
}

func NewSysRoleService(roleDao dao.SysRoleDao) SysRoleService {
	return &sysRoleService{
		sysRoleDao: roleDao,
	}
}

func (r *sysRoleService) GetRoleByID(ctx context.Context, roleID uint) (*po.SysRole, error) {
	role, err := r.sysRoleDao.GetRoleByID(common.Mysql, roleID)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.WithCtx(ctx).Errorf("[GetRoleByID] get role error, err = %v", err)
		return nil, e.ErrMysql
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, e.NewCustomMsg("The role is not exist")
	}
	return role, nil
}

func (r *sysRoleService) InsertSysRole(ctx context.Context, sysRole *po.SysRole) (uint, error) {
	// 添加
	err := r.sysRoleDao.InsertRole(common.Mysql, sysRole)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[InsertSysRole] insert role error, err = %v", err)
		return 0, e.ErrMysql
	}
	return sysRole.ID, nil
}

func (r *sysRoleService) UpdateSysRole(ctx context.Context, sysRole *po.SysRole) error {
	if err := r.sysRoleDao.UpdateRole(common.Mysql, sysRole); err != nil {
		logger.WithCtx(ctx).Errorf("[UpdateSysRole] update role error, err = %v", err)
		return e.ErrMysql
	}
	return nil
}

func (r *sysRoleService) DeleteSysRole(ctx context.Context, id uint) error {
	// 删除删除角色
	if err := r.sysRoleDao.DeleteRoleByID(common.Mysql, id); err != nil {
		logger.WithCtx(ctx).Errorf("[DeleteSysRole] delete role error, err = %v", err)
		return e.ErrMysql
	}
	return nil
}

func (r *sysRoleService) GetSysRoleList(ctx context.Context, query *dto.PageQuery) (*dto.PageInfo, error) {
	var roleQuery *po.SysRole
	if query.Query != nil {
		roleQuery = query.Query.(*po.SysRole)
	}
	// 获取角色列表
	sysSysRoles, err := r.sysRoleDao.GetRoleList(common.Mysql, query)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetSysRoleList] get role list error, err = %v", err)
		return nil, e.ErrMysql
	}
	newSysRoles := make([]*dto.SysRoleDtoForList, len(sysSysRoles))
	for i := 0; i < len(sysSysRoles); i++ {
		newSysRoles[i] = dto.NewSysRoleDtoForList(sysSysRoles[i])
	}
	// 获取所有角色总数目
	var count int64
	if count, err = r.sysRoleDao.GetRoleCount(common.Mysql, roleQuery); err != nil {
		logger.WithCtx(ctx).Errorf("[GetSysRoleList] get role count error, err = %v", err)
		return nil, e.ErrMysql
	}
	pageInfo := &dto.PageInfo{
		Total: count,
		Size:  int64(len(newSysRoles)),
		List:  newSysRoles,
	}
	return pageInfo, nil
}

func (r *sysRoleService) UpdateRoleApis(ctx context.Context, roleID uint, apiIDs []uint) error {
	tx := common.Mysql.Begin()
	if err := r.sysRoleDao.DeleteRoleAPIsByRoleID(tx, roleID); err != nil {
		logger.WithCtx(ctx).Errorf("[UpdateRoleApis] upadet role api error, err = %v", err)
		tx.Rollback()
		return e.ErrMysql
	}
	if err := r.sysRoleDao.InsertApisToRole(tx, roleID, apiIDs); err != nil {
		logger.WithCtx(ctx).Errorf("[UpdateRoleApis] insert role api error, err = %v", err)
		tx.Rollback()
		return e.ErrMysql
	}
	tx.Commit()
	return nil
}

func (r *sysRoleService) UpdateRoleMenus(ctx context.Context, roleID uint, menuIDs []uint) error {
	tx := common.Mysql.Begin()
	if err := r.sysRoleDao.DeleteRoleMenusByRoleID(tx, roleID); err != nil {
		logger.WithCtx(ctx).Errorf("[UpdateRoleApis] delete role menu error, err = %v", err)
		tx.Rollback()
		return e.ErrMysql
	}
	if err := r.sysRoleDao.InsertMenusToRole(tx, roleID, menuIDs); err != nil {
		logger.WithCtx(ctx).Errorf("[UpdateRoleApis] delete role menu error, err = %v", err)
		tx.Rollback()
		return e.ErrMysql
	}
	tx.Commit()
	return nil
}

func (r *sysRoleService) GetApiIDsByRoleID(ctx context.Context, roleID uint) ([]uint, error) {
	apiIDs, err := r.sysRoleDao.GetApiIDsByRoleID(common.Mysql, roleID)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetApiIDsByRoleID] get apis by role error, err = %v", err)
		return nil, e.ErrMysql
	}
	return apiIDs, nil
}

func (r *sysRoleService) GetMenuIDsByRoleID(ctx context.Context, roleID uint) ([]uint, error) {
	menuIDs, err := r.sysRoleDao.GetMenuIDsByRoleID(common.Mysql, roleID)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetMenuIDsByRoleID] get menus by role error, err = %v", err)
		return nil, e.ErrMysql
	}
	return menuIDs, nil
}

func (r *sysRoleService) GetApisByRoleID(ctx context.Context, roleID uint) ([]*po.SysApi, error) {
	apis, err := r.sysRoleDao.GetApisByRoleID(common.Mysql, roleID)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetMenuIDsByRoleID] get apis by role error, err = %v", err)
		return nil, e.ErrMysql
	}
	return apis, nil
}
