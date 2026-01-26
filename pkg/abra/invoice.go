package abra

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/anaryk/maximal-limit-abra-sync/pkg/internal"
)

func (c *Connector) CreateInvoice(customerCode, issueDate, dueDate string, internalNumber string, items []FakturaPolozka) (*APIResponse, error) {
	contactId, err := c.GetContactIDByShortName(customerCode)
	if err != nil {
		return nil, err
	}
	invoice := InvoiceRequest{
		Winstrom: struct {
			FakturaVydana []FakturaVydana `json:"faktura-vydana"`
		}{
			FakturaVydana: []FakturaVydana{
				{
					Kod:              internalNumber,
					DatVyst:          issueDate,
					DatSplat:         dueDate,
					StavUhrady:       "stavUhr.uhrazenoRucne",
					IdFirmy:          contactId,
					Polozky:          items,
					TypFaktury:       internal.AbraVydanaFakturaID,
					AccountingType:   internal.AbraAccountingOperationID,
					FormaUhradyCislo: internal.AbraPaymentCardID,
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

func (c *Connector) AddInvoiceItem(invoiceID string, item FakturaPolozka) (*APIResponse, error) {
	// Step 1: Set invoice to unpaid state (required to modify paid invoices)
	// Use empty string to remove payment status
	fmt.Printf("Step 1: Setting invoice %s to unpaid...\n", invoiceID)
	if err := c.setInvoicePaymentStatus(invoiceID, ""); err != nil {
		return nil, fmt.Errorf("failed to set invoice %s to unpaid: %w", invoiceID, err)
	}

	// Step 2: Add the item using PUT on the invoice
	type ItemAddRequest struct {
		Winstrom struct {
			FakturaVydana struct {
				PolozkyFaktury []struct {
					Nazev       string  `json:"nazev"`
					MnozMj      float64 `json:"mnozMj"`
					CenaMj      float64 `json:"cenaMj"`
					TypCenyDphK string  `json:"typCenyDphK"`
				} `json:"polozkyFaktury"`
			} `json:"faktura-vydana"`
		} `json:"winstrom"`
	}

	request := ItemAddRequest{}
	request.Winstrom.FakturaVydana.PolozkyFaktury = []struct {
		Nazev       string  `json:"nazev"`
		MnozMj      float64 `json:"mnozMj"`
		CenaMj      float64 `json:"cenaMj"`
		TypCenyDphK string  `json:"typCenyDphK"`
	}{
		{
			Nazev:       item.Popis,
			MnozMj:      item.Pocet,
			CenaMj:      item.CenaKus,
			TypCenyDphK: "typCeny.sDph",
		},
	}

	payload, err := json.Marshal(request)
	if err != nil {
		_ = c.setInvoicePaymentStatus(invoiceID, "stavUhr.uhrazenoRucne")
		return nil, err
	}

	fmt.Printf("Step 2: Adding item to invoice %s, payload: %s\n", invoiceID, string(payload))

	url := fmt.Sprintf("%s/c/%s/faktura-vydana/code:%s.json", internal.AbraBaseURL, internal.AbraCompany, invoiceID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		_ = c.setInvoicePaymentStatus(invoiceID, "stavUhr.uhrazenoRucne")
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		_ = c.setInvoicePaymentStatus(invoiceID, "stavUhr.uhrazenoRucne")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		_ = c.setInvoicePaymentStatus(invoiceID, "stavUhr.uhrazenoRucne")
		return nil, err
	}

	fmt.Printf("Step 2: AddInvoiceItem PUT response (HTTP %d): %s\n", resp.StatusCode, string(body))

	var apiResponse APIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		_ = c.setInvoicePaymentStatus(invoiceID, "stavUhr.uhrazenoRucne")
		return nil, fmt.Errorf("failed to unmarshal response: %w, body: %s", err, string(body))
	}

	// Step 3: Set invoice back to paid state
	fmt.Printf("Step 3: Setting invoice %s back to paid...\n", invoiceID)
	if err := c.setInvoicePaymentStatus(invoiceID, "stavUhr.uhrazenoRucne"); err != nil {
		fmt.Printf("Warning: Failed to set invoice %s back to paid: %v\n", invoiceID, err)
	}

	return &apiResponse, nil
}

func (c *Connector) setInvoicePaymentStatus(invoiceID string, status string) error {
	type StatusRequest struct {
		Winstrom struct {
			FakturaVydana struct {
				StavUhrK string `json:"stavUhrK"`
			} `json:"faktura-vydana"`
		} `json:"winstrom"`
	}

	request := StatusRequest{}
	request.Winstrom.FakturaVydana.StavUhrK = status

	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/c/%s/faktura-vydana/code:%s.json", internal.AbraBaseURL, internal.AbraCompany, invoiceID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("setInvoicePaymentStatus(%s, %q) response (HTTP %d): %s\n", invoiceID, status, resp.StatusCode, string(body))

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("failed to update invoice status (HTTP %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Connector) GetPDFInvoiceAsBase64(invoiceID string) (string, error) {
	url := fmt.Sprintf("%s/c/%s/faktura-vydana/%s.pdf?report-name=fakturaConfig", internal.AbraBaseURL, internal.AbraCompany, invoiceID)
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

// InvoiceItemResponse represents the response when fetching invoice items
type InvoiceItemResponse struct {
	Winstrom struct {
		FakturaVydanaPolozka []struct {
			ID       string `json:"id"`
			Nazev    string `json:"nazev"`
			SumCelkem string `json:"sumCelkem"`
		} `json:"faktura-vydana-polozka"`
	} `json:"winstrom"`
}

// VoucherItemInfo contains info about a voucher item on an invoice
type VoucherItemInfo struct {
	ID     string
	Nazev  string
	Amount float64
}

// GetInvoiceVoucherItems returns IDs of voucher discount items on an invoice
func (c *Connector) GetInvoiceVoucherItems(invoiceCode string) ([]string, error) {
	items, err := c.GetInvoiceVoucherItemsWithAmount(invoiceCode)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids, nil
}

// GetInvoiceVoucherItemsWithAmount returns voucher items with their amounts
func (c *Connector) GetInvoiceVoucherItemsWithAmount(invoiceCode string) ([]VoucherItemInfo, error) {
	url := fmt.Sprintf("%s/c/%s/faktura-vydana/code:%s.json?detail=full", internal.AbraBaseURL, internal.AbraCompany, invoiceCode)
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

	type InvoiceDetailResponse struct {
		Winstrom struct {
			FakturaVydana []struct {
				PolozkyFaktury []struct {
					ID        string  `json:"id"`
					Nazev     string  `json:"nazev"`
					SumCelkem float64 `json:"sumCelkem"`
				} `json:"polozkyFaktury"`
			} `json:"faktura-vydana"`
		} `json:"winstrom"`
	}

	var invoiceResp InvoiceDetailResponse
	if err := json.Unmarshal(body, &invoiceResp); err != nil {
		return nil, fmt.Errorf("failed to parse invoice response: %w", err)
	}

	var voucherItems []VoucherItemInfo
	if len(invoiceResp.Winstrom.FakturaVydana) > 0 {
		for _, item := range invoiceResp.Winstrom.FakturaVydana[0].PolozkyFaktury {
			if strings.Contains(item.Nazev, "Slevový poukaz") {
				voucherItems = append(voucherItems, VoucherItemInfo{
					ID:     item.ID,
					Nazev:  item.Nazev,
					Amount: item.SumCelkem,
				})
			}
		}
	}

	return voucherItems, nil
}

// HasVoucherItem checks if an invoice already has a voucher discount item
func (c *Connector) HasVoucherItem(invoiceCode string) (bool, error) {
	items, err := c.GetInvoiceVoucherItems(invoiceCode)
	if err != nil {
		return false, err
	}
	return len(items) > 0, nil
}

// RemoveInvoiceItem removes an item from an invoice by item ID using @action delete
func (c *Connector) RemoveInvoiceItem(invoiceCode string, itemID string) error {
	// Step 1: Set invoice to unpaid
	if err := c.setInvoicePaymentStatus(invoiceCode, ""); err != nil {
		return fmt.Errorf("failed to set invoice to unpaid: %w", err)
	}

	// Step 2: Delete the item using PUT with @action: delete
	type DeleteItemRequest struct {
		Winstrom struct {
			FakturaVydana struct {
				PolozkyFaktury []struct {
					ID     string `json:"id"`
					Action string `json:"@action"`
				} `json:"polozkyFaktury"`
			} `json:"faktura-vydana"`
		} `json:"winstrom"`
	}

	request := DeleteItemRequest{}
	request.Winstrom.FakturaVydana.PolozkyFaktury = []struct {
		ID     string `json:"id"`
		Action string `json:"@action"`
	}{
		{ID: itemID, Action: "delete"},
	}

	payload, err := json.Marshal(request)
	if err != nil {
		_ = c.setInvoicePaymentStatus(invoiceCode, "stavUhr.uhrazenoRucne")
		return err
	}

	url := fmt.Sprintf("%s/c/%s/faktura-vydana/code:%s.json", internal.AbraBaseURL, internal.AbraCompany, invoiceCode)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		_ = c.setInvoicePaymentStatus(invoiceCode, "stavUhr.uhrazenoRucne")
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		_ = c.setInvoicePaymentStatus(invoiceCode, "stavUhr.uhrazenoRucne")
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("RemoveInvoiceItem response (HTTP %d): %s\n", resp.StatusCode, string(body))

	// Step 3: Set invoice back to paid
	if err := c.setInvoicePaymentStatus(invoiceCode, "stavUhr.uhrazenoRucne"); err != nil {
		fmt.Printf("Warning: Failed to set invoice back to paid: %v\n", err)
	}

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if apiResponse.Winstrom.Success != "true" {
		return fmt.Errorf("failed to delete item: %v", apiResponse.Winstrom.Results)
	}

	return nil
}

// RemoveDuplicateVoucherItems removes all but one voucher item from an invoice
func (c *Connector) RemoveDuplicateVoucherItems(invoiceCode string) (int, error) {
	items, err := c.GetInvoiceVoucherItems(invoiceCode)
	if err != nil {
		return 0, err
	}

	if len(items) <= 1 {
		return 0, nil // No duplicates
	}

	// Keep the first one, remove the rest
	removedCount := 0
	for i := 1; i < len(items); i++ {
		if err := c.RemoveInvoiceItem(invoiceCode, items[i]); err != nil {
			fmt.Printf("Warning: Failed to remove duplicate item %s: %v\n", items[i], err)
			continue
		}
		removedCount++
	}

	return removedCount, nil
}
