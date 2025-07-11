package svc

import (
	"elevate-hub/config"
	cache2 "elevate-hub/db/cache"
	"elevate-hub/db/cache/memory"
	"elevate-hub/db/cache/redis"
	models2 "elevate-hub/db/models"
	"elevate-hub/db/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Mysql  Mysql
	Cache  cache2.Cache
}

type Mysql struct {
	DB                *gorm.DB
	APIModel          *models2.APIModel
	APPModel          *models2.APPModel
	AuthRelationModel *models2.AuthRelationModel
	MenuModel         *models2.MenuModel
	RoleModel         *models2.RoleModel
	UserModel         *models2.UserModel
	DepartmentModel   *models2.DepartmentModel
}

func NewServiceContext(conf config.Config) *ServiceContext {
	db := mysql.InitMysql(conf.Mysql, nil)
	var cache cache2.Cache
	if conf.Redis != nil {
		cache = redis.NewRedisCache(conf.Redis)
	}
	cache = memory.NewMemoryCache()
	return &ServiceContext{
		Config: conf,
		Mysql: Mysql{
			DB:                db,
			APIModel:          models2.NewAPIModel(db),
			APPModel:          models2.NewAPPModel(db),
			AuthRelationModel: models2.NewAuthRelationAuthRelationModel(db),
			MenuModel:         models2.NewMenuModel(db),
			RoleModel:         models2.NewRoleModel(db),
			UserModel:         models2.NewUserModel(db),
			DepartmentModel:   models2.NewDepartmentModel(db),
		},
		Cache: cache,
	}
}
