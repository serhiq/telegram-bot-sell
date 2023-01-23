package product

import "strconv"

type Product struct {
	Name        string `json:"name"`
	StoreID     string `json:"storeId,omitempty"`
	UUID        string `json:"uuid"`
	ParentUUID  string `json:"parentUuid" `
	Group       bool   `json:"group"`
	Image       string `json:"image"`
	MeasureName string `json:"measureName,omitempty"`
	Price       int
}

func (p *Product) PriceString() string {
	return strconv.FormatInt(int64(p.Price), 10) + " руб"
}
