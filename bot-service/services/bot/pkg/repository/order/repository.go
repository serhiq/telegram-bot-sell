package order

import (
	"bot/pkg/restoClient"
	domainOrder "bot/pkg/type/order"
)

type Repository struct {
	evoClient *restoClient.RestoClient
}

func New(c *restoClient.RestoClient) *Repository {
	return &Repository{
		evoClient: c,
	}
}

func (r Repository) Send(order *domainOrder.Order) error {

	_, err := r.evoClient.PostOrder(order)
	//postOrder, err := r.evoClient.PostOrder(order)
	if err != nil {
		return err
	}

	//	 todo check postorder and order

	return nil

}
