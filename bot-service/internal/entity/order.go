package entity

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type OrderRequest struct {
	ID       string `json:"id"`
	Contacts struct {
		Phone string `json:"phone"`
	} `json:"contacts"`
	Positions []PositionOrder `json:"positions"`
	State     string          `json:"state"`
	Comment   string          `json:"comment"`
}

type PositionOrder struct {
	ProductUUID       string `json:"product_uuid"`
	ProductName       string `json:"product_name"`
	Price             int    `json:"price"`
	PriceWithDiscount int    `json:"priceWithDiscount"`
	Quantity          int    `json:"quantity"`
}

func (c *OrderRequest) ToJson() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (c *OrderRequest) AddMenuItem(item *MenuItemDatabase) {

	for i, position := range c.Positions {
		if position.ProductUUID == item.UUID {
			c.Positions[i].Quantity = position.Quantity + 1
			return
		}
	}

	c.Positions = append(c.Positions, PositionOrder{
		ProductUUID:       item.UUID,
		ProductName:       item.Name,
		Price:             item.Price,
		PriceWithDiscount: item.Price,
		Quantity:          1,
	})
}

func (c *OrderRequest) DecreaseMenuItem(item *MenuItemDatabase) {

	for i, position := range c.Positions {
		if position.ProductUUID == item.UUID {
			if position.Quantity == 1 {
				c.Positions = append(c.Positions[:i], c.Positions[i+1:]...)
				return
			}

			c.Positions[i].Quantity = position.Quantity - 1
			return
		}
	}
}

func (c *OrderRequest) CountPosition() string {
	var count = 0
	for _, position := range c.Positions {
		count = count + position.Quantity
	}
	return strconv.FormatInt(int64(count), 10)
}

func (c *OrderRequest) CountItemPosition(uuid string) string {
	var count = 0
	for _, position := range c.Positions {
		if position.ProductUUID == uuid {
			count = count + position.Quantity
		}
	}
	return strconv.FormatInt(int64(count), 10)
}

func (c *OrderRequest) SumPositions() string {
	var sum = 0
	for _, position := range c.Positions {
		sum = sum + (position.Quantity * position.PriceWithDiscount)
	}

	return strconv.FormatInt(int64(sum), 10) + "руб"
}

// форматировать вывод явно не задача самого запроса на заказ
func (c *OrderRequest) OrderDescription() string {
	builder := strings.Builder{}

	builder.WriteString("\n\nСостав заказа:")

	for i, position := range c.Positions {
		builder.WriteString("\n\n" + strconv.FormatInt(int64(i+1), 10) + ". " + position.ProductName)
		builder.WriteString("\n" + "    Кол-во: " + strconv.FormatInt(int64(position.Quantity), 10))
		builder.WriteString("\n" + "    Цена: " + strconv.FormatInt(int64(position.Price), 10) + " руб")

	}
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString("\nОбщая сумма заказа: " + c.SumPositions())

	builder.WriteString("\n")
	builder.WriteString("\nЕсли все верно выбиерите действие:")

	return builder.String()
}

///////////////////////////////////////////////////////////////////////////////////////////
// уникальности до секунды может быть недостаточно, есть `time.UnixMilli()`/`time.UnixMicro()`/`time.UnixNano()`
// и всегда есть еще `uuid` и `nanoid`
func GetRandomOrderNumber() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
