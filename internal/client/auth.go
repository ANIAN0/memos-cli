// Package client provides Memos API clients.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ANIAN0/memos-cli/pkg/httpclient"
)

// AuthClient handles authentication operations.
type AuthClient struct {
	C *httpclient.Client
}

// User represents a Memos user.
type User struct {
	Name        string `json:"name"`        // "users/123"
	UID         string `json:"uid"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Role        string `json:"role"` // ADMIN, USER
	Avatar      string `json:"avatar"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// GetCurrentUserName returns the current username from the token.
// Since Memos doesn't have a direct /users/me endpoint that works with PAT,
// we extract the creator info from the list memos response.
func (a *AuthClient) GetCurrentUserName(ctx context.Context) (string, error) {
	// Try to list memos to get creator info
	resp, err := a.C.Get(ctx, "/api/v1/memos?pageSize=1")
	if err != nil {
		return "", fmt.Errorf("get current user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get current user failed: HTTP %d", resp.StatusCode)
	}

	var listResp struct {
		Memos []struct {
			Creator string `json:"creator"`
		} `json:"memos"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(listResp.Memos) == 0 {
		return "unknown (no memos found)", nil
	}

	// Creator format is "users/username"
	creator := listResp.Memos[0].Creator
	if len(creator) > 6 && creator[:6] == "users/" {
		return creator[6:], nil
	}
	return creator, nil
}
