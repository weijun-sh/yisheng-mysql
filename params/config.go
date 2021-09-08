package params

import (
	"encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/weijun-sh/yisheng-mysql/common"
	"github.com/weijun-sh/yisheng-mysql/log"
)

var (
	configFile string
	scanConfig = &Config{}
	mongodbConfig = &DBConfig{}
	mysqldbConfig = &DBConfig{}
)

type Config struct {
       MongoDB *DBConfig
       MysqlDB *DBConfig
}

// MongoDBConfig mongodb config
type DBConfig struct {
       DBURL      string
       DBName     string
       UserName   string `json:"-"`
       Password   string `json:"-"`
}

// GetMongodbConfig get mongodb config
func GetMongodbConfig() *DBConfig {
       return mongodbConfig
}

// GetMysqldbConfig get mysqldb config
func GetMysqldbConfig() *DBConfig {
       return mysqldbConfig
}

// GetScanConfig get scan config
func GetScanConfig() *Config {
	return scanConfig
}

// LoadConfig load config
func LoadConfig(filePath string) *Config {
	log.Println("LoadConfig Config file is", filePath)
	if !common.FileExist(filePath) {
		log.Fatalf("LoadConfig error: config file '%v' not exist", filePath)
	}

	config := &Config{}
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		log.Fatalf("LoadConfig error (toml DecodeFile): %v", err)
	}

	var bs []byte
	if log.JSONFormat {
		bs, _ = json.Marshal(config)
	} else {
		bs, _ = json.MarshalIndent(config, "", "  ")
	}
	log.Println("LoadConfig finished.", string(bs))

       mongodbConfig = config.MongoDB
       mysqldbConfig = config.MysqlDB

	configFile = filePath // init config file path
	return scanConfig
}

