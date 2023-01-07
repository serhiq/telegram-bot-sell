package database

import (
	"bot/internal/entity"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

func Init() (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		//dsn := fmt.Sprintf("%s:%s@tcp(%s:3333)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		getenv("MYSQL_USER", "user_app"),
		getenv("MYSQL_PASSWORD", "password"),
		getenv("MYSQL_HOST", "localhost"),
		getenv("MYSQL_DATABASE", "a_bot_db"))
	config := createGornConfig()

	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&entity.Chat{}, &entity.MenuItemDatabase{})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	sqlDB.SetConnMaxLifetime(time.Second)
	return &Database{
		Db: db,
	}, nil
}

func (r *Database) DBHealthCheck() error {
	sqlDB, err := r.Db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

func (r *Database) ShutDown() {
	db, err := r.Db.DB()
	if err != nil {
		log.Printf("database: error close database, %s", err)
	}

	err = db.Close()
	if err != nil {
		log.Printf("database: error close database, %s", err)
		return
	}
	log.Print("database: close")
}

func createGornConfig() *gorm.Config {
	return &gorm.Config{
		PrepareStmt: false,
	}
}

type Database struct {
	Db *gorm.DB
}

func getenv(key string, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}

	return v
}
