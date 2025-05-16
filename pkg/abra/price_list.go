package abra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/anaryk/maximal-limit-abra-sync/pkg/internal"
)

func (c *Connector) CheckIfPriceItemExists(eanKod string) (bool, int, error) {
	url := fmt.Sprintf("%s/c/%s/cenik/(eanKod='%s').json", internal.AbraBaseURL, internal.AbraCompany, eanKod)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, 0, err
	}

	var priceItem CenikWrapper

	json.Unmarshal(body, &priceItem)

	if len(priceItem.Winstrom.Cenik) == 0 {
		return false, 0, nil
	} else {
		id, err := strconv.Atoi(priceItem.Winstrom.Cenik[0].ID)
		if err != nil {
			return false, 0, fmt.Errorf("failed to convert ID to int: %w", err)
		}
		return true, id, nil
	}
}

func (c *Connector) CreatePriceItem(priceItem Cenik) (*APIResponse, error) {
	url := fmt.Sprintf("%s/c/%s/cenik.json", internal.AbraBaseURL, internal.AbraCompany)

	cenW := CenikWrapper{
		Winstrom: CenikWrapperWinstrom{
			Version: "1.0",
			Cenik:   []Cenik{priceItem},
		},
	}

	jsonData, err := json.Marshal(cenW)
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

func (c *Connector) UpdatePriceItem(priceItem Cenik, priceItemID int) (*APIResponse, error) {
	url := fmt.Sprintf("%s/c/%s/cenik/%s.json", internal.AbraBaseURL, internal.AbraCompany, strconv.Itoa(priceItemID))

	cenW := CenikWrapper{
		Winstrom: CenikWrapperWinstrom{
			Version: "1.0",
			Cenik:   []Cenik{priceItem},
		},
	}

	jsonData, err := json.Marshal(cenW)
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
