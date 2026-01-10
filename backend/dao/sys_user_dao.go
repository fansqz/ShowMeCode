package dao

import (
	"errors"
	"github.com/fansqz/fancode-backend/models/dto"
	"github.com/fansqz/fancode-backend/models/po"
	"gorm.io/gorm"
)

type SysUserDao interface {
	// InsertUser 创建用户
	InsertUser(db *gorm.DB, user *po.SysUser) error
	// UpdateUser
	UpdateUser(db *gorm.DB, user *po.SysUser) error
	// DeleteUserByID 删除用户
	DeleteUserByID(db *gorm.DB, id uint) error
	// GetUserByID 通过用户id获取用户
	GetUserByID(db *gorm.DB, id uint) (*po.SysUser, error)
	// GetUserNameByID 通过用户id获取用户名称
	GetUserNameByID(db *gorm.DB, id uint) (string, error)
	// GetUserList 获取用户列表
	GetUserList(db *gorm.DB, pageQuery *dto.PageQuery) ([]*po.SysUser, error)
	// GetUserCount 获取所有用户数量
	GetUserCount(db *gorm.DB, user *po.SysUser) (int64, error)
	// GetUserByLoginName 根据用户登录名称获取用户
	GetUserByLoginName(db *gorm.DB, loginName string) (*po.SysUser, error)
	// CheckEmail 检测邮箱是否已经存在
	CheckEmailIsExist(db *gorm.DB, email string) (bool, error)
	// DeleteUserRoleByUserID 清除所有与userID关联的userID-roleID数据
	DeleteUserRoleByUserID(db *gorm.DB, userID uint) error
	// GetRolesByUserID 获取用户关联的所有role的名称
	GetRolesByUserID(db *gorm.DB, userID uint) ([]*po.SysRole, error)
	// GetUserByEmail 根据用户邮箱获取用户信息
	GetUserByEmail(db *gorm.DB, email string) (*po.SysUser, error)
	// CheckLoginName 检测loginname是否存在
	CheckLoginName(db *gorm.DB, loginname string) (bool, error)
	// GetRoleIDsByUserID 获取用户关联的所有role的id
	GetRoleIDsByUserID(db *gorm.DB, userID uint) ([]uint, error)
	// InsertRolesToUser 给用户设置角色
	InsertRolesToUser(db *gorm.DB, userID uint, roleIDs []uint) error
	// UpdateUserPassword 更新用户密码
	UpdateUserPassword(db *gorm.DB, userID uint, password string) error
}

type sysUserDao struct {
}

func NewSysUserDao() SysUserDao {
	return &sysUserDao{}
}

func (s *sysUserDao) InsertUser(db *gorm.DB, user *po.SysUser) error {
	return db.Create(user).Error
}

func (s *sysUserDao) UpdateUser(db *gorm.DB, user *po.SysUser) error {
	return db.Model(user).Updates(user).Error
}

func (s *sysUserDao) DeleteUserByID(db *gorm.DB, id uint) error {
	return db.Delete(&po.SysUser{}, id).Error
}

func (s *sysUserDao) GetUserByID(db *gorm.DB, id uint) (*po.SysUser, error) {
	user := &po.SysUser{}
	err := db.Preload("Roles").First(&user, id).Error
	return user, err
}

func (s *sysUserDao) GetUserNameByID(db *gorm.DB, id uint) (string, error) {
	user := &po.SysUser{}
	err := db.Select("username").First(user, id).Error
	return user.Username, err
}

func (s *sysUserDao) GetUserList(db *gorm.DB, pageQuery *dto.PageQuery) ([]*po.SysUser, error) {
	var user *po.SysUser
	if pageQuery.Query != nil {
		user = pageQuery.Query.(*po.SysUser)
	}
	offset := (pageQuery.Page - 1) * pageQuery.PageSize
	var users []*po.SysUser
	if user != nil && user.Username != "" {
		db = db.Where("username LIKE ?", "%"+user.Username+"%")
	}
	if user != nil && user.Sex != 0 {
		db = db.Where("sex = ?", user.Sex)
	}
	if user != nil && user.LoginName != "" {
		db = db.Where("login_name LIKE ?", "%"+user.LoginName+"%")
	}
	if user != nil && user.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+user.Phone+"%")
	}
	if user != nil && user.Email != "" {
		db = db.Where("email LIKE ?", "%"+user.Email+"%")
	}
	err := db.Limit(pageQuery.PageSize).Offset(offset).Find(&users).Error
	return users, err
}

