package main

import (
	"bot/pkg/restoClient"
	"bot/pkg/store/mysql"
	"bot/services/bot/pkg/config"
	repositoryChat "bot/services/bot/pkg/repository/chat"
	repositoryOrder "bot/services/bot/pkg/repository/order"
	"bot/services/bot/pkg/repository/product"
	"bot/services/bot/pkg/worker"
	"bot/services/ssbot/internal/delivery/bot"
	r "github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Initializing bot...")

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	store, err := mysql.New(cfg.DBConfig)

	if err != nil {
		panic(err)
	}

	err = store.Db.AutoMigrate(&repositoryChat.Chat{}, &product.Product{})
	if err != nil {
		panic(err)
	}

	client := r.New()
	evoClient := restoClient.New(client, &restoClient.Options{
		Auth:    cfg.RestaurantAPI.Auth,
		Store:   cfg.RestaurantAPI.Store,
		BaseUrl: cfg.RestaurantAPI.BaseURL,
	})

	var repoProduct = product.New(store.Db)
	var repoChat = repositoryChat.New(store.Db)
	var repoOrder = repositoryOrder.New(evoClient)

	syncWorker := worker.New(repoProduct, evoClient, r.New())
	syncWorker.EnqueueUniquePeriodicWork()

	sBot, err := bot.New(bot.Options{
		Token: cfg.Telegram.Token,
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
