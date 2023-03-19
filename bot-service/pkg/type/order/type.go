package order

import (
	_type "bot/pkg/type"
	"bot/pkg/type/product"
	"encoding/json"
	"fmt"
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
	Price             uint64
	PriceWithDiscount uint64
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
	var count uint64 = 0
	for _, position := range c.Positions {
		count = count + uint64(position.Quantity)
	}
	return strconv.FormatUint(count, 10)
}

func (c *Order) CountItemPosition(uuid string) string {
	var count uint64 = 0
	for _, position := range c.Positions {
		if position.ProductUUID == uuid {
			count = count + uint64(position.Quantity)
		}
	}
	return strconv.FormatUint(count, 10)
}

func (c *Order) SumPositions() string {
	var sum uint64 = 0
	for _, position := range c.Positions {
		sum = sum + (uint64(position.Quantity) * uint64(position.PriceWithDiscount))
	}

	return strconv.FormatUint(sum, 10) + "руб"
}

func (c *Order) calculateTotal() uint64 {
	var total uint64
	for _, position := range c.Positions {
		// Используем цену с учетом скидки, если она задана, иначе используем базовую цену
		var price uint64
		if position.PriceWithDiscount != 0 {
			price = position.PriceWithDiscount
		} else {
			price = position.Price
		}

		total += price * uint64(position.Quantity)
	}

	return total
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
	builder.WriteString("\nЕсли все верно выбиерите действие:")

	return builder.String()
}

func (c *Order) OrderDescriptionNew() string {
	b := strings.Builder{}

	b.WriteString("\n\nСостав заказа:")

	for i, pos := range c.Positions {
		var itemPrice uint64
		if pos.PriceWithDiscount != 0 {
			itemPrice = pos.PriceWithDiscount
		} else {
			itemPrice = pos.Price
		}
		itemStr := fmt.Sprintf("\n\n%d. %s\n    Кол-во: %d\n    Цена: %s \n", i+1, pos.ProductName, pos.Quantity, _type.FormatPrice(itemPrice))

		b.WriteString(itemStr)
	}

	b.WriteString("\nОбщая сумма заказа: " + _type.FormatPrice(c.calculateTotal()))

	b.WriteString("\n")
	b.WriteString("\nЕсли все верно выберите действие:")

	return b.String()
}

func (c *Order) ToJson() (string, error) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
