package common

import (
	//"fmt"
	"gin/model"
	//"github.com/jinzhu/gorm"
	"gorm.io/gorm"
	//"github.com/spf13/viper"
	"gorm.io/driver/postgres"
)

var DB *gorm.DB

func InitDB() * gorm.DB{
	//driverName := viper.GetString("datasource.driverName")
	//host := viper.GetString("datasource.host")
	//port := viper.GetString("datasource.port")
	//database := viper.GetString("datasource.database")
	//username := viper.GetString("datasource.username")
	//password := viper.GetString("datasource.password")
	//charset := viper.GetString("datasource.charset")
	//args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
	//	username,
	//	password,
	//	host,
	//	port,
	//	database,
	//	charset)
	dsn := "host=120.77.221.202 port=8888 user=gaussdb dbname=postgres password=Enmo@123 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//db, err := gorm.Open(driverName,args)
	if err != nil {
		panic( "failed to connect database, err:" + err.Error())
	}
	db.AutoMigrate(&model.User{})

	DB = db

	return db
}

func GetDB() *gorm.DB {
	return DB
}
