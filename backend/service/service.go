package service

import (
	"github.com/fansqz/fancode-backend/service/common_service"
	"github.com/fansqz/fancode-backend/service/system_service"
	"github.com/fansqz/fancode-backend/service/user_coding_service"
	"github.com/fansqz/fancode-backend/service/user_saved_code_service"
	"github.com/fansqz/fancode-backend/service/visual_debug_servcie"
	"github.com/fansqz/fancode-backend/service/visual_document_service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	common_service.NewAccountService,
	common_service.NewAuthService,
	common_service.NewCommonService,
	system_service.NewSysApiService,
	system_service.NewSysMenuService,
	system_service.NewSysRoleService,
	system_service.NewSysUserService,
	user_coding_service.NewUserCodeService,
	user_saved_code_service.NewUserSavedCodeService,
	visual_debug_servcie.NewDebugService,
	visual_debug_servcie.NewVisualService,
	visual_document_service.NewVisualDocumentService,
	visual_document_service.NewVisualDocumentBankService,
)
