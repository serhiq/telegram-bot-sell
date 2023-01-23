package main

import (
	"bot/config"
	"bot/pkg/restoClient"
	"bot/pkg/store/mysql"
	repositoryChat "bot/services/bot/pkg/repository/chat"
	repositoryOrder "bot/services/bot/pkg/repository/order"
	"bot/services/bot/pkg/repository/product"
	"bot/services/bot/pkg/worker"
	"bot/services/ssbot/internal/delivery/bot"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.Info("Initializing bot...")

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	store, err := mysql.New(mysql.Settings{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     3306,
		Database: os.Getenv("MYSQL_DATABASE"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
	})

	if err != nil {
		panic(err)
	}

	err = store.Db.AutoMigrate(&repositoryChat.Chat{}, &product.Product{})
	if err != nil {
		panic(err)
	}

	restyClient := resty.New()
	evoClient := restoClient.New(restyClient, &restoClient.Options{
		Auth:    cfg.Auth,
		Store:   cfg.Store,
		BaseUrl: cfg.BaseUrl,
	})

	var repoProduct = product.New(store.Db)
	var repoChat = repositoryChat.New(store.Db)
	var repoOrder = repositoryOrder.New(evoClient)

	syncWorker := worker.New(repoProduct, evoClient, restyClient)
	syncWorker.EnqueueUniquePeriodicWork()

	sBot, err := bot.New(bot.Options{
		Token: cfg.Token,
	}, repoProduct, repoChat, repoOrder)

	if err != nil {
		log.Fatalf("Cannot initialize sBot: %s", err)
	}

	err = sBot.Start()
	if err != nil {
		log.Fatalf("Cannot start sBot: %s", err)
	}

	//configureCloseSignal(an)
}

//
//func configureCloseSignal(an *app.An) {
//	c := make(chan os.Signal)
//	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
//	go func() {
//		<-c
//		log.Println("Gracefully shutting down...")
//		an.Shutdown()
//		os.Exit(0)
//	}()
//}
