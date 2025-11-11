package database

import (
	"fmt"
	"log"

	"github.com/xiaoyuer1231231/gin_mysql_grom_project/config"
	"github.com/xiaoyuer1231231/gin_mysql_grom_project/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDataBase(configParm *config.Config) (*gorm.DB, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configParm.Database.User,
		configParm.Database.Password,
		configParm.Database.Host,
		configParm.Database.Port,
		configParm.Database.Database)
	fmt.Printf("数据库配置: %+v\n", configParm.Database)
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	log.Println("Database connected successfully")
	return db, nil
}
