package meta

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strconv"
)

func Param2Int(key string, c *gin.Context) (int64, error) {
	id := c.Param(key)
	if id == "" {
		return 0, errors.Errorf("%s is error", key)
	}
	idInt, err := strconv.ParseInt(id, 10, 64)
	return idInt, err
}
