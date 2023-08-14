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
	"io"
	"mime/multipart"
)

// FileReturn is to decode json data
type FileReturn struct {
	Id string `json:"id"`
}

// AddFile is to upload a file
func (c *Client) AddFile(ctx context.Context, r io.Reader, name string) (FileReturn, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	filePart, err := writer.CreateFormFile("file", name)
	if err != nil {
		return FileReturn{}, err
	}

	_, err = io.Copy(filePart, r)
	if err != nil {
		return FileReturn{}, err
	}

	_ = writer.WriteField("type", "voucher")

	err = writer.Close()
	if err != nil {
		return FileReturn{}, err
	}

	var fr FileReturn
	err = c.Request("/v1/files/").
		ContentType(writer.FormDataContentType()).
		BodyReader(body).
		ToJSON(&FileReturn{}).
		Fetch(ctx)
	if err != nil {
		return fr, err
	}

	return fr, nil
}
