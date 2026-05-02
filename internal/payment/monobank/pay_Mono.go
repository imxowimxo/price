package monobank

import (
	p "Price/internal/payment"
	"context"
	"encoding/json"
	"net/http"
)

type MonoBank struct {
	token  string // сделал с маленькой буквы чтобы был недоступен в других пакетах,так лучше?
	client *http.Client
	apiURL string
}

type createInvoiceRequest struct {
	Amount     int64  `json:"amount"`
	UserID     string `json:"reference"`
	Ccy        int    `json:"ccy"`
	WebHookUrl string `json:"web_hook_url"`
}

type createInvoiceResponse struct {
	PageUrl   string `json:"page_url"`
	InvoiceId string `json:"invoice_id"`
}

func NewMonoBank(token string, client *http.Client, apiURL string) *MonoBank {
	return &MonoBank{
		token:  token,
		client: client,
		apiURL: apiURL,
	}
}

func (m *MonoBank) CreateInvoice(ctx context.Context, userID int64, typeSub string) (string, error) {
	invoiceReq := createInvoiceRequest{
		Amount:     1500,
		Ccy:        980,
		WebHookUrl: "https://pay.monobank.io",
	}

	body, err := json.Marshal(invoiceReq)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.apiURL, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Token", m.token)

	resp, err := m.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	json.Unmarshal(,resp)

}

func (m *MonoBank) ParseCallback(ctx context.Context, res []byte, bankSign string) (*p.PaymentResult, error) {
}
