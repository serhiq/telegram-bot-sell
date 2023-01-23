package restoClient

import (
	"bot/pkg/restoClient/order"
	domainOrder "bot/pkg/type/order"
	"fmt"
	"strconv"
	"time"
)

func (c RestoClient) PostOrder(o *domainOrder.Order) (*order.Response, error) {
	request := order.ToResponse(o)

	orderResponse := &order.Response{}

	request.State = "new"
	request.ID = GetRandomOrderNumber()

	endpoint := c.options.BaseUrl + "/order/" + c.options.Store
	resp, err := c.client.R().
		SetHeader("Authorization", c.options.Auth).
		SetBody(request).
		SetResult(&orderResponse).
		Post(endpoint)

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("client: postOrder is ERROR: code == %d", resp.StatusCode())
	}

	// todo compare order and orderResponse

	return orderResponse, nil

}

func GetRandomOrderNumber() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
	//return strconv.FormatInt(time.Now().UnixNano(), 10)
}
