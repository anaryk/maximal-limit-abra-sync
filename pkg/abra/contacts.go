package abra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/anaryk/maximal-limit-abra-sync/pkg/internal"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/utils"
)

func GenerateContactJSON(data ContactData) ([]byte, error) {
	shortCode := utils.GenerateShortCode(data.Name)
	contact := map[string]interface{}{
		"winstrom": map[string]interface{}{
			"adresar": []map[string]interface{}{
				{
					"nazev": data.Name,
					"kod":   shortCode,
					"ulice": data.Street,
					"mesto": data.City,
					"psc":   data.PostalCode,
					"email": data.Email,
					"mobil": data.Mobile,
				},
			},
		},
	}

	jsonData, err := json.MarshalIndent(contact, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (c *Connector) CreateContact(data ContactData) (*APIResponse, error) {
	generatedJson, err := GenerateContactJSON(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/c/%s/%s", internal.AbraBaseURL, internal.AbraCompany, internal.AbraCustomerContantURL), bytes.NewBuffer(generatedJson))
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

func (c *Connector) CheckIfContactExist(contactCode string) (*APIResponseContacts, error) {
	url := fmt.Sprintf("%s/c/%s/%s/(kod == '%s').json", internal.AbraBaseURL, internal.AbraCompany, internal.AbraCustomerContantURL, contactCode)
	req, err := http.NewRequest("GET", url, nil)
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

	var apiResponse APIResponseContacts
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	return &apiResponse, nil
}

func (c *Connector) GetContactIDByShortName(shortName string) (string, error) {
	url := fmt.Sprintf("%s/c/%s/%s/(kod == '%s').json", internal.AbraBaseURL, internal.AbraCompany, internal.AbraCustomerContantURL, shortName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var apiResponse APIResponseContacts
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return "", err
	}

	if len(apiResponse.Winstrom.Adresar) == 0 {
		return "", fmt.Errorf("contact with short name '%s' not found", shortName)
	}

	return apiResponse.Winstrom.Adresar[0].ID, nil
}
