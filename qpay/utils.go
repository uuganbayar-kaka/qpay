package qpay

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func RequestExtract(requestBody io.ReadCloser) ([]byte, error) {
	decoder, decErr := io.ReadAll(requestBody)
	if decErr != nil {
		return nil, decErr
	}
	return decoder, nil
}

// For SuccessFull Response
func Response(resWriter http.ResponseWriter, resp ResponseBuilder) {
	resWriter.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	resWriter.WriteHeader(resp.StatusCode)
	encoder := json.NewEncoder(resWriter)
	DataRet := resp.Data
	encoder.Encode(DataRet)
}

// Response Exception Handler
func ResponseExceptionHandler(resWriter http.ResponseWriter, respExc ResponseExceptionBuilder) {
	resWriter.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	resWriter.WriteHeader(respExc.Error.StatusCode)
	encoder := json.NewEncoder(resWriter)
	encoder.Encode(respExc)
}
