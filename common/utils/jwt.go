package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func GetUserID(c *gin.Context) (int64, error) {
	userId, ok := c.Get("userId")
	if !ok {
		return 0, errors.New("user id is not exist")
	}
	number, ok := userId.(json.Number)
	if !ok {
		return 0, errors.New("user id is not json number")
	}
	return number.Int64()
}

/*
*
roleIds. isAdmin, error
*/
func GetRoles(c *gin.Context) ([]int64, bool, error) {

	isadmin := isSupperManage(c)

	roleIds, ok := c.Get("roleIds")
	if !ok {
		return []int64{}, isadmin, errors.New("roleIds is not exist")
	}
	ids, ok := roleIds.([]interface{})
	if !ok {
		return []int64{}, isadmin, errors.New("roleIds is not json array")
	}
	var intIds = make([]int64, 0, len(ids))
	for _, roleId := range ids {
		id, _ := roleId.(json.Number).Int64()
		intIds = append(intIds, id)
	}
	return intIds, isadmin, nil

}

func isSupperManage(c *gin.Context) bool {
	supperManage, ok := c.Get("supperManage")
	if !ok {
		return false
	}
	number, _ := supperManage.(json.Number).Int64()
	return number == 1

}
