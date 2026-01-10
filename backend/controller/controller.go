package controller

import (
	"github.com/fansqz/fancode-backend/controller/admin"
	"github.com/fansqz/fancode-backend/controller/user"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewAccountController,
	NewAuthController,
	NewFileController,
	user.NewUserSavedCodeHandler,
	admin.NewSysApiController,
	admin.NewSysMenuController,
	admin.NewSysRoleController,
	admin.NewSysUserController,
	admin.NewVisualDocumentManageController,
	admin.NewVisualDocumentBankManageController,
	user.NewDebugController,
	user.NewVisualController,
	user.NewVisualDocumentController,
	user.NewVisualDocumentBankController,
	NewCommonController,
)
