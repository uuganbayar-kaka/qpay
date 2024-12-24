package qpay

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

// Constants
const (
	ContentTypeJSON = "application/json"
	QPayURL         = "https://merchant-sandbox.qpay.mn"
	QPAY_TOKEN_URL  = "https://merchant-sandbox.qpay.mn/v2/auth/token"
	QPAY_USERNAME   = "TEST_MERCHANT"
	QPAY_PASSWORD   = "123456"
)

// Cache instance (in-memory, replace with Redis or other if needed)
var appCache = cache.New(5*time.Minute, 10*time.Minute)

// Request and Response Structs
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type InvoiceRequest struct {
	InvoiceCode         string  `json:"invoice_code"`
	SenderInvoiceNo     string  `json:"sender_invoice_no"`
	InvoiceReceiverCode string  `json:"invoice_receiver_code"`
	InvoiceDescription  string  `json:"invoice_description"`
	Amount              float64 `json:"amount"`
	CallbackURL         string  `json:"callback_url"`
}

type InvoiceResponse struct {
	InvoiceID  string `json:"invoice_id"`
	QRImage    string `json:"qr_image"`
	QRShortURL string `json:"qPay_shortUrl"`
	QRText     string `json:"qr_text"`
}

// QPayGetToken retrieves a QPay access token
func QPayGetToken(invoiceNo string) (string, error) {
	terminalID := "91909029"
	username := QPAY_USERNAME
	password := QPAY_PASSWORD
	tokenURL := QPAY_TOKEN_URL

	authData := fmt.Sprintf("%s:%s", username, password)
	authB64 := base64.StdEncoding.EncodeToString([]byte(authData))

	log.Printf("QPay access token: %s-%s-%s", tokenURL, authB64, terminalID)

	requestBody := map[string]string{
		"terminal_id": terminalID,
	}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic "+authB64)
	req.Header.Set("Content-Type", ContentTypeJSON)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to retrieve token")
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	appCache.Set(invoiceNo, tokenResp, cache.DefaultExpiration)
	return tokenResp.AccessToken, nil
}

// CreateInvoiceSimple creates a simple invoice
func CreateInvoiceSimple(data map[string]interface{}) (InvoiceResponse, error) {
	amount := data["amount"].(float64)
	if amount <= 0 {
		return InvoiceResponse{}, errors.New("amount must be greater than zero")
	}

	invoiceNo := time.Now().Format("20060102150405")
	log.Printf("Invoice No: %s", invoiceNo)

	bearerToken, err := QPayGetToken(invoiceNo)
	if err != nil {
		return InvoiceResponse{}, err
	}

	payload := InvoiceRequest{
		InvoiceCode:         "TEST_INVOICE",
		SenderInvoiceNo:     invoiceNo,
		InvoiceReceiverCode: "91909029",
		InvoiceDescription:  "test",
		Amount:              amount,
		CallbackURL:         "https://example.com/payments?payment_id=91909029",
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", QPayURL+"/v2/invoice", bytes.NewBuffer(body))
	if err != nil {
		return InvoiceResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", ContentTypeJSON)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return InvoiceResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return InvoiceResponse{}, errors.New("failed to create invoice")
	}

	var invoiceResp InvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&invoiceResp); err != nil {
		return InvoiceResponse{}, err
	}

	return invoiceResp, nil
}
