package system_service

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fansqz/fancode-backend/common"
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/dao/mock"
	"github.com/fansqz/fancode-backend/models/dto"
	"github.com/fansqz/fancode-backend/models/po"
	"github.com/fansqz/fancode-backend/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestSysRoleService_GetRoleByID(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()

	// mock roleDao
	roleDao := mock.NewMockSysRoleDao(mockCtl)
	role := &po.SysRole{
		Name:        "role名",
		Description: "role描述",
	}
	role.ID = 1
	roleDao.EXPECT().GetRoleByID(gomock.Any(), uint(1)).Return(role, nil)
	roleDao.EXPECT().GetRoleByID(gomock.Any(), uint(2)).Return(nil, gorm.ErrRecordNotFound)
	roleDao.EXPECT().GetRoleByID(gomock.Any(), uint(3)).Return(nil, gorm.ErrInvalidDB)

	// 测试
	roleService := NewSysRoleService(roleDao)
	role2, err := roleService.GetRoleByID(ctx, 1)
	assert.Equal(t, role, role2)
	assert.Nil(t, err)

	role3, err := roleService.GetRoleByID(ctx, 2)
	assert.Nil(t, role3)
	assert.NotNil(t, err)

	role4, err := roleService.GetRoleByID(ctx, 3)
	assert.Nil(t, role4)
	assert.Equal(t, err, e.ErrMysql)
}

