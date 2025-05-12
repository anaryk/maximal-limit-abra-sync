package abra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/anaryk/maximal-limit-abra-sync/pkg/internal"
)

func (c *Connector) CreateSalesReceipt(salesReceipt SaleReceipt) (*APIResponse, error) {
	url := fmt.Sprintf("%s/c/%s/prodejka.json", internal.AbraBaseURL, internal.AbraCompany)

	jsonData, err := json.Marshal(salesReceipt)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse APIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	return &apiResponse, nil
}
