package middleware

import (
	"elevate-hub/config"
	"elevate-hub/db/meta"
	"elevate-hub/db/models"
	jwt "elevate-hub/middleware/jwt"
	"elevate-hub/svc"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// AuthInit jwt验证new
func AuthInit(conf config.Config, ctx *svc.ServiceContext) (*jwt.GinJWTMiddleware, error) {
	timeout := time.Hour

	if conf.JWT.Timeout != 0 {
		timeout = time.Duration(conf.JWT.Timeout) * time.Hour
	}
	return jwt.New(&jwt.GinJWTMiddleware{
		SvCtx:           ctx,
		Realm:           "elevate-hub",
		Key:             []byte(conf.JWT.Secret),
		Timeout:         timeout,
		MaxRefresh:      time.Hour,
		SendCookie:      true,
		PayloadFunc:     PayloadFunc,     // success
		IdentityHandler: IdentityHandler, // success
		Authenticator:   Authenticator,   // success
		Authorizator:    Authorizator,
		Unauthorized:    Unauthorized, // success
		//TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenLookup:   "header: Authorization, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

}

func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		user := v["user"].(models.User)
		return jwt.MapClaims{
			jwt.IdentityKey:   user.ID,
			jwt.RoleIdKey:     v["role"],
			jwt.NiceKey:       user.NickName,
			jwt.DeptId:        user.DeptId,
			jwt.Username:      user.Username,
			jwt.SupperManager: user.SupperManager,
		}
	}
	return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return map[string]interface{}{
		"IdentityKey":  claims[jwt.IdentityKey],
		"UserName":     claims[jwt.Username],
		"UserId":       claims[jwt.IdentityKey],
		"RoleIds":      claims[jwt.RoleIdKey],
		"SupperManage": claims[jwt.SupperManager],
		//"DataScope":   claims["datascope"],
	}
}

func Authenticator(c *gin.Context, svCtx *svc.ServiceContext) (interface{}, error) {
	var userInfo struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := c.ShouldBind(&userInfo)
	if err != nil {
		log.Printf("login failed: %v", err)
		return nil, errors.Errorf("登录需要username, password字段请确认。")
	}

	userInfo.Password = jwt.MD5Password(strings.TrimSpace(userInfo.Password))

	user, err := svCtx.Mysql.UserModel.Get(meta.GetOption{
		Condition: map[string]interface{}{
			"username": strings.TrimSpace(userInfo.Username),
			"password": userInfo.Password,
		},
	})
	if err != nil {
		return nil, errors.New("账号或密码错误,请重新登录")
	}
	roles, err := svCtx.Mysql.AuthRelationModel.GetUserRoles(user.ID)
	if err != nil {
		return nil, errors.New("系统出现了些小问题,请稍后重试")
	}
	return map[string]interface{}{
		"user": *user,
		"role": roles,
	}, nil

	//return nil, jwt.ErrFailedAuthentication
}

// 过期 啥啥啥的校验
func Authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(map[string]interface{}); ok {
		c.Set("roleIds", v["RoleIds"])
		c.Set("userId", v["UserId"])
		c.Set("userName", v["Username"])
		c.Set("supperManage", v["SupperManage"])
		return true
	}
	return true
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"code": code,
		"msg":  message,
	})
}