func (s *sysUserDao) GetUserCount(db *gorm.DB, user *po.SysUser) (int64, error) {
	var count int64
	if user != nil && user.Username != "" {
		db = db.Where("username LIKE ?", "%"+user.Username+"%")
	}
	if user != nil && user.Sex != 0 {
		db = db.Where("sex = ?", user.Sex)
	}
	if user != nil && user.LoginName != "" {
		db = db.Where("login_name LIKE ?", "%"+user.LoginName+"%")
	}
	if user != nil && user.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+user.Phone+"%")
	}
	if user != nil && user.Email != "" {
		db = db.Where("email LIKE ?", "%"+user.Email+"%")
	}
	err := db.Model(&po.SysUser{}).Count(&count).Error
	return count, err
}

func (s *sysUserDao) GetUserByLoginName(db *gorm.DB, loginName string) (*po.SysUser, error) {
	var user po.SysUser
	err := db.Where("login_name = ?", loginName).Preload("Roles").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *sysUserDao) GetUserByEmail(db *gorm.DB, email string) (*po.SysUser, error) {
	var user po.SysUser
	err := db.Where("email = ?", email).Preload("Roles").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *sysUserDao) CheckLoginName(db *gorm.DB, loginname string) (bool, error) {
	var user *po.SysUser
	err := db.Where("login_name = ?", loginname).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *sysUserDao) ListLoginName(db *gorm.DB, loginName string) ([]string, error) {
	var users []*po.SysUser
	err := db.Model(&po.SysUser{}).Where("login_name like ?", loginName+"%").
		Select("login_name").Find(&users).Error
	if err != nil {
		return nil, err
	}
	answer := make([]string, len(users))
	for i := 0; i < len(users); i++ {
		answer[i] = users[i].LoginName
	}
	return answer, nil
}

func (s *sysUserDao) CheckEmailIsExist(db *gorm.DB, email string) (bool, error) {
	var user *po.SysUser
	err := db.Where("email = ? and deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return user.ID != 0, nil
}

func (s *sysUserDao) DeleteUserRoleByUserID(db *gorm.DB, userID uint) error {
	user := po.SysUser{}
	user.ID = userID
	if err := db.Model(&user).Association("Roles").Clear(); err != nil {
		return err
	}
	return nil
}

func (s *sysUserDao) GetRolesByUserID(db *gorm.DB, userID uint) ([]*po.SysRole, error) {
	user := po.SysUser{}
	user.ID = userID
	if err := db.Model(&user).Association("Roles").Find(&user.Roles); err != nil {
		return nil, err
	}
	return user.Roles, nil
}

func (s *sysUserDao) GetRoleIDsByUserID(db *gorm.DB, userID uint) ([]uint, error) {
	user := po.SysUser{}
	user.ID = userID
	if err := db.Model(&user).Association("Roles").Find(&user.Roles); err != nil {
		return nil, err
	}
	roleIDs := make([]uint, len(user.Roles))
	for i, role := range user.Roles {
		roleIDs[i] = role.ID
	}
	return roleIDs, nil
}

func (s *sysUserDao) InsertRolesToUser(db *gorm.DB, userID uint, roleIDs []uint) error {
	user := &po.SysUser{}
	user.ID = userID
	for _, roleID := range roleIDs {
		role := &po.SysRole{}
		role.ID = roleID
		err := db.Model(user).Association("Roles").Append(role)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sysUserDao) UpdateUserPassword(db *gorm.DB, userID uint, password string) error {
	user := po.SysUser{}
	user.ID = userID
	return db.Model(&user).Update("password", password).Error
}
