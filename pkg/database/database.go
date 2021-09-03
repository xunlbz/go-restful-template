package database

import (
	"os"
	"runtime"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DEFAULTDB  数据库实例
var DEFAULTDB *gorm.DB
var models []interface{}
var DBName = "edge_admin.db"

func init() {
	if runtime.GOOS == "linux" {
		dir := "/etc/edge/"
		if exist, _ := pathExists(dir); !exist {
			os.Mkdir(dir, os.ModeDir)
		}
		DBName = "/etc/edge/edge_admin.db"
	}
}

// Open 连接数据库
func Open() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(DBName), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DEFAULTDB = db
	AutoMigrationAllModels()
	return db
}

// Close 连接数据库
func Close() {
}

func AutoMigration(dst ...interface{}) {
	DEFAULTDB.AutoMigrate(dst...)
}

//AutoMigrate 数据库升级
func AutoMigrationAllModels() {
	DEFAULTDB.AutoMigrate(models...)
	models = nil
}
