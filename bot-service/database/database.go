package database

import (
	"bot/internal/entity"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// классическое наиманование функции для создания экземпляра структура - `New` или `NewStruct`, где `Struct` - наименование структуры
func Init() (*Database, error) {
	// загрузка настроек из переменных окружения - явно не задача кода работы с БД
	// настройки должны приходить "свыше"
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
	// не совсем понятно зачем так делать, вся суть пула соединений в возможности удержания подключения и его повторного использования
	sqlDB.SetConnMaxLifetime(time.Second)
	return &Database{
		Db: db,
	}, nil
}

// если у нас структура называется `Database`, то зачем в наименовании метода писать `DB`?
func (r *Database) DBHealthCheck() error {
	sqlDB, err := r.Db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// `ShutDown` громко сказано - как будто мы выключем всю базу данных, хотя по факту - просто закрываем соединение
func (r *Database) ShutDown() {
	// в действительности все это зря, ибо "под капотом" в Gorm используется пул соединений, а потому соединений может быть несколько и этот код в лучшем случае закроет одно из них,
	// а в худшем - сначала откроет новое, а потом сразу его закрое
	// насколько я помню, у Gorm нет внятного способа закрыть все соединения пула
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
