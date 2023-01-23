package repository

import (
	domainOrder "bot/pkg/type/order"
	"bot/pkg/type/product"
	"bot/services/bot/pkg/repository/chat"
	product2 "bot/services/bot/pkg/repository/product"
)

type OrderRepository interface {
	Send(order *domainOrder.Order) error
}

type ProductRepository interface {
	ImportMenu(items []*product2.Product) error
	GetMenuItemByParent(parent string) ([]*product.Product, error)
	GetMenu(id string) (*product.Product, error)
}

type ChatRepository interface {
	InsertChat(chat *chat.Chat) error
	UpdateChat(chat *chat.Chat) error
	GetChat(id int64) (*chat.Chat, error)
	GetOrCreateChat(id int64) (*chat.Chat, error)
	DeleteChat(id string) error
}
