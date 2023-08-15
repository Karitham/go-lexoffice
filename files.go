//********************************************************************************************************************//
//
// This file is part of golexoffice.
// All code may be used. Feel free and maybe code something better.
//
// Author: Jonas Kwiedor (aka gowizzard)
//
//********************************************************************************************************************//

package golexoffice

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
)

type CreateFileResponse struct {
	ID string `json:"id"`
}

// CreateFile uploads a file
// <https://developers.lexoffice.io/docs/?shell#files-endpoint-upload-a-file>
func (c *Client) CreateFile(ctx context.Context, r io.Reader, name string) (CreateFileResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	filePart, err := writer.CreateFormFile("file", name)
	if err != nil {
		return CreateFileResponse{}, err
	}

	_, err = io.Copy(filePart, r)
	if err != nil {
		return CreateFileResponse{}, err
	}

	_ = writer.WriteField("type", "voucher")

	err = writer.Close()
	if err != nil {
		return CreateFileResponse{}, err
	}

	var er LegacyErrorResponse
	var fr CreateFileResponse
	err = c.Request("/v1/files").
		ContentType(writer.FormDataContentType()).
		BodyReader(body).
		ToJSON(&fr).
		ErrorJSON(&er).
		Fetch(ctx)
	if err != nil {
		return fr, fmt.Errorf("error while request (%s): %w", er, err)
	}

	return fr, nil
}

// DownloadFile downloads a file
// <https://developers.lexoffice.io/docs/?shell#files-endpoint-download-a-file>
func (c *Client) DownloadFile(ctx context.Context, out io.Writer, id string) error {
	var er LegacyErrorResponse
	err := c.Requestf("/v1/files/%s", id).
		ErrorJSON(&er).
		ToWriter(out).
		Fetch(ctx)
	if err != nil {
		return fmt.Errorf("error while request (%s): %w", er, err)
	}

	return nil
}
