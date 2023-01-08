package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// не используйте ссылку на текущий каталог, поскольку приложение может быть запущено с любым рабочим каталогом, что приведет к неочевидному поведению
// я бы рекомендовал использовать пути относительно исполняемого файла
const TempPatch = "./tmp/"
const PreviewCachePatch = "./imageCache/"

type Config struct {
	Token   string
	Auth    string
	Store   string
	BaseUrl string
}

func New() (*Config, error) {
	config := &Config{}

	// такой код предполагает, что пустое значение тоже допустимо, хотя по логике - нет
	token, tokenExists := os.LookupEnv("TELEGRAM_TOKEN")
	if !tokenExists {
		return nil, fmt.Errorf("config: %s is not set", "TELEGRAM_TOKEN")
	}
	config.Token = token

	auth, authExist := os.LookupEnv("EVOTOR_TOKEN")
	if !authExist {
		return nil, fmt.Errorf("config: %s is not set", "EVOTOR_TOKEN")
	}

	config.Auth = auth

	store, storeExist := os.LookupEnv("EVOTOR_STORE")
	if !storeExist {
		return nil, fmt.Errorf("config: %s is not set", "EVOTOR_STORE")
	}

	config.Store = store

	baseUrl, baseUrlExist := os.LookupEnv("BASE_URL")
	if !baseUrlExist {
		return nil, fmt.Errorf("config: %s is not set", "BASE_URL")
	}

	config.BaseUrl = baseUrl

	if err := os.MkdirAll(filepath.Dir(TempPatch), fs.ModeDir); err != nil {
		return nil, fmt.Errorf("config: failed creating tmp path %s (%s)", filepath.Dir(TempPatch), err)
	}

	if err := os.MkdirAll(filepath.Dir(PreviewCachePatch), fs.ModeDir); err != nil {
		return nil, fmt.Errorf("config: failed creating cache path %s (%s)", filepath.Dir(TempPatch), err)
	}

	return config, nil
}
