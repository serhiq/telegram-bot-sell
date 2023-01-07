package entity

type PostOrderResponse struct {
	UUID      string      `json:"uuid"`
	ID        string      `json:"id"`
	UserUUID  string      `json:"user_uuid"`
	StoreUUID string      `json:"store_uuid"`
	Number    string      `json:"number"`
	State     string      `json:"state"`
	Contacts  interface{} `json:"contacts"`
	Comment   string      `json:"comment"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
	Positions []struct {
		OrderUUID               string      `json:"order_uuid"`
		Position                string      `json:"position"`
		ProductUUID             string      `json:"product_uuid"`
		ProductName             string      `json:"product_name"`
		ProductCode             string      `json:"product_code"`
		ProductMeasureName      string      `json:"product_measureName"`
		ProductMeasurePrecision string      `json:"product_measurePrecision"`
		ProductType             string      `json:"product_type"`
		ProductTax              interface{} `json:"product_tax"`
		Price                   string      `json:"price"`
		PriceWithDiscount       string      `json:"priceWithDiscount"`
		Quantity                string      `json:"quantity"`
	} `json:"positions"`
}
