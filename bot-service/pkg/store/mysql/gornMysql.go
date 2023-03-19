package mysql

import (
	"bot/services/bot/pkg/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Store struct {
	Db *gorm.DB
}

func New(s config.DBConfig) (*Store, error) {

	cfg := &gorm.Config{
		PrepareStmt: false,
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		s.Username,
		s.Password,
		s.Host,
		s.Port,
		s.DatabaseName)

	db, err := gorm.Open(mysql.Open(dsn), cfg)

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	return &Store{Db: db}, nil
}
