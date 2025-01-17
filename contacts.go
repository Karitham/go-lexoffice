//**********************************************************
//
// This file is part of lexoffice.
// All code may be used. Feel free and maybe code something better.
//
// Author: Jonas Kwiedor
//
//**********************************************************

package golexoffice

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aarondl/opt/omit"
)

// ContactsReturn is to decode json data
type ContactsReturn struct {
	Content          []ContactsContent    `json:"content"`
	First            bool                 `json:"first"`
	Last             bool                 `json:"last"`
	TotalPages       int                  `json:"totalPages"`
	TotalElements    int                  `json:"totalElements"`
	NumberOfElements int                  `json:"numberOfElements"`
	Size             int                  `json:"size"`
	Number           int                  `json:"number"`
	Sort             []ContactsReturnSort `json:"sort"`
}

type ContactsContent struct {
	Id             string                    `json:"id,omitempty"`
	Version        int                       `json:"version,omitempty"`
	Roles          ContactBodyRoles          `json:"roles"`
	Company        *ContactBodyCompany       `json:"company,omitempty"`
	Person         *ContactBodyPerson        `json:"person,omitempty"`
	Addresses      ContactBodyAddresses      `json:"addresses"`
	EmailAddresses ContactBodyEmailAddresses `json:"emailAddresses"`
	PhoneNumbers   ContactBodyPhoneNumbers   `json:"phoneNumbers"`
	Note           string                    `json:"note"`
	Archived       bool                      `json:"archived,omitempty"`
}

type ContactsReturnRoles struct {
	Customer ContactsReturnCustomer `json:"customer"`
	Vendor   ContactsReturnVendor   `json:"vendor"`
}

type ContactsReturnCustomer struct {
	Number int `json:"number,omitempty"`
}

type ContactsReturnVendor struct {
	Number int `json:"number,omitempty"`
}

type ContactsReturnCompany struct {
	Name                 string                         `json:"name"`
	TaxNumber            string                         `json:"taxNumber"`
	VatRegistrationId    string                         `json:"vatRegistrationId"`
	AllowTaxFreeInvoices bool                           `json:"allowTaxFreeInvoices"`
	ContactPersons       []ContactsReturnContactPersons `json:"contactPersons"`
}

type ContactsReturnContactPersons struct {
	Salutation   string `json:"salutation"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	EmailAddress string `json:"emailAddress"`
	PhoneNumber  string `json:"phoneNumber"`
}

type ContactsReturnAddresses struct {
	Billing  []ContactsReturnBilling  `json:"billing"`
	Shipping []ContactsReturnShipping `json:"shipping"`
}

type ContactsReturnBilling struct {
	Supplement  string `json:"supplement"`
	Street      string `json:"street"`
	Zip         string `json:"zip"`
	City        string `json:"city"`
	CountryCode string `json:"countryCode"`
}

type ContactsReturnShipping struct {
	Supplement  string `json:"supplement"`
	Street      string `json:"street"`
	Zip         string `json:"zip"`
	City        string `json:"city"`
	CountryCode string `json:"countryCode"`
}

type ContactsReturnEmailAddresses struct {
	Business []string `json:"business"`
	Office   []string `json:"office"`
	Private  []string `json:"private"`
	Other    []string `json:"other"`
}

type ContactsReturnPhoneNumbers struct {
	Business []string `json:"business"`
	Office   []string `json:"office"`
	Mobile   []string `json:"mobile"`
	Private  []string `json:"private"`
	Fax      []string `json:"fax"`
	Other    []string `json:"other"`
}

type ContactsReturnSort struct {
	Property     string `json:"property"`
	Direction    string `json:"direction"`
	IgnoreCase   bool   `json:"ignoreCase"`
	NullHandling string `json:"nullHandling"`
	Ascending    bool   `json:"ascending"`
}

// ContactBody is to create a new contact
type ContactBody struct {
	Id             string                     `json:"id,omitempty"`
	Version        int                        `json:"version"`
	Roles          ContactBodyRoles           `json:"roles"`
	Company        *ContactBodyCompany        `json:"company,omitempty"`
	Person         *ContactBodyPerson         `json:"person,omitempty"`
	Addresses      *ContactBodyAddresses      `json:"addresses,omitempty"`
	EmailAddresses *ContactBodyEmailAddresses `json:"emailAddresses,omitempty"`
	PhoneNumbers   *ContactBodyPhoneNumbers   `json:"phoneNumbers,omitempty"`
	Note           string                     `json:"note"`
	Archived       bool                       `json:"archived,omitempty"`
}

type ContactBodyRoles struct {
	Customer *ContactBodyCustomer `json:"customer,omitempty"`
	Vendor   *ContactBodyVendor   `json:"vendor,omitempty"`
}

type ContactBodyCustomer struct {
	Number int `json:"number,omitempty"`
}

type ContactBodyVendor struct {
	Number int `json:"number,omitempty"`
}

type ContactBodyCompany struct {
	Name                 string                       `json:"name"`
	TaxNumber            string                       `json:"taxNumber,omitempty"`
	VatRegistrationId    string                       `json:"vatRegistrationId,omitempty"`
	AllowTaxFreeInvoices bool                         `json:"allowTaxFreeInvoices"`
	ContactPersons       []*ContactBodyContactPersons `json:"contactPersons"`
}

type ContactBodyPerson struct {
	Salutation string `json:"salutation"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
}

