package restoClient

import "github.com/go-resty/resty/v2"

type RestoClient struct {
	client *resty.Client

	options *Options
}

type Options struct {
	Auth    string
	Store   string
	BaseUrl string
}

func New(client *resty.Client, options *Options) *RestoClient {
	return &RestoClient{
		client:  client,
		options: options,
	}
}
