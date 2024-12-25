package qpay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// CreateInvoiceSimple creates a simple invoice
func CreateInvoiceSimple(res http.ResponseWriter, req *http.Request) {

	var reqDta InvoiceSimpleData
	var responseException ResponseExceptionBuilder
	var response ResponseBuilder

	bodyByte, bodyErr := RequestExtract(req.Body)
	if bodyErr != nil {
		fmt.Printf("Request body extraction error [%s]\n", bodyErr.Error())
		responseException.Error.Detail = JSONInvalid
		responseException.Error.StatusCode = StatusBadRequest
		responseException.Error.Data = "JSON Error"
		ResponseExceptionHandler(res, responseException)
		return
	}
	err := json.Unmarshal(bodyByte, &reqDta)
	if err != nil {
		responseException.Error.Detail = InternalServerError
		responseException.Error.StatusCode = StatusInternalServerError
		responseException.Error.Data = "JSON Error"
		ResponseExceptionHandler(res, responseException)
		return
	}

	var amount float64
	strAmount := strings.TrimSpace(reqDta.Amount)
	amount, intErr := strconv.ParseFloat(strAmount, 64)
	if intErr != nil {
		return
	}
	if amount <= 0 {
		return
	}

	invoiceNo := time.Now().Format("20060102150405")
	log.Printf("Invoice No: %s", invoiceNo)

	bearerToken, _ := QPayGetToken(invoiceNo)

	payload := InvoiceRequest{
		InvoiceCode:         "TEST_INVOICE",
		SenderInvoiceNo:     invoiceNo,
		InvoiceReceiverCode: "91909029",
		InvoiceDescription:  "test",
		Amount:              amount,
		CallbackURL:         "https://example.com/payments?payment_id=91909029",
	}
	body, _ := json.Marshal(payload)

	request, err := http.NewRequest("POST", QPayURL+"/v2/invoice", bytes.NewBuffer(body))
	if err != nil {
		return
	}

	request.Header.Set("Authorization", "Bearer "+bearerToken)
	request.Header.Set("Content-Type", ContentTypeJSON)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var invoiceResp InvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&invoiceResp); err != nil {
		return
	}

	response.Data = invoiceResp
	response.StatusCode = StatusOK
	Response(res, response)
}
