package restoClient

import (
	"bot/pkg/restoClient/menu"
	"fmt"
	"log"
)

func (c RestoClient) GetMenu() (menu.Menu, error) {
	result := menu.Menu{}

	endpoint := c.options.BaseUrl + "/product/" + c.options.Store
	resp, err := c.client.R().
		SetHeader("Authorization", c.options.Auth).
		SetResult(&result).
		Get(endpoint)

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		log.Println(resp)
		return nil, fmt.Errorf("client: getMenu is ERROR: code == %d, resp: == %s", resp.StatusCode(), resp)
	}

	return result, nil
}
