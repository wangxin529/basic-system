package application

import (
	"elevate-hub/common/set"
	"elevate-hub/common/utils"
	"elevate-hub/db/meta"
	"elevate-hub/db/models"
	"elevate-hub/middleware/constant"
	"elevate-hub/svc"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type ApplicationMiddleware struct {
	Apps    map[string]*models.APP
	APPAPIs map[int64]set.Set[string]
	ctx     *svc.ServiceContext
}

func NewApplicationMiddleware(ctx *svc.ServiceContext) *ApplicationMiddleware {
	am := ApplicationMiddleware{
		ctx: ctx,
	}
	// 初始化数据
	am.refresh()
	go am.cronJob(1 * time.Minute)

	return &am
}

func (a *ApplicationMiddleware) cronJob(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		a.refresh()
	}
}

func (a *ApplicationMiddleware) refresh() {
	var am = ApplicationMiddleware{
		Apps:    make(map[string]*models.APP),
		APPAPIs: make(map[int64]set.Set[string]),
	}
	apps, _, err := a.ctx.Mysql.APPModel.List(meta.EmptyListOption)
	if err != nil {
		log.Println(errors.Wrap(err, "mysql list app failed"))
	}
	for _, app := range apps {
		am.Apps[app.Key] = app
	}

	// 筛选可用API
	var apiSet = make(map[int64]string)
	apis, _, err := a.ctx.Mysql.APIModel.List(meta.EmptyListOption)
	if err != nil {
		log.Println(errors.Wrap(err, "mysql list api failed"))
		return
	}
	for _, api := range apis {
		if api.Status != 0 && api.APIType != 1 {
			continue
		}
		apiSet[api.ID] = api.Path
	}
	// APP 授权API绑定
	appApis, _, err := a.ctx.Mysql.AuthRelationModel.List(meta.ListOption{
		DisableCount: true,
		Condition: map[string]interface{}{
			"role_id": 0,
			"user_id": 0,
			"menu_id": 0,
		},
	})
	if err != nil {
		log.Println(errors.Wrap(err, "mysql list api failed"))
		return
	}
	for _, aa := range appApis {

		if _, ok := am.APPAPIs[aa.AppID]; !ok {
			am.APPAPIs[aa.AppID] = set.NewSet[string]()
		}

		if _, ok := apiSet[aa.ApiID]; !ok {
			continue
		}

		appAPI := am.APPAPIs[aa.AppID]
		appAPI.Add(apiSet[aa.ApiID])
	}

	a.APPAPIs = am.APPAPIs
	a.Apps = am.Apps
	return
}

func (a *ApplicationMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if val, ok := c.Get(constant.Verify_Status); ok && val == "success" {
			c.Next()
			return
		}
		err := c.Request.ParseForm()
		if err != nil {
			c.Next()
			return
		}
		var params = make(map[string]interface{})
		for name := range c.Request.Form {
			params[name] = c.Request.Form.Get(name)
		}
		// 非公共请求接口
		if !isExist(params, "sign") ||
			!isExist(params, "app_key") ||
			!isExist(params, "ts") ||
			!isExist(params, "sign_method") ||
			!isExist(params, "v") {
			c.Next()
			return
		}
		//API权限校验
		originSign := params["sign"]
		delete(params, "sign")

		appKey := params["app_key"].(string)
		var (
			app *models.APP
			ok  bool
		)
		// 应用校验
		if app, ok = a.Apps[appKey]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "app is not current",
			})
			c.Abort()
			return
		}

		// 接口权限校验
		apiPath := c.Request.URL.Path
		if aa, ok := a.APPAPIs[app.ID]; !ok || !aa.Contains(apiPath) {
			c.JSON(http.StatusForbidden, map[string]interface{}{
				"code": 403,
				"msg":  "api is valid",
			})
			c.Abort()
			return
		}

		// 签名校验
		sign, err := a.calcSign(params, app)
		if err != nil {
			c.JSON(http.StatusForbidden, map[string]interface{}{
				"code": 403,
				"msg":  err.Error(),
			})
			c.Abort()
			return
		}
		if sign != originSign {
			c.JSON(http.StatusForbidden, map[string]interface{}{
				"code": 403,
				"msg":  "sign is not current",
			})
			c.Abort()
			return
		}
		c.Set(constant.Verify_Status, "success")
		c.Next()
		return
	}
}
func isExist(params map[string]interface{}, key string) bool {
	_, ok := params[key]
	return ok
}

func (a *ApplicationMiddleware) calcSign(params map[string]interface{}, app *models.APP) (string, error) {

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var paramsStr string
	for _, k := range keys {
		paramsStr += fmt.Sprintf("%s%s", k, params[k])
	}
	paramsStr += app.Secret + paramsStr + app.Secret

	md5 := strings.ToLower(utils.MD5(paramsStr))
	return md5, nil
}
