package entity

type MenuResponse struct {
	Name string `json:"name"`
	//Code          string        `json:"code"`
	//Barcodes      []interface{} `json:"barcodes,omitempty"`
	//ID string `json:"id"`
	//StoreID       string        `json:"storeId,omitempty"`
	//UserID        string        `json:"userId,omitempty"`
	//CreatedAt     string        `json:"createdAt,omitempty"`
	//UpdatedAt     string        `json:"updatedAt,omitempty"`
	UUID string `json:"uuid"`
	//UserUUID      string        `json:"user_uuid,omitempty"`
	//StoreUUID     string        `json:"store_uuid,omitempty"`
	//BarCodes      []interface{} `json:"barCodes,omitempty"`
	ParentUUID string `json:"parentUuid"`
	Group      bool   `json:"group"`
	//Parents       []interface{} `json:"parents"`
	//Fullname      string        `json:"fullname"`
	Image         string `json:"image"`
	IsUnavailable string `json:"isUnavailable"`
	//ParentID      string        `json:"parentId,omitempty"`
	//Type          string        `json:"type,omitempty"`
	Price int `json:"price,omitempty"`	// а разве цена не может содержать копейки?
	//Tax           string        `json:"tax,omitempty"`
	MeasureName string `json:"measureName,omitempty"`
	//TaxNumber     string        `json:"taxNumber,omitempty"`
	ProductType string `json:"productType,omitempty"`
}

type Menu []MenuResponse

func (m *MenuResponse) CanAddToOrder() bool {
	// я бы изобрел свой тип `bool` через интерфейсы `json.Marshaler`/`json.Unmarshaler`
	if m.Group && (m.IsUnavailable == "" || m.IsUnavailable == "0") {
		return true
	}

	if !m.Group && (m.IsUnavailable == "" || m.IsUnavailable == "0") && m.ProductType == "NORMAL" {
		return true
	}

	return false
}
