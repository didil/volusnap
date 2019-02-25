package api

import (
	"database/sql"
	"fmt"

	// postgres driver
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func openDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable", viper.GetString("db.host"), viper.GetInt("db.port"), viper.GetString("db.user"), viper.GetString("db.name"), viper.GetString("db.password"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
