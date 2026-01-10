package dao

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewSysApiDao,
	NewSysMenuDao,
	NewSysRoleDao,
	NewSysUserDao,
	NewUserCodeDao,
	NewUserSavedCodeDao,
	NewVisualDocumentDao,
	NewVisualDocumentBankDao,
)
