package rest

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sebastianreh/chatroom-bots/stocks/internal/config"
	"github.com/sebastianreh/chatroom/pkg/logger"
	str "github.com/sebastianreh/chatroom/pkg/strings"
)

const (
	endpoint              = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"
	GetStockCSVMethodName = "GetStockCSV"
	apiClientName         = "api_client"
)

type (
	StooqClient interface {
		GetStockCSV(stockName string) ([]byte, error)
	}

	stooqClient struct {
		configs    config.Config
		logger     logger.Logger
		restClient *resty.Client
	}
)

func NewStooqClient(logger logger.Logger, restClient *resty.Client) StooqClient {
	return &stooqClient{
		logger:     logger,
		restClient: restClient,
	}
}

func (client *stooqClient) GetStockCSV(stockName string) ([]byte, error) {
	var res []byte
	var err error
	var resp *resty.Response

	req := client.restClient.R()
	url := fmt.Sprintf(endpoint, stockName)
	resp, err = req.Get(url)
	if err != nil {
		client.logger.Error(str.ErrorConcat(err, apiClientName, GetStockCSVMethodName))
		return res, err
	}

	if !resp.IsSuccess() {
		err = fmt.Errorf("error update api company with https status code: %d, body %s",
			resp.StatusCode(), string(resp.Body()))
		client.logger.Error(str.ErrorConcat(err, apiClientName, GetStockCSVMethodName))
		return res, err
	}

	return resp.Body(), nil
}
