package product

import _type "bot/pkg/type"

type Product struct {
	Name        string `json:"name"`
	StoreID     string `json:"storeId,omitempty"`
	UUID        string `json:"uuid"`
	ParentUUID  string `json:"parentUuid" `
	Group       bool   `json:"group"`
	Image       string `json:"image"`
	MeasureName string `json:"measureName,omitempty"`
	Price       uint64 // Цена в копейках
}

func (p *Product) PriceString() string {
	return _type.FormatPrice(p.Price)
}
