package conf

type Mysql struct {
	User     string `json:"user"`
	Passwd   string `json:"passwd"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Logger   bool   `json:"logger"`
	//Charset  string `json:"charset"`
}

type Redis struct {
	Addr             []string `json:"addr"`
	Password         string   `json:"password"`
	Port             string   `json:"port"`
	MasterName       string   `json:"masterName"`
	SentinelPassword string   `json:"sentinelPassword"`
	DB               int      `json:"db"`
}
