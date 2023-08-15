//**********************************************************
//
// Copyright (C) 2018 - 2023 J&J Ideenschmiede GmbH <info@jj-ideenschmiede.de>
//
// This file is part of golexoffice.
// All code may be used. Feel free and maybe code something better.
//
// Author: Jonas Kwiedor
//
//**********************************************************

package golexoffice

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/shopspring/decimal"
)

const DateFormat = "2006-01-02T15:04:05.000-07:00"

type Date time.Time

func (t *Date) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, "\"")
	tt, err := time.Parse(DateFormat, string(b))
	if err != nil {
		return err
	}
	*t = Date(tt)
	return nil
}

func (t Date) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "\"%s\"", time.Time(t).Format(DateFormat)), nil
}

// InvoiceBody is to define body data
type InvoiceBody struct {
	ID                 string                        `json:"id,omitempty"`
	OrganizationID     string                        `json:"organizationId,omitempty"`
	CreateDate         Date                          `json:"createDate,omitempty"`
	UpdatedDate        Date                          `json:"updatedDate,omitempty"`
	Version            int                           `json:"version,omitempty"`
	Archived           bool                          `json:"archived,omitempty"`
	VoucherStatus      string                        `json:"voucherStatus,omitempty"`
	VoucherNumber      string                        `json:"voucherNumber,omitempty"`
	VoucherDate        Date                          `json:"voucherDate,omitempty"`
	DueDate            Date                          `json:"dueDate,omitempty"`
	Address            InvoiceBodyAddress            `json:"address,omitempty"`
	LineItems          []InvoiceBodyLineItems        `json:"lineItems,omitempty"`
	TotalPrice         InvoiceBodyTotalPrice         `json:"totalPrice,omitempty"`
	TaxAmounts         []InvoiceBodyTaxAmounts       `json:"taxAmounts,omitempty"`
	TaxConditions      InvoiceBodyTaxConditions      `json:"taxConditions,omitempty"`
	PaymentConditions  InvoiceBodyPaymentConditions  `json:"paymentConditions,omitempty"`
	ShippingConditions InvoiceBodyShippingConditions `json:"shippingConditions,omitempty"`
	Title              string                        `json:"title,omitempty"`
	Introduction       string                        `json:"introduction,omitempty"`
	Remark             string                        `json:"remark,omitempty"`
}

