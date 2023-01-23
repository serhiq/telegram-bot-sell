package order

import (
	"bot/pkg/type/product"
	"encoding/json"
	"strconv"
	"strings"
)

type Order struct {
	ID       string
	Contacts struct {
		Phone string
	}
	Positions []Position
	State     string
	Comment   string
}

type Position struct {
	ProductUUID       string
	ProductName       string
	Price             int
	PriceWithDiscount int
	Quantity          int
}

func (c *Order) AddItem(item *product.Product) {

	for i, position := range c.Positions {
		if position.ProductUUID == item.UUID {
			c.Positions[i].Quantity = position.Quantity + 1
			return
		}
	}

	c.Positions = append(c.Positions, Position{
		ProductUUID:       item.UUID,
		ProductName:       item.Name,
		Price:             item.Price,
		PriceWithDiscount: item.Price,
		Quantity:          1,
	})
}

func (c *Order) DecreaseMenuItem(item *product.Product) {

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

func (c *Order) CountPosition() string {
	var count = 0
	for _, position := range c.Positions {
		count = count + position.Quantity
	}
	return strconv.FormatInt(int64(count), 10)
}

func (c *Order) CountItemPosition(uuid string) string {
	var count = 0
	for _, position := range c.Positions {
		if position.ProductUUID == uuid {
			count = count + position.Quantity
		}
	}
	return strconv.FormatInt(int64(count), 10)
}

func (c *Order) SumPositions() string {
	var sum = 0
	for _, position := range c.Positions {
		sum = sum + (position.Quantity * position.PriceWithDiscount)
	}

	return strconv.FormatInt(int64(sum), 10) + "руб"
}

func (c *Order) OrderDescription() string {
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

func (c *Order) ToJson() (string, error) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
