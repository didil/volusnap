package api

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// 
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

func openDB() (*gorm.DB, error) {
	connStr := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable", viper.GetString("db.host"), viper.GetInt("db.port"), viper.GetString("db.user"), viper.GetString("db.name"), viper.GetString("db.password"))
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}).Error

	return err
}
