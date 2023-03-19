package order

import (
	domainOrder "bot/pkg/type/order"
	evo "github.com/softc24/evotor-resto-go"
	"strconv"
	"time"
)

func ToResponse(response *domainOrder.Order) *evo.Order {

	positions := []evo.OrderPosition{}

	for _, p := range response.Positions {
		positions = append(positions, evo.OrderPosition{
			ProductUUID:       p.ProductUUID,
			ProductName:       p.ProductName,
			Price:             evo.Money(p.Price),
			PriceWithDiscount: evo.Money(p.Price),
			Quantity:          evo.Quantity(p.Quantity),
		})
	}
	// создаем заказ
	order := evo.MakeOrder(strconv.FormatInt(time.Now().UnixMilli(), 32), response.Comment, evo.Contacts{
		Phone: response.Contacts.Phone,
	}, positions)

	return &order
}
