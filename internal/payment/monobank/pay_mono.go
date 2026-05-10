package monobank

import (
	p "Price/internal/payment"
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type MonoBank struct {
	token  string
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

type monoWebhookRequest struct {
	Id     string `json:"invoiceId"`
	Status string `json:"status"`
	Amount int    `json:"amount"`
	UserId string `json:"reference"`
}

type publicKey struct {
	PublicKey string `json:"key"`
}

func NewMonoBank(token string, client *http.Client, apiURL string) *MonoBank {
	return &MonoBank{
		token:  token,
		client: client,
		apiURL: apiURL,
	}
}

func (m *MonoBank) getPublicKey(ctx context.Context) (string, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.monobank.ua/api/merchant/pubkey", nil)
	if err != nil {
		return "", err
	}

	res, err := m.client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		byt, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("ошибка MonoBank, статус: %s текст ошибки:%s", res.Status, string(byt))
	}

	pub := publicKey{}

	if err := json.NewDecoder(res.Body).Decode(&pub); err != nil {
		return "", err
	}

	return pub.PublicKey, nil

}

func (m *MonoBank) CreateInvoice(ctx context.Context, userID int64) (string, error) {
	user := fmt.Sprintf("%d", userID)
	invoiceReq := createInvoiceRequest{
		Amount:     5000,
		UserID:     user,
		Ccy:        980,
		WebHookUrl: "пока нет сервера,будет позже и добавлю сюда адрес",
	}

	body, err := json.Marshal(invoiceReq)
	if err != nil {
		return "", err
	}

	data := bytes.NewReader(body)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.apiURL, data)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Token", m.token)

	resp, err := m.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		byt, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ошибка MonoBank, статус: %s текст ошибки:%s", resp.Status, string(byt))
	}

	invoiceResp := createInvoiceResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&invoiceResp); err != nil {
		return "", err
	}

	return invoiceResp.PageUrl, nil
}

func (m *MonoBank) ParseCallback(ctx context.Context, res []byte, bankSign string) (*p.PaymentResult, error) {
	result, err := m.getPublicKey(ctx)
	if err != nil {
		return nil, err
	}

	pubKeyBytes, err := base64.StdEncoding.DecodeString(result)
	if err != nil {
		return nil, err
	}

	bankBytes, err := base64.StdEncoding.DecodeString(bankSign)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pubKeyBytes)
	if block == nil {
		return nil, errors.New("ошибка парсинга PEM")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("ключ не являеться ECDSA")
	}

	hash := sha256.Sum256(res)

	verif := ecdsa.VerifyASN1(ecdsaPub, hash[:], bankBytes)
	if !verif {
		return nil, errors.New("неверная подпись MonoBank")
	}

	req := monoWebhookRequest{}
	if err := json.Unmarshal(res, &req); err != nil {
		return nil, err
	}

	paymentRes := p.PaymentResult{
		PaymentID: req.Id,
		UserID:    req.UserId,
		Status:    req.Status,
		Money:     req.Amount,
		BankName:  "MonoBank",
	}
	return &paymentRes, nil
}