func TestSysRoleService_InsertSysRole(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()

	// mock roleDao
	roleDao := mock.NewMockSysRoleDao(mockCtl)
	role := &po.SysRole{
		Name:        "menu名",
		Description: "menu描述",
	}
	roleDao.EXPECT().InsertRole(gomock.Any(), role).
		DoAndReturn(func(db *gorm.DB, sysRole *po.SysRole) error {
			assert.Equal(t, role, sysRole)
			sysRole.ID = 1
			return nil
		})
	roleDao.EXPECT().InsertRole(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

	// 测试
	roleService := NewSysRoleService(roleDao)
	id, err := roleService.InsertSysRole(ctx, role)
	assert.Equal(t, id, uint(1))
	assert.Nil(t, err)
	id, err = roleService.InsertSysRole(ctx, &po.SysRole{})
	assert.Equal(t, id, uint(0))
	assert.Equal(t, err, e.ErrMysql)
}

func TestSysRoleService_UpdateSysRole(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()

	// mock roleDao
	roleDao := mock.NewMockSysRoleDao(mockCtl)
	role := &po.SysRole{
		Name:        "role名",
		Description: "role描述",
	}
	role.ID = 1
	roleDao.EXPECT().UpdateRole(gomock.Any(), role).
		DoAndReturn(func(db *gorm.DB, sysRole *po.SysRole) error {
			assert.Equal(t, role, sysRole)
			return nil
		})
	roleDao.EXPECT().UpdateRole(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

	// 测试
	roleService := NewSysRoleService(roleDao)
	err := roleService.UpdateSysRole(ctx, role)
	assert.Nil(t, err)
	err = roleService.UpdateSysRole(ctx, &po.SysRole{})
	assert.Equal(t, err, e.ErrMysql)
}

func TestSysRoleService_DeleteSysRole(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()

	// mock roleDao
	roleDao := mock.NewMockSysRoleDao(mockCtl)
	roleDao.EXPECT().DeleteRoleByID(gomock.Any(), uint(1)).Return(nil)
	roleDao.EXPECT().DeleteRoleByID(gomock.Any(), uint(2)).Return(gorm.ErrInvalidDB)

	// 测试
	roleService := NewSysRoleService(roleDao)
	err := roleService.DeleteSysRole(ctx, 1)
	assert.Nil(t, err)
	err = roleService.DeleteSysRole(ctx, 2)
	assert.Equal(t, err, e.ErrMysql)
}

func TestSysRoleService_GetSysRoleList(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()

	roleDao := mock.NewMockSysRoleDao(mockCtl)

	// 测试1
	testRoleList := []*po.SysRole{
		{Name: "role1", Description: "description1", Model: gorm.Model{ID: 0, UpdatedAt: time.Now()}},
		{Name: "role2", Description: "description2", Model: gorm.Model{ID: 1, UpdatedAt: time.Now()}},
		{Name: "role3", Description: "description3", Model: gorm.Model{ID: 2, UpdatedAt: time.Now()}},
	}
	roleDao.EXPECT().GetRoleList(gomock.Any(), &dto.PageQuery{
		Page: 1, PageSize: 3, SortProperty: "name", SortRule: "desc", Query: &po.SysRole{Name: "role", Description: "description"},
	}).Return(testRoleList, nil)
	roleDao.EXPECT().GetRoleCount(gomock.Any(), &po.SysRole{Name: "role", Description: "description"}).Return(int64(10), nil)
	roleService := NewSysRoleService(roleDao)
	roleList, err := roleService.GetSysRoleList(ctx, &dto.PageQuery{
		Page: 1, PageSize: 3, SortProperty: "name", SortRule: "desc", Query: &po.SysRole{Name: "role", Description: "description"},
	})
	assert.Nil(t, err)
	assert.Equal(t, &dto.PageInfo{Total: 10, Size: 3, List: []*dto.SysRoleDtoForList{
		{Name: "role1", Description: "description1", ID: 0, UpdatedAt: utils.Time(testRoleList[0].UpdatedAt)},
		{Name: "role2", Description: "description2", ID: 1, UpdatedAt: utils.Time(testRoleList[1].UpdatedAt)},
		{Name: "role3", Description: "description3", ID: 2, UpdatedAt: utils.Time(testRoleList[2].UpdatedAt)},
	}}, roleList)

	// 测试2
	roleDao.EXPECT().GetRoleList(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrInvalidDB)
	roleList, err = roleService.GetSysRoleList(ctx, &dto.PageQuery{})
	assert.Nil(t, roleList)
	assert.Equal(t, err, e.ErrMysql)

	// 测试3
	roleDao.EXPECT().GetRoleList(gomock.Any(), gomock.Any()).Return([]*po.SysRole{}, nil)
	roleDao.EXPECT().GetRoleCount(gomock.Any(), gomock.Any()).Return(int64(0), gorm.ErrInvalidDB)
	roleList, err = roleService.GetSysRoleList(ctx, &dto.PageQuery{})
	assert.Nil(t, roleList)
	assert.Equal(t, err, e.ErrMysql)

}

func TestSysRoleService_GetApiIDsByRoleID(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()

	roleDao := mock.NewMockSysRoleDao(mockCtl)

	// mock数据库
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{})
	assert.Nil(t, err)
	common.Mysql = gormDB

	mock.ExpectBegin()
	mock.ExpectCommit()
	roleDao.EXPECT().DeleteRoleAPIsByRoleID(gomock.Any(), uint(1)).Return(nil)
	roleDao.EXPECT().InsertApisToRole(gomock.Any(), uint(1), []uint{1, 2, 3}).Return(nil)

	mock.ExpectBegin()
	mock.ExpectRollback()
	roleDao.EXPECT().DeleteRoleAPIsByRoleID(gomock.Any(), uint(2)).Return(gorm.ErrInvalidDB)

	mock.ExpectBegin()
	mock.ExpectRollback()
	roleDao.EXPECT().DeleteRoleAPIsByRoleID(gomock.Any(), uint(3)).Return(nil)
	roleDao.EXPECT().InsertApisToRole(gomock.Any(), uint(3), []uint{1, 2, 3}).Return(gorm.ErrInvalidDB)

	roleService := NewSysRoleService(roleDao)
	err2 := roleService.UpdateRoleApis(ctx, 1, []uint{1, 2, 3})
	assert.Nil(t, err2)

	err2 = roleService.UpdateRoleApis(ctx, 2, []uint{1, 2, 3})
	assert.Equal(t, err2, e.ErrMysql)

	err2 = roleService.UpdateRoleApis(ctx, 3, []uint{1, 2, 3})
	assert.Equal(t, err2, e.ErrMysql)
}

func TestSysRoleService_GetMenuIDsByRoleID(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()

	roleDao := mock.NewMockSysRoleDao(mockCtl)

	// mock数据库
	db, mock, err := sqlmock.New()
	defer db.Close()
	assert.Nil(t, err)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{})
	assert.Nil(t, err)
	common.Mysql = gormDB

	mock.ExpectBegin()
	mock.ExpectCommit()
	roleDao.EXPECT().DeleteRoleMenusByRoleID(gomock.Any(), uint(1)).Return(nil)
	roleDao.EXPECT().InsertMenusToRole(gomock.Any(), uint(1), []uint{1, 2, 3}).Return(nil)

	mock.ExpectBegin()
	mock.ExpectRollback()
	roleDao.EXPECT().DeleteRoleMenusByRoleID(gomock.Any(), uint(2)).Return(gorm.ErrInvalidDB)

	mock.ExpectBegin()
	mock.ExpectRollback()
	roleDao.EXPECT().DeleteRoleMenusByRoleID(gomock.Any(), uint(3)).Return(nil)
	roleDao.EXPECT().InsertMenusToRole(gomock.Any(), uint(3), []uint{1, 2, 3}).Return(gorm.ErrInvalidDB)

	roleService := NewSysRoleService(roleDao)
	err2 := roleService.UpdateRoleMenus(ctx, 1, []uint{1, 2, 3})
	assert.Nil(t, err2)

	err2 = roleService.UpdateRoleMenus(ctx, 2, []uint{1, 2, 3})
	assert.Equal(t, err2, e.ErrMysql)

	err2 = roleService.UpdateRoleMenus(ctx, 3, []uint{1, 2, 3})
	assert.Equal(t, err2, e.ErrMysql)
}
