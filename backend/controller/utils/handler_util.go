package utils

import (
	e "github.com/fansqz/fancode-backend/common/error"
	"github.com/fansqz/fancode-backend/models/dto"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetIntQueryOrDefault(ctx *gin.Context, key string, i int) int {
	a := ctx.Query(key)
	return AtoiOrDefault(a, i)
}

func GetIntParamOrDefault(ctx *gin.Context, key string, i int) int {
	a := ctx.Param(key)
	return AtoiOrDefault(a, i)
}

func GetBoolQuery(ctx *gin.Context, key string) bool {
	a := ctx.Query(key)
	return Atob(a)
}

func GetBoolParam(ctx *gin.Context, key string) bool {
	a := ctx.Param(key)
	return Atob(a)
}

func Atob(a string) bool {
	if a == "false" || a == "0" {
		return false
	} else {
		return true
	}
}

func AtoiOrDefault(a string, i int) int {
	if a == "" {
		return i
	}
	b, err := strconv.Atoi(a)
	if err != nil {
		return i
	} else {
		return b
	}
}
func GetPageQueryByQuery(ctx *gin.Context) (*dto.PageQuery, error) {
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("pageSize")
	var page int
	var pageSize int
	var convertErr error
	page, convertErr = strconv.Atoi(pageStr)
	if convertErr != nil {
		return nil, e.ErrBadRequest
	}
	pageSize, convertErr = strconv.Atoi(pageSizeStr)
	if convertErr != nil {
		return nil, e.ErrBadRequest
	}
	if pageSize > 50 {
		pageSize = 50
	}
	sortProperty := ctx.Query("sortProperty")
	sortRule := ctx.Query("sortRule")
	answer := &dto.PageQuery{
		Page:         page,
		PageSize:     pageSize,
		SortProperty: sortProperty,
		SortRule:     sortRule,
	}
	return answer, nil
}
