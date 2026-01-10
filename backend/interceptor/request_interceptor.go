package interceptor

import (
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/common/logger"
	"github.com/fansqz/fancode-backend/constants"
	r "github.com/fansqz/fancode-backend/models/vo"
	"github.com/fansqz/fancode-backend/service/system_service"
	"github.com/fansqz/fancode-backend/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

type RequestInterceptor struct {
	roleService system_service.SysRoleService
	userService system_service.SysUserService
}

func NewRequestInterceptor(roleService system_service.SysRoleService, userService system_service.SysUserService) *RequestInterceptor {
	return &RequestInterceptor{
		roleService: roleService,
		userService: userService,
	}
}

// TokenAuthorize
//
//	@Description: token拦截器
//	@return gin.HandlerFunc
func (i *RequestInterceptor) TokenAuthorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.WithCtx(c).Infof("[requst] url: %s", c.Request.URL.Path)
		result := r.NewResult(c)

		// 从token中获取用户信息
		token := c.Request.Header.Get("token")
		var userID uint
		if token != "" {
			userID = i.getAndSaveUserID(c, token)
		}

		// 检验路径是否在游客路径中
		allow, _ := i.checkIsAllowTouristReq(c)
		if allow {
			i.checkAndUpdateToken(c, userID)
			c.Next()
			return
		}

		// 不在游客放行路径内，且没有登陆，直接拦截
		if userID == 0 {
			logger.WithCtx(c).Warnf("[TokenAuthorize] userInfo is nil")
			result.Error(e.ErrSessionExpire)
			c.Abort()
			return
		}

		// 进行权限校验
		allow, err := i.checkIsRoleAllowReq(c, userID)
		if err != nil {
			result.Error(err)
			c.Abort()
			return
		}
		if allow {
			i.checkAndUpdateToken(c, userID)
			c.Next()
			return
		}
		rejectRequest(c)
	}
}

// checkAndUpdateToken 校验token时间，如果要过期则更新token
func (i *RequestInterceptor) checkAndUpdateToken(ctx *gin.Context, userID uint) {
	if userID == 0 {
		return
	}
	token := ctx.Request.Header.Get("token")
	ok, err := utils.ShouldUpdateToken(token)
	if err == nil && ok {
		// 续签token
		newToken, err := utils.GenerateToken(utils.Claims{
			ID: userID,
		})
		if err == nil {
			ctx.Header(constants.TokenHeader, newToken)
		}
	}
}

// checkIsRoleAllowReq 校验是否是用户角色
func (i *RequestInterceptor) checkIsRoleAllowReq(ctx *gin.Context, userID uint) (bool, error) {
	rules, err := i.userService.GetRoleIDsByUserID(ctx, userID)
	if err != nil {
		logger.WithCtx(ctx).Warnf("[TokenAuthorize] get apis error, err = %v", err)
		return false, err
	}
	path := ctx.Request.URL.Path
	method := ctx.Request.Method
	for _, ruleID := range rules {
		apis, _ := i.roleService.GetApisByRoleID(ctx, ruleID)
		for _, api := range apis {
			if matchPath(path, api.Path) {
				if strings.EqualFold(method, constants.AllMethod) {
					return true, nil
				} else if strings.EqualFold(method, api.Method) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

// checkIsAllowTouristReq 检查路径是否在放行游客放行的白名单中
func (i *RequestInterceptor) checkIsAllowTouristReq(ctx *gin.Context) (bool, error) {
	path := ctx.Request.URL.Path
	apis, err := i.roleService.GetApisByRoleID(ctx, constants.TouristID)
	if err != nil {
		logger.WithCtx(ctx).Warnf("[TokenAuthorize] get apis error, err = %v", err)
		return false, err
	}
	logger.WithCtx(ctx).Infof("[TokenAuthorize] api auth success")
	method := ctx.Request.Method
	for _, api := range apis {
		if matchPath(path, api.Path) {
			if strings.EqualFold(method, constants.AllMethod) {
				return true, nil
			} else if strings.EqualFold(method, api.Method) {
				return true, nil
			}
		}
	}
	return false, nil
}

// getUserAndSaveInfo 从token中读取用户信息并保存到ctx中
func (i *RequestInterceptor) getAndSaveUserID(ctx *gin.Context, token string) uint {
	claims, err := utils.ParseToken(token)
	var userID uint
	if err != nil {
		logger.WithCtx(ctx).Warnf("[TokenAuthorize] token authorize err, err = %v", err)
	} else {
		userID = claims.ID
		if ctx.Keys == nil {
			ctx.Keys = make(map[string]interface{}, 1)
		}
		ctx.Set(utils.CtxUserIDKey, userID)
		// 记录 用户id
		ctx.Set(logger.USER_ID_KEY, userID)
	}
	return userID
}

// matchPath 判断请求路径是否和规则相匹配
func matchPath(requestPath, pattern string) bool {
	routeSegments := strings.Split(requestPath, "/")
	patternSegments := strings.Split(pattern, "/")

	if len(routeSegments) != len(patternSegments) {
		return false
	}

	for i := 0; i < len(routeSegments); i++ {
		if patternSegments[i] != "" && patternSegments[i] != routeSegments[i] {
			if !strings.HasPrefix(patternSegments[i], ":") {
				return false
			}
		}
	}

	return true
}

func rejectRequest(ctx *gin.Context) {
	result := r.NewResult(ctx)
	result.Error(e.ErrPermissionInvalid)
	ctx.Abort()
}
