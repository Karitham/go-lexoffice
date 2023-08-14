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
	"time"
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
	Currency          string  `json:"currency,omitempty"`
	NetAmount         float64 `json:"netAmount,omitempty"`
	GrossAmount       float64 `json:"grossAmount,omitempty"`
	TaxRatePercentage float64 `json:"taxRatePercentage,omitempty"`
}

type InvoiceBodyTotalPrice struct {
	Currency                string  `json:"currency,omitempty"`
	TotalNetAmount          float64 `json:"totalNetAmount,omitempty"`
	TotalGrossAmount        float64 `json:"totalGrossAmount,omitempty"`
	TaxRatePercentage       float64 `json:"taxRatePercentage,omitempty"`
	TotalTaxAmount          float64 `json:"totalTaxAmount,omitempty"`
	TotalDiscountAbsolute   float64 `json:"totalDiscountAbsolute,omitempty"`
	TotalDiscountPercentage float64 `json:"totalDiscountPercentage,omitempty"`
}

type InvoiceBodyTaxAmounts struct {
	TaxRatePercentage float64 `json:"taxRatePercentage,omitempty"`
	TaxAmount         float64 `json:"taxAmount,omitempty"`
	Amount            float64 `json:"amount,omitempty"`
}

type InvoiceBodyTaxConditions struct {
	TaxType     string `json:"taxType,omitempty"`
	TaxTypeNote string `json:"taxTypeNote,omitempty"`
}

type InvoiceBodyPaymentConditions struct {
	PaymentTermLabel          string                               `json:"paymentTermLabel,omitempty"`
	PaymentTermDuration       int                                  `json:"paymentTermDuration,omitempty"`
	PaymentDiscountConditions InvoiceBodyPaymentDiscountConditions `json:"paymentDiscountConditions"`
}

type InvoiceBodyPaymentDiscountConditions struct {
	DiscountPercentage int `json:"discountPercentage"`
	DiscountRange      int `json:"discountRange"`
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

func (c *Client) GetInvoice(ctx context.Context, id string) (InvoiceBody, error) {
	var ib InvoiceBody
	var er ErrorResponse
	err := c.Request("/v1/invoices/" + id).ToJSON(&ib).ErrorJSON(&er).Fetch(ctx)
	if err != nil {
		return ib, fmt.Errorf("error getting invoice (%s): %w", er.String(), err)
	}

	return ib, nil
}

func (c *Client) CreateInvoice(ctx context.Context, body InvoiceBody, finalize bool) (InvoiceReturn, error) {
	var ir InvoiceReturn
	var er ErrorResponse
	qb := c.Request("/v1/invoices").BodyJSON(body).ToJSON(&ir).Post().ErrorJSON(&er)
	if finalize {
		qb.Param("finalize", "true")
	}

	err := qb.Fetch(ctx)
	if err != nil {
		return ir, fmt.Errorf("error creating invoice (%s): %w", er.String(), err)
	}

	return ir, nil
}
