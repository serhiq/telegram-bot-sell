package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Store struct {
	Db *gorm.DB
}

type Settings struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Password string
}

func (s Settings) toDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		s.User,
		s.Password,
		s.Host,
		s.Port,
		s.Database)
}

func New(settings Settings) (*Store, error) {

	config := &gorm.Config{
		PrepareStmt: false,
	}

	db, err := gorm.Open(mysql.Open(settings.toDSN()), config)

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
