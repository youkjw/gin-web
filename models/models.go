package models

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

type Database struct {
	Type string
	User string
	Password string
	Host string
	Name string
	TablePrefix string
}

type Yqw struct {
	Database
}

var DatabaseSetting = &Database{}
var YqwSetting = &Yqw{}

var db map[string]*gorm.DB

var (
	Cfg *ini.File
)

type Model struct {
	ID int `grom:"primary_key;autoIncrement" json:"id"`
	CreateAt int `json:"create_at"`
	UpdateAt int `json:"update_at"`
}

func Setup() {
	err := Cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo DatabaseSetting err: %v", err)
	}

	err = Cfg.Section("yqw").MapTo(YqwSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo YqwSetting err: %v", err)
	}

	DatabaseSetting.connect()
	YqwSetting.connect()
}

func (database *Database) connect() {
	var (
		err error
		dbType, dbName, user, password, host, tablePrefix string
	)

	dbType = database.Type
	dbName = database.Name
	user = database.User
	password = database.Password
	host = database.Host
	tablePrefix = database.TablePrefix

	db[database.Name], err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName))

	if err != nil {
		log.Println(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	db[database.Name].SingularTable(true)
	db[database.Name].LogMode(true)
	db[database.Name].DB().SetMaxIdleConns(10)
	db[database.Name].DB().SetMaxOpenConns(100)
	db[database.Name].DB().SetConnMaxLifetime(time.Hour)
}

func CloseDB(db *gorm.DB) {
	defer db.Close()
}