package configs

import (
	"fmt"
	"stripe-subscription/models"
	"stripe-subscription/shared/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type connection struct {
	db *gorm.DB
}

var db *gorm.DB

func InitDB() {
	var err error
	User := Username()
	Password := Password()
	Host := Host()
	Port := DBPort()
	DbName := DbName()
	createDBDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", User, Password, Host, Port)
	database, _ := gorm.Open(mysql.Open(createDBDsn), &gorm.Config{})
	_ = database.Exec("CREATE DATABASE IF NOT EXISTS " + DbName + ";")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", User, Password, Host, Port, DbName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		panic(err.Error())
	}
}

func MigrateModels() {
	conn := NewConnection()
	conn.GetDB().AutoMigrate(
		&models.Customer{},
	)
}

func NewConnection() *connection {
	return &connection{
		db,
	}
}

func (self *connection) GetDB() *gorm.DB {
	return self.db
}

func Close() error {
	dbSQL, err := db.DB()
	if err != nil {
		return err
	}
	return dbSQL.Close()
}
