package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ANIAN0/memos-cli/pkg/httpclient"
)

// Comment represents a Memos comment.
type Comment struct {
	Name       string `json:"name"`
	UID        string `json:"uid"`
	Content    string `json:"content"`
	Creator    string `json:"creator"`
	CreateTime string `json:"createTime"`
}

// CommentClient handles comment operations.
type CommentClient struct {
	C *httpclient.Client
}

// createCommentRequest is the request body for creating a comment.
type createCommentRequest struct {
	Content string `json:"content"`
}

// List returns comments for a memo.
func (c *CommentClient) List(ctx context.Context, memoID string) ([]Comment, error) {
	memoID = normalizeID(memoID)
	resp, err := c.C.Get(ctx, "/api/v1/"+memoID+"/comments")
	if err != nil {
		return nil, fmt.Errorf("list request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("list failed: HTTP %d", resp.StatusCode)
	}

	var respBody struct {
		Comments []Comment `json:"comments"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return respBody.Comments, nil
}

// Create creates a new comment on a memo.
func (c *CommentClient) Create(ctx context.Context, memoID, content string) (*Comment, error) {
	memoID = normalizeID(memoID)
	body, err := json.Marshal(createCommentRequest{Content: content})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.C.Post(ctx, "/api/v1/"+memoID+"/comments", "application/json", strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("create failed: HTTP %d", resp.StatusCode)
	}

	var comment Comment
	if err := json.NewDecoder(resp.Body).Decode(&comment); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &comment, nil
}