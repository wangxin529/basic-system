package config

import "elevate-hub/db/conf"

type Config struct {
	Port  int         `json:"port"`
	JWT   JWT         `json:"jwt"`
	Mysql conf.Mysql  `json:"mysql"`
	Redis *conf.Redis `json:"redis"`
}

type JWT struct {
	Secret  string `json:"secret"`
	Timeout int    `json:"timeout"`
}
