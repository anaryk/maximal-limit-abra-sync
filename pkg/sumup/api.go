package sumup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/anaryk/maximal-limit-abra-sync/pkg/internal"
)

func (c *Connector) GetTransactionInSpecificDate(date, merchant string) (*SumUpTransaction, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.sumup.com/v2.1/merchants/%s/transactions/history?limit=%s&order=%s&oldest_time=%s", merchant, strconv.Itoa(internal.SumUpAPITransactionLimit), internal.SumUpAPIOrdering, date), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get transaction: %s", resp.Status)
	}

	var transaction SumUpTransaction
	if err := json.NewDecoder(resp.Body).Decode(&transaction); err != nil {
		return nil, err
	}

	return &transaction, nil
}

// GetReceiptByTransactionID retrieves the receipt for a specific transaction ID.
// exmaple https://api.sumup.com/v1.1/receipts/bf41901c-e522-4bd3-90b9-1c2210b24ee0?mid=ME4Z76XR

func (c *Connector) GetReceiptByTransactionID(transactionID, merchant string) (*SumUpReceipts, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.sumup.com/v1.1/receipts/%s?mid=%s", transactionID, merchant), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get receipt: %s", resp.Status)
	}

	var receipt SumUpReceipts
	if err := json.NewDecoder(resp.Body).Decode(&receipt); err != nil {
		return nil, err
	}

	return &receipt, nil
}
