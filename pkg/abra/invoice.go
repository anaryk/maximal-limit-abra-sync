package abra

import (
	"bytes"
	"database/sql"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/anaryk/maximal-limit-abra-sync/pkg/internal"
	"github.com/anaryk/maximal-limit-abra-sync/pkg/utils"
)

func (c *Connector) CreateInvoice(customerCode, issueDate, dueDate string, internalNumber sql.NullString, items []FakturaPolozka) (*APIResponse, error) {
	contactId, err := c.GetContactIDByShortName(customerCode)
	if err != nil {
		return nil, err
	}
	if internalNumber.String == "" {
		internalNumber.String = fmt.Sprintf("FV-%s", utils.GenerateRandomString(8))
	}
	invoice := InvoiceRequest{
		Winstrom: struct {
			FakturaVydana []FakturaVydana `json:"faktura-vydana"`
		}{
			FakturaVydana: []FakturaVydana{
				{
					Kod:            internalNumber.String,
					DatVyst:        issueDate,
					DatSplat:       dueDate,
					StavUhrady:     "stavUhr.uhrazenoRucne",
					IdFirmy:        contactId,
					Polozky:        items,
					TypFaktury:     internal.AbraVydanaFakturaID,
					AccountingType: internal.AbraAccountingOperationID,
				},
			},
		},
	}

	payload, err := json.Marshal(invoice)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/c/%s/faktura-vydana.json", internal.AbraBaseURL, internal.AbraCompany)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
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

func (c *Connector) GetPDFInvoiceAsBase64(invoiceID string) (string, error) {
	url := fmt.Sprintf("%s/c/%s/faktura-vydana/%s.pdf", internal.AbraBaseURL, internal.AbraCompany, invoiceID)
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

	return b64.StdEncoding.EncodeToString([]byte(body)), nil
}
