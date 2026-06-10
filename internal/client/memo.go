package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ANIAN0/memos-cli/pkg/httpclient"
)

// Memo represents a Memos memo.
type Memo struct {
	Name       string   `json:"name"`       // "memos/123"
	UID        string   `json:"uid"`
	Content    string   `json:"content"`
	Visibility string   `json:"visibility"` // PRIVATE / PROTECTED / PUBLIC
	Pinned     bool     `json:"pinned"`
	Tags       []string `json:"tags"`
	Resources  []any    `json:"resources"`
	Creator    string   `json:"creator"`
	CreateTime string   `json:"createTime"`
	UpdateTime string   `json:"updateTime"`
}

// MemoClient handles memo operations.
type MemoClient struct {
	C *httpclient.Client
}

// createMemoRequest is the request body for creating/updating a memo.
type createMemoRequest struct {
	Content    string   `json:"content"`
	Visibility string   `json:"visibility,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}

// normalizeID converts "123" to "memos/123" if needed.
func normalizeID(id string) string {
	if strings.HasPrefix(id, "memos/") {
		return id
	}
	return "memos/" + id
}

// Create creates a new memo.
func (m *MemoClient) Create(ctx context.Context, content, visibility string, tags []string) (*Memo, error) {
	body, err := json.Marshal(createMemoRequest{
		Content:    content,
		Visibility: visibility,
		Tags:       tags,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := m.C.Post(ctx, "/api/v1/memos", "application/json", strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("create failed: HTTP %d", resp.StatusCode)
	}

	var memo Memo
	if err := json.NewDecoder(resp.Body).Decode(&memo); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &memo, nil
}

// Get retrieves a memo by ID.
func (m *MemoClient) Get(ctx context.Context, id string) (*Memo, error) {
	resp, err := m.C.Get(ctx, "/api/v1/"+normalizeID(id))
	if err != nil {
		return nil, fmt.Errorf("get request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get failed: HTTP %d", resp.StatusCode)
	}

	var memo Memo
	if err := json.NewDecoder(resp.Body).Decode(&memo); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &memo, nil
}

// List returns a list of memos.
func (m *MemoClient) List(ctx context.Context, pageSize int, pageToken, filter, sort string) ([]Memo, error) {
	params := url.Values{}
	if pageSize > 0 {
		params.Set("pageSize", fmt.Sprintf("%d", pageSize))
	}
	if pageToken != "" {
		params.Set("pageToken", pageToken)
	}
	if filter != "" {
		params.Set("filter", filter)
	}
	if sort != "" {
		params.Set("sort", sort)
	}

	u := "/api/v1/memos"
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	resp, err := m.C.Get(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("list request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list failed: HTTP %d", resp.StatusCode)
	}

	var listResp struct {
		Memos []Memo `json:"memos"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return listResp.Memos, nil
}

// Update updates a memo.
func (m *MemoClient) Update(ctx context.Context, id, content, visibility string, tags []string) (*Memo, error) {
	body, err := json.Marshal(createMemoRequest{
		Content:    content,
		Visibility: visibility,
		Tags:       tags,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	u := "/api/v1/" + normalizeID(id)
	mask := []string{}
	if content != "" {
		mask = append(mask, "content")
	}
	if visibility != "" {
		mask = append(mask, "visibility")
	}
	if len(tags) > 0 {
		mask = append(mask, "tags")
	}
	if len(mask) > 0 {
		u += "?updateMask=" + strings.Join(mask, ",")
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", m.C.BaseURL+u, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.C.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("update request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update failed: HTTP %d", resp.StatusCode)
	}

	var memo Memo
	if err := json.NewDecoder(resp.Body).Decode(&memo); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &memo, nil
}

// Delete deletes a memo by ID.
func (m *MemoClient) Delete(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", m.C.BaseURL+"/api/v1/"+normalizeID(id), nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := m.C.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete failed: HTTP %d", resp.StatusCode)
	}
	return nil
}

// Search searches for memos containing the query.
func (m *MemoClient) Search(ctx context.Context, query string, pageSize int) ([]Memo, error) {
	filter := fmt.Sprintf(`content.contains("%s")`, query)
	return m.List(ctx, pageSize, "", filter, "")
}