package system_service

import (
	"context"
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/dao/mock"
	"github.com/fansqz/fancode-backend/models/dto"
	"github.com/fansqz/fancode-backend/models/po"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSysApiService_GetApiCount(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()
	apiDao := mock.NewMockSysApiDao(mockCtl)
	apiDao.EXPECT().GetApiCount(gomock.Any()).Return(int64(10), nil)
	apiDao.EXPECT().GetApiCount(gomock.Any()).Return(int64(0), gorm.ErrInvalidDB)
	apiService := NewSysApiService(apiDao)
	count, err := apiService.GetApiCount(ctx)
	assert.Equal(t, int64(10), count)
	assert.Nil(t, err)
	count, err = apiService.GetApiCount(ctx)
	assert.Equal(t, err, e.ErrMysql)
	assert.Equal(t, int64(0), count)

}

func TestSysApiService_GetApiByID(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()
	// mock apiDao
	apiDao := mock.NewMockSysApiDao(mockCtl)
	api := &po.SysApi{
		Name:        "api名",
		Description: "api描述",
		Path:        "请求路径",
		Method:      "请求方法",
		ParentApiID: 10,
	}
	api.ID = 1
	apiDao.EXPECT().GetApiByID(gomock.Any(), uint(1)).Return(api, nil)
	apiDao.EXPECT().GetApiByID(gomock.Any(), uint(2)).Return(nil, gorm.ErrRecordNotFound)
	apiDao.EXPECT().GetApiByID(gomock.Any(), uint(3)).Return(nil, gorm.ErrInvalidDB)

	// 测试
	apiService := NewSysApiService(apiDao)
	api2, err := apiService.GetApiByID(ctx, 1)
	assert.Equal(t, api, api2)
	assert.Nil(t, err)

	api3, err := apiService.GetApiByID(ctx, 2)
	assert.Nil(t, api3)
	assert.NotNil(t, err)

	api4, err := apiService.GetApiByID(ctx, 3)
	assert.Nil(t, api4)
	assert.Equal(t, err, e.ErrMysql)
}

func TestSysApiService_UpdateApi(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()

	// mock apiDao
	apiDao := mock.NewMockSysApiDao(mockCtl)
	api := &po.SysApi{
		Name:        "api名",
		Description: "api描述",
		Path:        "请求路径",
		Method:      "请求方法",
		ParentApiID: 10,
	}
	api.ID = 1
	apiDao.EXPECT().UpdateApi(gomock.Any(), api).
		DoAndReturn(func(db *gorm.DB, sysApi *po.SysApi) error {
			assert.Equal(t, api, sysApi)
			return nil
		})
	apiDao.EXPECT().UpdateApi(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

	// 测试
	apiService := NewSysApiService(apiDao)
	err := apiService.UpdateApi(ctx, api)
	assert.Nil(t, err)
	err = apiService.UpdateApi(ctx, &po.SysApi{})
	assert.Equal(t, err, e.ErrMysql)
}

func TestSysApiService_InsertApi(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()
	// mock apiDao
	apiDao := mock.NewMockSysApiDao(mockCtl)
	api := &po.SysApi{
		Name:        "api名",
		Description: "api描述",
		Path:        "请求路径",
		Method:      "请求方法",
		ParentApiID: 10,
	}
	apiDao.EXPECT().InsertApi(gomock.Any(), api).
		DoAndReturn(func(db *gorm.DB, sysApi *po.SysApi) error {
			assert.Equal(t, api, sysApi)
			sysApi.ID = 1
			return nil
		})
	apiDao.EXPECT().InsertApi(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidDB)

	// 测试
	apiService := NewSysApiService(apiDao)
	id, err := apiService.InsertApi(ctx, api)
	assert.Equal(t, id, uint(1))
	assert.Nil(t, err)
	id, err = apiService.InsertApi(ctx, &po.SysApi{})
	assert.Equal(t, id, uint(0))
	assert.Equal(t, err, e.ErrMysql)
}

func TestSysApiService_GetApiTree(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	ctx := context.Background()
	// mock apiDao
	apiDao := mock.NewMockSysApiDao(mockCtl)
	apis := make([]*po.SysApi, 4)
	for i := 0; i < 4; i++ {
		api := &po.SysApi{}
		api.Name = "api" + strconv.Itoa(i)
		api.Description = "api描述" + strconv.Itoa(i)
		api.Path = "apiPath" + strconv.Itoa(i)
		api.Method = "apiMethod" + strconv.Itoa(i)
		api.ID = uint(i + 1)
		apis[i] = api
	}
	apis[1].ParentApiID = 1
	apis[2].ParentApiID = 1
	apis[3].ParentApiID = 1
	apiDao.EXPECT().GetAllApi(gomock.Any()).Return(apis, nil)
	apiDao.EXPECT().GetAllApi(gomock.Any()).Return([]*po.SysApi{}, gorm.ErrInvalidDB)

	// 测试
	apiService := NewSysApiService(apiDao)

	treeDtos, err := apiService.GetApiTree(ctx)
	treeDto := dto.NewSysApiTreeDto(apis[0])
	for i := 1; i < 4; i++ {
		treeDto.Children = append(treeDto.Children, dto.NewSysApiTreeDto(apis[i]))
	}
	assert.Equal(t, []*dto.SysApiTreeDto{treeDto}, treeDtos)
	assert.Nil(t, err)

	treeDtos, err = apiService.GetApiTree(ctx)
	assert.Nil(t, treeDtos)
	assert.Equal(t, err, e.ErrMysql)
}