type InvoiceBodyAddress struct {
	ContactID   string `json:"contactId,omitempty"`
	Name        string `json:"name,omitempty"`
	Supplement  string `json:"supplement,omitempty"`
	Street      string `json:"street,omitempty"`
	City        string `json:"city,omitempty"`
	Zip         string `json:"zip,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
}

type InvoiceBodyLineItems struct {
	Id                 string               `json:"id,omitempty"`
	Type               string               `json:"type,omitempty"`
	Name               string               `json:"name,omitempty"`
	Description        string               `json:"description,omitempty"`
	Quantity           float64              `json:"quantity,omitempty"`
	UnitName           string               `json:"unitName,omitempty"`
	UnitPrice          InvoiceBodyUnitPrice `json:"unitPrice,omitempty"`
	DiscountPercentage float64              `json:"discountPercentage,omitempty"`
	LineItemAmount     float64              `json:"lineItemAmount,omitempty"`
}

type InvoiceBodyUnitPrice struct {
	Currency          string          `json:"currency,omitempty"`
	NetAmount         decimal.Decimal `json:"netAmount,omitempty"`
	GrossAmount       decimal.Decimal `json:"grossAmount,omitempty"`
	TaxRatePercentage float64         `json:"taxRatePercentage,omitempty"`
}

type InvoiceBodyTotalPrice struct {
	Currency                string          `json:"currency,omitempty"`
	TotalNetAmount          decimal.Decimal `json:"totalNetAmount,omitempty"`
	TotalGrossAmount        decimal.Decimal `json:"totalGrossAmount,omitempty"`
	TaxRatePercentage       float64         `json:"taxRatePercentage,omitempty"`
	TotalTaxAmount          decimal.Decimal `json:"totalTaxAmount,omitempty"`
	TotalDiscountAbsolute   decimal.Decimal `json:"totalDiscountAbsolute,omitempty"`
	TotalDiscountPercentage float64         `json:"totalDiscountPercentage,omitempty"`
}

type InvoiceBodyTaxAmounts struct {
	TaxRatePercentage float64         `json:"taxRatePercentage,omitempty"`
	TaxAmount         decimal.Decimal `json:"taxAmount,omitempty"`
	Amount            decimal.Decimal `json:"amount,omitempty"`
}

type InvoiceBodyTaxConditions struct {
	TaxType     string `json:"taxType,omitempty"`
	TaxTypeNote string `json:"taxTypeNote,omitempty"`
}

type InvoiceBodyPaymentConditions struct {
	PaymentTermLabel          string                                         `json:"paymentTermLabel,omitempty"`
	PaymentTermDuration       int                                            `json:"paymentTermDuration,omitempty"`
	PaymentDiscountConditions omit.Val[InvoiceBodyPaymentDiscountConditions] `json:"paymentDiscountConditions,omitempty"`
}

type InvoiceBodyPaymentDiscountConditions struct {
	DiscountPercentage int `json:"discountPercentage,omitempty"`
	DiscountRange      int `json:"discountRange,omitempty"`
}

type InvoiceBodyShippingConditions struct {
	ShippingDate    Date   `json:"shippingDate,omitempty"`
	ShippingEndDate Date   `json:"shippingEndDate,omitempty"`
	ShippingType    string `json:"shippingType,omitempty"`
}

// InvoiceReturn is to decode json data
type InvoiceReturn struct {
	ID          string `json:"id,omitempty"`
	ResourceURI string `json:"resourceUri,omitempty"`
	CreatedDate Date   `json:"createdDate,omitempty"`
	UpdatedDate Date   `json:"updatedDate,omitempty"`
	Version     int    `json:"version,omitempty"`
}

// GetInvoice is to get a invoice
// <https://developers.lexoffice.io/docs/?shell#invoices-endpoint-retrieve-an-invoice>
func (c *Client) GetInvoice(ctx context.Context, id string) (InvoiceBody, error) {
	var ib InvoiceBody
	var er ErrorResponse
	err := c.Requestf("/v1/invoices/%s", id).ToJSON(&ib).ErrorJSON(&er).Fetch(ctx)
	if err != nil {
		return ib, fmt.Errorf("error getting invoice (%s): %w", er.String(), err)
	}

	return ib, nil
}

// CreateInvoiceOptions represent the set of possible options when creating an invoice.
// if you provide a body, then the invoice will be created with the given body.
// if you provide a preceding sales voucher id,
// then the invoice will be created from the sales voucher with the given id.
type CreateInvoiceOptions struct {
	Finalize                bool
	PrecedingSalesVoucherID string
	Body                    InvoiceBody
}

// CreateInvoice is to create a new invoice, or to pursue a sales voucher to an invoice
// <https://developers.lexoffice.io/docs/?shell#invoices-endpoint-create-an-invoice> and
// <https://developers.lexoffice.io/docs/?shell#invoices-endpoint-pursue-to-an-invoice>
func (c *Client) CreateInvoice(ctx context.Context, o CreateInvoiceOptions) (InvoiceReturn, error) {
	var ir InvoiceReturn
	var er ErrorResponse
	qb := c.Request("/v1/invoices").ToJSON(&ir).Post().ErrorJSON(&er)
	if o.Finalize {
		qb = qb.Param("finalize", "true")
	}

	if o.PrecedingSalesVoucherID != "" {
		qb = qb.Param("precedingSalesVoucherId", o.PrecedingSalesVoucherID)
	} else {
		qb = qb.BodyJSON(o.Body)
	}

	err := qb.Fetch(ctx)
	if err != nil {
		return ir, fmt.Errorf("error creating invoice (%s): %w", er.String(), err)
	}

	return ir, nil
}

type RenderResponse struct {
	ID string `json:"documentFileId,omitempty"`
}

// RenderInvoicePDF is to render a invoice as pdf
// <https://developers.lexoffice.io/docs/?shell#invoices-endpoint-render-an-invoice-document-pdf>
func (c *Client) RenderInvoicePDF(ctx context.Context, invoiceID string) (RenderResponse, error) {
	var df RenderResponse
	var er ErrorResponse
	err := c.Requestf("/v1/invoices/%s/document", invoiceID).ToJSON(&df).ErrorJSON(&er).Fetch(ctx)
	if err != nil {
		return RenderResponse{}, fmt.Errorf("error getting document file id (%s): %w", er.String(), err)
	}

	return RenderResponse{}, nil
}

// DeeplinkInvoiceURL is to get the deeplink url for a invoice
// <https://developers.lexoffice.io/docs/?shell#invoices-endpoint-deeplink-to-an-invoice>
func (c *Client) DeeplinkInvoiceURL(ctx context.Context, invoiceID string, edit bool) (string, error) {
	arg := "view"
	if edit {
		arg = "edit"
	}

	p, _ := url.JoinPath(c.baseUrl, "/permalink/invoices", arg, invoiceID)
	return p, nil
}
