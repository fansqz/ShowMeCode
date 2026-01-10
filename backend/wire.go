//go:build wireinject
// +build wireinject

package main

import (
	"github.com/fansqz/fancode-backend/common/config"
	"github.com/fansqz/fancode-backend/controller"
	"github.com/fansqz/fancode-backend/dao"
	"github.com/fansqz/fancode-backend/interceptor"
	"github.com/fansqz/fancode-backend/routers"
	"github.com/fansqz/fancode-backend/service"
	"github.com/google/wire"
	"net/http"
)

func initApp(*config.AppConfig) (*http.Server, error) {
	panic(wire.Build(
		dao.ProviderSet,
		service.ProviderSet,
		controller.ProviderSet,
		interceptor.ProviderSet,
		routers.SetupRouter,
		newApp))
}
