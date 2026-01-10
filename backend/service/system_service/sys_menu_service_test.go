package system_service

import (
	"context"
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/dao/mock"
	"github.com/fansqz/fancode-backend/models/dto"
	"github.com/fansqz/fancode-backend/models/po"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"strconv"
	"testing"
)

func TestSysMenuService_GetMenuCount(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()
	menuDao := mock.NewMockSysMenuDao(mockCtl)
	menuDao.EXPECT().GetMenuCount(gomock.Any()).Return(int64(10), nil)
	menuDao.EXPECT().GetMenuCount(gomock.Any()).Return(int64(0), gorm.ErrInvalidDB)
	menuService := NewSysMenuService(menuDao)
	count, err := menuService.GetMenuCount(ctx)
	assert.Equal(t, int64(10), count)
	assert.Nil(t, err)
	count, err = menuService.GetMenuCount(ctx)
	assert.Equal(t, err, e.ErrMysql)
	assert.Equal(t, int64(0), count)

}

func TestSysMenuService_GetMenuByID(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()
	// mock apiDao
	menuDao := mock.NewMockSysMenuDao(mockCtl)
	menu := &po.SysMenu{
		Name:         "menu名",
		Description:  "menu描述",
		Code:         "menu权限值",
		ParentMenuID: 10,
	}
	menu.ID = 1
	menuDao.EXPECT().GetMenuByID(gomock.Any(), uint(1)).Return(menu, nil)
	menuDao.EXPECT().GetMenuByID(gomock.Any(), uint(2)).Return(nil, gorm.ErrRecordNotFound)
	menuDao.EXPECT().GetMenuByID(gomock.Any(), uint(3)).Return(nil, gorm.ErrInvalidDB)

	// 测试
	menuService := NewSysMenuService(menuDao)
	menu2, err := menuService.GetMenuByID(ctx, 1)
	assert.Equal(t, menu, menu2)
	assert.Nil(t, err)

	menu3, err := menuService.GetMenuByID(ctx, 2)
	assert.Nil(t, menu3)
	assert.NotNil(t, err)

	menu4, err := menuService.GetMenuByID(ctx, 3)
	assert.Nil(t, menu4)
	assert.Equal(t, err, e.ErrMysql)
}

func TestSysMenuService_UpdateMenu(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()
	// mock apiDao
	menuDao := mock.NewMockSysMenuDao(mockCtl)
	menu := &po.SysMenu{
		Name:         "menu名",
		Description:  "menu描述",
		Code:         "menuCode",
		ParentMenuID: 10,
	}
	menu.ID = 1
	menuDao.EXPECT().UpdateMenu(gomock.Any(), menu).
		DoAndReturn(func(db *gorm.DB, sysMenu *po.SysMenu) error {
			assert.Equal(t, menu, sysMenu)
			return nil
		})
	menuDao.EXPECT().UpdateMenu(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

	// 测试
	menuService := NewSysMenuService(menuDao)
	err := menuService.UpdateMenu(ctx, menu)
	assert.Nil(t, err)
	err = menuService.UpdateMenu(ctx, &po.SysMenu{})
	assert.Equal(t, err, e.ErrMysql)
}

func TestSysMenuService_InsertMenu(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()
	// mock apiDao
	menuDao := mock.NewMockSysMenuDao(mockCtl)
	menu := &po.SysMenu{
		Name:         "menu名",
		Description:  "menu描述",
		Code:         "menuCode",
		ParentMenuID: 10,
	}
	menuDao.EXPECT().InsertMenu(gomock.Any(), menu).
		DoAndReturn(func(db *gorm.DB, sysMenu *po.SysMenu) error {
			assert.Equal(t, menu, sysMenu)
			sysMenu.ID = 1
			return nil
		})
	menuDao.EXPECT().InsertMenu(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

	// 测试
	menuService := NewSysMenuService(menuDao)
	id, err := menuService.InsertMenu(ctx, menu)
	assert.Equal(t, id, uint(1))
	assert.Nil(t, err)
	id, err = menuService.InsertMenu(ctx, &po.SysMenu{})
	assert.Equal(t, id, uint(0))
	assert.Equal(t, err, e.ErrMysql)
}

func TestSysMenuService_GetMenuTree(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()
	// mock apiDao
	menuDao := mock.NewMockSysMenuDao(mockCtl)
	menus := make([]*po.SysMenu, 4)
	for i := 0; i < 4; i++ {
		menu := &po.SysMenu{}
		menu.Name = "menu" + strconv.Itoa(i)
		menu.Description = "menu描述" + strconv.Itoa(i)
		menu.Code = "menuPath" + strconv.Itoa(i)
		menu.ID = uint(i + 1)
		menus[i] = menu
	}
	menus[1].ParentMenuID = 1
	menus[2].ParentMenuID = 1
	menus[3].ParentMenuID = 1
	menuDao.EXPECT().GetAllMenu(gomock.Any()).Return(menus, nil)
	menuDao.EXPECT().GetAllMenu(gomock.Any()).Return([]*po.SysMenu{}, gorm.ErrInvalidDB)

	// 测试
	menuService := NewSysMenuService(menuDao)

	treeDtos, err := menuService.GetMenuTree(ctx)
	treeDto := dto.NewSysMenuTreeDto(menus[0])
	for i := 1; i < 4; i++ {
		treeDto.Children = append(treeDto.Children, dto.NewSysMenuTreeDto(menus[i]))
	}
	assert.Equal(t, []*dto.SysMenuTreeDto{treeDto}, treeDtos)
	assert.Nil(t, err)

	treeDtos, err = menuService.GetMenuTree(ctx)
	assert.Nil(t, treeDtos)
	assert.Equal(t, err, e.ErrMysql)
}
