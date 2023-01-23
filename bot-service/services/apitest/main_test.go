package apitest

import (
	restoClient "bot/pkg/restoClient"
	"bot/pkg/restoClient/menu"
	domainOrder "bot/pkg/type/order"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

const base = ""
const store = ""
const auth = ""

func Test_main(t *testing.T) {
	c := resty.New()
	options := restoClient.Options{
		Auth:    auth,
		Store:   store,
		BaseUrl: base,
	}

	client := restoClient.New(c, &options)
	items, err := client.GetMenu()
	assert.Equal(t, nil, err, err)
	assert.True(t, len(items) != 0)

	//	 делаем заказ 1 позицию normal with price
	order := createOrder(items)

	postOrder, err := client.PostOrder(order)
	assert.Equal(t, nil, err, err)
	assert.Equal(t, order.State, postOrder.State, postOrder.State)
	assert.Equal(t, order.ID, postOrder.ID, postOrder.ID)
}

func createOrder(m menu.Menu) *domainOrder.Order {
	order := domainOrder.Order{
		ID: restoClient.GetRandomOrderNumber(),
		Contacts: struct {
			Phone string
		}{
			Phone: "89239990599",
		},
		Positions: nil,
		State:     "new",
		Comment:   "happy_patch",
	}

	for _, position := range m {
		if position.ProductType == "NORMAL" && position.Price != 0 {
			order.Positions = append(order.Positions, domainOrder.Position{
				ProductUUID:       position.UUID,
				ProductName:       position.Name,
				Price:             position.Price,
				PriceWithDiscount: position.Price,
				Quantity:          1,
			})
			return &order
		}
	}
	return &order
}
