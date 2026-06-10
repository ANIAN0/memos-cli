package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ANIAN0/memos-cli/pkg/httpclient"
)

// Attachment represents a Memos attachment.
type Attachment struct {
	Name       string `json:"name"`
	UID        string `json:"uid"`
	Filename   string `json:"filename"`
	Size       int64  `json:"size"`
	Type       string `json:"type"`
	CreateTime string `json:"createTime"`
}

// AttachmentClient handles attachment operations.
type AttachmentClient struct {
	C *httpclient.Client
}

// normalizeAttachmentID converts "123" to "attachments/123" if needed.
func normalizeAttachmentID(id string) string {
	if strings.HasPrefix(id, "attachments/") {
		return id
	}
	return "attachments/" + id
}

// Upload uploads a file as an attachment.
func (a *AttachmentClient) Upload(ctx context.Context, filePath string) (*Attachment, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	// Create multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}

	if _, err := io.Copy(part, f); err != nil {
		return nil, fmt.Errorf("copy file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close writer: %w", err)
	}

	resp, err := a.C.Post(ctx, "/api/v1/attachments", writer.FormDataContentType(), body)
	if err != nil {
		return nil, fmt.Errorf("upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("upload failed: HTTP %d", resp.StatusCode)
	}

	var att Attachment
	if err := json.NewDecoder(resp.Body).Decode(&att); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &att, nil
}

// List returns a list of attachments.
func (a *AttachmentClient) List(ctx context.Context, pageSize int) ([]Attachment, error) {
	u := "/api/v1/attachments"
	if pageSize > 0 {
		u += fmt.Sprintf("?pageSize=%d", pageSize)
	}

	resp, err := a.C.Get(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("list request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("list failed: HTTP %d", resp.StatusCode)
	}

	// Try to decode as array first, then as wrapped object
	var atts []Attachment
	if err := json.NewDecoder(resp.Body).Decode(&atts); err != nil {
		// Reset reader and try wrapped format
		var respBody struct {
			Attachments []Attachment `json:"attachments"`
		}
		if err2 := json.NewDecoder(resp.Body).Decode(&respBody); err2 != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}
		atts = respBody.Attachments
	}
	return atts, nil
}

// Get downloads an attachment to a local file.
func (a *AttachmentClient) Get(ctx context.Context, id, outputPath string) error {
	id = normalizeAttachmentID(id)

	resp, err := a.C.Get(ctx, "/api/v1/"+id)
	if err != nil {
		return fmt.Errorf("get request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("get failed: HTTP %d", resp.StatusCode)
	}

	// Create parent directory if needed
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

// Delete deletes an attachment by ID.
func (a *AttachmentClient) Delete(ctx context.Context, id string) error {
	id = normalizeAttachmentID(id)

	req, err := http.NewRequestWithContext(ctx, "DELETE", a.C.BaseURL+"/api/v1/"+id, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := a.C.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("delete failed: HTTP %d", resp.StatusCode)
	}
	return nil
}