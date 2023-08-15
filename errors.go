package golexoffice

import "strings"

// source: https://developers.lexoffice.io/docs/#error-codes-legacy-error-response
// files, profile, contacts
//
//	{
//		"requestId":"3fb21ee4-ad26-4e2f-82af-a1197af02d08",
//		"IssueList":[
//			{"i18nKey":"invalid_value","source":"company and person","type":"validation_failure"},
//			{"i18nKey":"missing_entity","source":"company.name","type":"validation_failure"}
//		]
//	}
type LegacyErrorResponse struct {
	RequestID string `json:"requestId"`
	IssueList []struct {
		Key    string `json:"i18nKey"`
		Source string `json:"source"`
		Type   string `json:"type"`
	} `json:"IssueList"`
}

func (e LegacyErrorResponse) Error() string {
	return e.String()
}

func (e LegacyErrorResponse) String() string {
	builder := &strings.Builder{}
	for i, issue := range e.IssueList {
		builder.WriteString(issue.Key)
		builder.WriteString(": ")
		builder.WriteString(issue.Source)
		if issue.Type != "" {
			builder.WriteString(" (")
			builder.WriteString(issue.Type)
			builder.WriteString(")")
		}

		if i < len(e.IssueList)-1 {
			builder.WriteString(", ")
		}
	}

	return builder.String()
}

// source: https://developers.lexoffice.io/docs/#error-codes-regular-error-response
// event-subscription, invoices
//
//	{
//		"timestamp": "2017-05-11T17:12:31.233+02:00",
//		"status": 406,
//		"error": "Not Acceptable",
//		"path": "/v1/invoices",
//		"traceId": "90d78d0777be",
//		"message": "Validation failed for request. Please see details list for specific causes.",
//		"details": [
//			{
//				"violation": "NOTNULL",
//				"field": "lineItems[0].unitPrice.taxRatePercentage",
//				"message": "darf nicht leer sein"
//			}
//		]
//	}
type ErrorResponse struct {
	Timestamp   Date   `json:"timestamp"`
	Status      int    `json:"status"`
	ErrorString string `json:"error"`
	Path        string `json:"path"`
	TraceID     string `json:"traceId"`
	Message     string `json:"message"`
	Details     []struct {
		Violation string `json:"violation"`
		Field     string `json:"field"`
		Message   string `json:"message"`
	} `json:"details"`
}

func (e ErrorResponse) Error() string {
	return e.String()
}

func (e ErrorResponse) String() string {
	builder := &strings.Builder{}
	builder.WriteString(e.Message)

	if len(e.Details) == 0 {
		return builder.String()
	}

	builder.WriteString(" (")
	for i, detail := range e.Details {
		builder.WriteString(detail.Field)
		builder.WriteString(": ")
		builder.WriteString(detail.Message)
		if detail.Violation != "" {
			builder.WriteString(" (")
			builder.WriteString(detail.Violation)
			builder.WriteString(")")
		}

		if i < len(e.Details)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString(")")

	return builder.String()
}
