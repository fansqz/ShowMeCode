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

type SysMenuService interface {

	// GetMenuCount 获取menu数目
	GetMenuCount(ctx context.Context) (int64, error)
	// DeleteMenuByID 删除menu
	DeleteMenuByID(ctx context.Context, id uint) error
	// UpdateMenu 更新menu
	UpdateMenu(ctx context.Context, menu *po.SysMenu) error
	// GetMenuByID 根据id获取menu
	GetMenuByID(ctx context.Context, id uint) (*po.SysMenu, error)
	// GetMenuTree 获取menu树
	GetMenuTree(ctx context.Context) ([]*dto.SysMenuTreeDto, error)
	// InsertMenu 添加menu
	InsertMenu(ctx context.Context, menu *po.SysMenu) (uint, error)
}

type sysMenuService struct {
	sysMenuDao dao.SysMenuDao
}

func NewSysMenuService(menuDao dao.SysMenuDao) SysMenuService {
	return &sysMenuService{
		sysMenuDao: menuDao,
	}
}

func (s *sysMenuService) GetMenuCount(ctx context.Context) (int64, error) {
	count, err := s.sysMenuDao.GetMenuCount(common.Mysql)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[GetMenuCount] get menu count error, err = %v", err)
		return 0, e.ErrMysql
	}
	return count, nil
}

// DeleteMenuByID 根据menu的id进行删除
func (s *sysMenuService) DeleteMenuByID(ctx context.Context, id uint) error {
	err := common.Mysql.Transaction(func(tx *gorm.DB) error {
		// 递归删除API
		return s.deleteMenusRecursive(ctx, tx, id)
	})

	if err != nil {
		logger.WithCtx(ctx).Errorf("[DeleteMenuByID] get menu by id error, err = %v", err)
		return e.ErrMysql
	}

	return nil
}

// deleteMenusRecursive 递归删除API
func (s *sysMenuService) deleteMenusRecursive(ctx context.Context, db *gorm.DB, parentID uint) error {
	childMenus, err := s.sysMenuDao.GetChildMenusByParentID(db, parentID)
	if err != nil {
		return err
	}
	for _, childAPI := range childMenus {
		// 删除子menu的子menu
		if err = s.deleteMenusRecursive(ctx, db, childAPI.ID); err != nil {
			return err
		}
	}
	// 当前menu
	if err = s.sysMenuDao.DeleteMenuByID(db, parentID); err != nil {
		logger.WithCtx(ctx).Errorf("[deleteMenusRecursive] delete menu by id error, err = %v", err)
		return err
	}
	return nil
}

func (s *sysMenuService) UpdateMenu(ctx context.Context, menu *po.SysMenu) error {
	err := s.sysMenuDao.UpdateMenu(common.Mysql, menu)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[UpdateMenu] update menu error, err = %v", err)
		return e.ErrMysql
	}
	return nil
}

func (s *sysMenuService) GetMenuByID(ctx context.Context, id uint) (*po.SysMenu, error) {
	menu, err := s.sysMenuDao.GetMenuByID(common.Mysql, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, e.NewCustomMsg("The Menu is not exist")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.WithCtx(ctx).Errorf("[GetMenuByID] get menu error, err = %v", err)
		return nil, e.ErrMysql
	}
	return menu, nil
}

func (s *sysMenuService) GetMenuTree(ctx context.Context) ([]*dto.SysMenuTreeDto, error) {
	var menuList []*po.SysMenu
	var err error
	if menuList, err = s.sysMenuDao.GetAllMenu(common.Mysql); err != nil {
		logger.WithCtx(ctx).Errorf("[GetMenuTree] get all menu error, err = %v", err)
		return nil, e.ErrMysql
	}

	menuMap := make(map[uint]*dto.SysMenuTreeDto)
	var rootMenus []*dto.SysMenuTreeDto

	// 添加到map中保存
	for _, menu := range menuList {
		menuMap[menu.ID] = dto.NewSysMenuTreeDto(menu)
	}

	// 遍历并添加到父节点中
	for _, menu := range menuList {
		if menu.ParentMenuID == 0 {
			rootMenus = append(rootMenus, menuMap[menu.ID])
		} else {
			parentMenu, exists := menuMap[menu.ParentMenuID]
			if !exists {
				return nil, e.NewCustomMsg("The Menu is not exist")
			}
			parentMenu.Children = append(parentMenu.Children, menuMap[menu.ID])
		}
	}

	return rootMenus, nil
}

func (s *sysMenuService) InsertMenu(ctx context.Context, menu *po.SysMenu) (uint, error) {
	err := s.sysMenuDao.InsertMenu(common.Mysql, menu)
	if err != nil {
		logger.WithCtx(ctx).Errorf("[InsertMenu] insert menu error, err = %v", err)
		return 0, e.ErrMysql
	}
	return menu.ID, nil
}
