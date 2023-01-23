package order

import (
	domainOrder "bot/pkg/type/order"
)

func ToResponse(response *domainOrder.Order) *OrderRequest {

	positions := []PositionOrder{}

	for _, p := range response.Positions {
		positions = append(positions, PositionOrder{
			ProductUUID:       p.ProductUUID,
			ProductName:       p.ProductName,
			Price:             p.Price,
			PriceWithDiscount: p.PriceWithDiscount,
			Quantity:          p.Quantity,
		})
	}

	return &OrderRequest{
		ID: response.ID,
		Contacts: struct {
			Phone string `json:"phone"`
		}{
			Phone: response.Contacts.Phone,
		},
		Positions: positions,
		State:     "",
		Comment:   "",
	}
}
