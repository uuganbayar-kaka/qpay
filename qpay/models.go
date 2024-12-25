package qpay

import (
	"time"
)

type Invoice struct {
	ID                string    `json:"id"`
	Description       string    `json:"description,omitempty"`
	InvoiceID         string    `json:"invoice_id,omitempty"`
	InvoiceData       string    `json:"invoice_data,omitempty"`
	Amount            string    `json:"amount,omitempty"`
	ExpireDate        string    `json:"expire_date,omitempty"`
	CreatedBy         string    `json:"created_by,omitempty"`
	CreatedDate       time.Time `json:"created_date"`
	UpdatedBy         string    `json:"updated_by,omitempty"`
	UpdatedDate       time.Time `json:"updated_date,omitempty"`
	InvoiceStatusCode string    `json:"invoice_status_code,omitempty"`
	Status            string    `json:"status,omitempty"`
	ResponseInvoiceID string    `json:"response_invoice_id,omitempty"`
	QRShortURL        string    `json:"qr_short_url,omitempty"`
	QRText            string    `json:"qr_text,omitempty"`
	QRImage           string    `json:"qr_image,omitempty"`
}

func (i *Invoice) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":                  i.ID,
		"description":         i.Description,
		"invoice_id":          i.InvoiceID,
		"invoice_data":        i.InvoiceData,
		"amount":              i.Amount,
		"expire_date":         i.ExpireDate,
		"created_by":          i.CreatedBy,
		"created_date":        i.CreatedDate,
		"updated_by":          i.UpdatedBy,
		"updated_date":        i.UpdatedDate,
		"invoice_status_code": i.InvoiceStatusCode,
		"status":              i.Status,
		"response_invoice_id": i.ResponseInvoiceID,
		"qr_short_url":        i.QRShortURL,
		"qr_text":             i.QRText,
		"qr_image":            i.QRImage,
	}
}