type ContactBodyContactPersons struct {
	Salutation   string `json:"salutation"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	EmailAddress string `json:"emailAddress"`
	PhoneNumber  string `json:"phoneNumber"`
}

type ContactBodyAddresses struct {
	Billing  []*ContactBodyBilling  `json:"billing"`
	Shipping []*ContactBodyShipping `json:"shipping"`
}

type ContactBodyBilling struct {
	Supplement  string `json:"supplement"`
	Street      string `json:"street"`
	Zip         string `json:"zip"`
	City        string `json:"city"`
	CountryCode string `json:"countryCode"`
}

type ContactBodyShipping struct {
	Supplement  string `json:"supplement"`
	Street      string `json:"street"`
	Zip         string `json:"zip"`
	City        string `json:"city"`
	CountryCode string `json:"countryCode"`
}

type ContactBodyEmailAddresses struct {
	Business []string `json:"business"`
	Office   []string `json:"office"`
	Private  []string `json:"private"`
	Other    []string `json:"other"`
}

type ContactBodyPhoneNumbers struct {
	Business []string `json:"business"`
	Office   []string `json:"office"`
	Mobile   []string `json:"mobile"`
	Private  []string `json:"private"`
	Fax      []string `json:"fax"`
	Other    []string `json:"other"`
}

type ContactsResponse struct {
	ID          string `json:"id"`
	ResourceUri string `json:"resourceUri"`
	CreatedDate string `json:"createdDate"`
	UpdatedDate string `json:"updatedDate"`
	Version     int    `json:"version"`
}

type GetContactsParams struct {
	Page   omit.Val[int]
	Email  omit.Val[string]
	Name   omit.Val[string]
	Number omit.Val[int]

	// filtering goes like this:
	// 	- unspecified -> no filter
	// 	- true -> only customer
	// 	- false -> only non-customer
	Customer omit.Val[bool]

	// filtering goes like this:
	// 	- unspecified -> no filter
	// 	- true -> only vendor
	// 	- false -> only non-vendor
	Vendor omit.Val[bool]
}

// GetContacts is to get a list of all contacts
// <https://developers.lexoffice.io/docs/?shell#contacts-endpoint-filtering-contacts>
func (c *Client) GetContacts(ctx context.Context, p GetContactsParams) (ContactsReturn, error) {
	var er LegacyErrorResponse
	var cr ContactsReturn

	qb := c.Request("/v1/contacts").
		ToJSON(&cr).
		ErrorJSON(&er)

	if p.Page.IsSet() {
		qb.ParamInt("page", p.Page.MustGet())
	}

	if p.Page.IsSet() {
		qb.Param("email", p.Email.MustGet())
	}

	if p.Page.IsSet() {
		qb.Param("name", p.Name.MustGet())
	}

	if p.Page.IsSet() {
		qb.ParamInt("number", p.Number.MustGet())
	}

	if p.Customer.IsSet() {
		qb.Param("customer", strconv.FormatBool(p.Customer.GetOrZero()))
	}

	if p.Vendor.IsSet() {
		qb.Param("vendor", strconv.FormatBool(p.Vendor.GetOrZero()))
	}

	err := qb.Fetch(ctx)
	if err != nil {
		return ContactsReturn{}, fmt.Errorf("error getting contacts (%s): %w", er.String(), err)
	}

	return cr, nil

}

// GetContact is to get a contact by id
// <https://developers.lexoffice.io/docs/?shell#contacts-endpoint-retrieve-a-contact>
func (c *Client) GetContact(ctx context.Context, id string) (ContactsContent, error) {
	var er LegacyErrorResponse
	var crc ContactsContent
	err := c.Requestf("/v1/contacts/%s", id).ToJSON(&crc).ErrorJSON(&er).Fetch(ctx)
	if err != nil {
		return crc, fmt.Errorf("error getting contact (%s): %w", er.String(), err)
	}
	return crc, nil

}

// CreateContact creates a new contact
// <https://developers.lexoffice.io/docs/?shell#contacts-endpoint-create-a-contact>
func (c *Client) CreateContact(ctx context.Context, body ContactBody) (ContactsResponse, error) {
	var er LegacyErrorResponse
	var cr ContactsResponse
	err := c.Request("/v1/contacts").
		BodyJSON(body).
		ToJSON(&cr).
		ErrorJSON(&er).
		Post().
		Fetch(ctx)
	if err != nil {
		return cr, fmt.Errorf("error creating contacts (%s): %w", er.String(), err)
	}

	return cr, nil
}

// UpdateContact updates existing contact
// <https://developers.lexoffice.io/docs/?shell#contacts-endpoint-update-a-contact>
func (c *Client) UpdateContact(ctx context.Context, body ContactBody) (ContactsResponse, error) {
	var er LegacyErrorResponse
	var cr ContactsResponse
	err := c.Requestf("/v1/contacts/%s", body.Id).
		BodyJSON(body).
		ToJSON(&cr).
		ErrorJSON(&er).
		Put().
		Fetch(ctx)
	if err != nil {
		return cr, fmt.Errorf("error updating contacts (%s): %w", er.String(), err)
	}

	return cr, nil
}
