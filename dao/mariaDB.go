package dao

import (
	"os"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func InitDb() (*gorm.DB, error) {

	db, err := gorm.Open(os.Getenv("DB_TYPE"), os.Getenv("DB_CONNEC_STRING"))

	if err != nil {
		return nil, err
	}

	return db, nil
}