package system_service

import (
	"context"
	"github.com/fansqz/fancode-backend/common"
	conf "github.com/fansqz/fancode-backend/common/config"
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/dao"
	"github.com/fansqz/fancode-backend/models/dto"
	"github.com/fansqz/fancode-backend/models/po"
	"github.com/fansqz/fancode-backend/utils"
	"gorm.io/gorm"
	"time"
)

type SysUserService interface {
	// GetUserByID 根据用户id获取用户信息
	GetUserByID(ctx context.Context, userID uint) (*po.SysUser, error)
	// InsertSysUser 添加用户
	InsertSysUser(ctx context.Context, sysUser *po.SysUser) (uint, error)
	// UpdateSysUser 更新用户，但是不更新密码
	UpdateSysUser(ctx context.Context, sysUser *po.SysUser) error
	// DeleteSysUser 删除用户
	DeleteSysUser(ctx context.Context, id uint) error
	// GetSysUserList 获取用户列表
	GetSysUserList(ctx context.Context, pageQuery *dto.PageQuery) (*dto.PageInfo, error)
	// UpdateUserRoles 更新角色roleIDs
	UpdateUserRoles(ctx context.Context, userID uint, roleIDs []uint) error
	// GetRoleIDsByUserID 通过用户id获取所有角色id
	GetRoleIDsByUserID(ctx context.Context, userID uint) ([]uint, error)
	// GetAllSimpleRole
	GetAllSimpleRole(ctx context.Context) ([]*dto.SimpleRoleDto, error)
}

type sysUserService struct {
	config     *conf.AppConfig
	sysUserDao dao.SysUserDao
	sysRoleDao dao.SysRoleDao
}

func NewSysUserService(config *conf.AppConfig, userDao dao.SysUserDao, roleDao dao.SysRoleDao) SysUserService {
	return &sysUserService{
		config:     config,
		sysUserDao: userDao,
		sysRoleDao: roleDao,
	}
}

func (s *sysUserService) GetUserByID(ctx context.Context, userID uint) (*po.SysUser, error) {
	user, err := s.sysUserDao.GetUserByID(common.Mysql, userID)
	if err == gorm.ErrRecordNotFound {
		return nil, e.ErrUserNotExist
	}
	if err != nil {
		return nil, e.ErrMysql
	}
	return user, nil
}

func (s *sysUserService) InsertSysUser(ctx context.Context, sysUser *po.SysUser) (uint, error) {
	// 设置默认用户名
	if sysUser.Username == "" {
		sysUser.Username = "fancoder"
	}
	// 随机登录名称
	if sysUser.LoginName == "" {
		sysUser.LoginName = sysUser.LoginName + utils.GetUUID()
	}
	// 设置默认密码
	if sysUser.Password == "" {
		sysUser.Password = s.config.DefaultPassword
	}
	// 设置默认出生时间
	t := time.Time{}
	if sysUser.BirthDay == t {
		sysUser.BirthDay = time.Now()
	}
	// 设置默认性别
	if sysUser.Sex != 1 && sysUser.Sex != 2 {
		sysUser.Sex = 1
	}
	p, err := utils.GetPwd(sysUser.Password)
	if err != nil {
		return 0, e.ErrMysql
	}
	sysUser.Password = string(p)
	err = s.sysUserDao.InsertUser(common.Mysql, sysUser)
	if err != nil {
		return 0, e.ErrMysql
	}
	return sysUser.ID, nil
}

func (s *sysUserService) UpdateSysUser(ctx context.Context, sysUser *po.SysUser) error {
	err := s.sysUserDao.UpdateUser(common.Mysql, sysUser)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (s *sysUserService) DeleteSysUser(ctx context.Context, id uint) error {
	err := s.sysUserDao.DeleteUserByID(common.Mysql, id)
	if err != nil {
		return e.ErrMysql
	}
	return nil
}

func (s *sysUserService) GetSysUserList(ctx context.Context, pageQuery *dto.PageQuery) (*dto.PageInfo, error) {
	var pageInfo *dto.PageInfo
	var userQuery *po.SysUser
	if pageQuery.Query != nil {
		userQuery = pageQuery.Query.(*po.SysUser)
	}
	err := common.Mysql.Transaction(func(tx *gorm.DB) error {
		userList, err := s.sysUserDao.GetUserList(tx, pageQuery)
		if err != nil {
			return err
		}
		userDtoList := make([]*dto.SysUserDtoForList, len(userList))
		for i, user := range userList {
			user.Roles, err = s.sysUserDao.GetRolesByUserID(tx, user.ID)
			if err != nil {
				return err
			}
			userDtoList[i] = dto.NewSysUserDtoForList(user)
		}
		var count int64
		count, err = s.sysUserDao.GetUserCount(tx, userQuery)
		if err != nil {
			return err
		}
		pageInfo = &dto.PageInfo{
			Total: count,
			Size:  int64(len(userDtoList)),
			List:  userDtoList,
		}
		return nil
	})
	if err != nil {
		return nil, e.ErrMysql
	}
	return pageInfo, nil
}

func (s *sysUserService) UpdateUserRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	tx := common.Mysql.Begin()
	err := s.sysUserDao.DeleteUserRoleByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return e.ErrMysql
	}
	err = s.sysUserDao.InsertRolesToUser(tx, userID, roleIDs)
	if err != nil {
		tx.Rollback()
		return e.ErrMysql
	}
	tx.Commit()
	return nil
}

func (s *sysUserService) GetRoleIDsByUserID(ctx context.Context, userID uint) ([]uint, error) {
	roleIDs, err := s.sysUserDao.GetRoleIDsByUserID(common.Mysql, userID)
	if err != nil {
		return nil, e.ErrMysql
	}
	return roleIDs, nil
}

func (s *sysUserService) GetAllSimpleRole(ctx context.Context) ([]*dto.SimpleRoleDto, error) {
	roles, err := s.sysRoleDao.GetAllSimpleRoleList(common.Mysql)
	if err != nil {
		return nil, e.ErrMysql
	}
	simpleRoles := make([]*dto.SimpleRoleDto, len(roles))
	for i, role := range roles {
		simpleRoles[i] = dto.NewSimpleRoleDto(role)
	}
	return simpleRoles, nil
}
