package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gabuladze/tolling/types"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) AggregateDistance(ctx context.Context, aggReq *types.AggregateRequest) error {
	b, err := json.Marshal(aggReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("the service responded with non 200 status code %d", res.StatusCode)
	}
	return nil
}

func (c *HTTPClient) GetInvoice(ctx context.Context, aggReq *types.GetInvoiceRequest) (*types.Invoice, error) {
	endpoint := fmt.Sprintf("%s/%s?obu=%d", c.Endpoint, "invoice", aggReq.ObuID)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("the service responded with non 200 status code %d", res.StatusCode)
	}

	var inv types.Invoice
	if err := json.NewDecoder(res.Body).Decode(&inv); err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return &inv, nil
}
