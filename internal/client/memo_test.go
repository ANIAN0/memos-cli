package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ANIAN0/memos-cli/pkg/httpclient"
)

func setupMockMemos(handler http.HandlerFunc) (*httptest.Server, *httpclient.Client) {
	srv := httptest.NewServer(handler)
	c := httpclient.New(srv.URL)
	return srv, c
}

func TestMemoClient_Create_Success(t *testing.T) {
	srv, c := setupMockMemos(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/memos" && r.Method == "POST" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"name":"memos/123","uid":"abc","content":"test","visibility":"PRIVATE","tags":["work"]}`))
			return
		}
		w.WriteHeader(404)
	})
	defer srv.Close()

	mc := &MemoClient{C: c}
	memo, err := mc.Create(context.Background(), "test", "PRIVATE", []string{"work"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if memo.Name != "memos/123" {
		t.Errorf("name = %q, want %q", memo.Name, "memos/123")
	}
}

func TestMemoClient_Get_Success(t *testing.T) {
	srv, c := setupMockMemos(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/memos/123" && r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"name":"memos/123","content":"hello"}`))
			return
		}
		w.WriteHeader(404)
	})
	defer srv.Close()

	mc := &MemoClient{C: c}
	memo, err := mc.Get(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if memo.Content != "hello" {
		t.Errorf("content = %q, want %q", memo.Content, "hello")
	}
}

func TestMemoClient_Get_WithPrefix(t *testing.T) {
	srv, c := setupMockMemos(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/memos/456" && r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"name":"memos/456","content":"world"}`))
			return
		}
		w.WriteHeader(404)
	})
	defer srv.Close()

	mc := &MemoClient{C: c}
	memo, err := mc.Get(context.Background(), "memos/456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if memo.Content != "world" {
		t.Errorf("content = %q, want %q", memo.Content, "world")
	}
}

func TestMemoClient_List_Success(t *testing.T) {
	srv, c := setupMockMemos(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/memos" && r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"memos":[{"name":"memos/1","content":"a"},{"name":"memos/2","content":"b"}]}`))
			return
		}
		w.WriteHeader(404)
	})
	defer srv.Close()

	mc := &MemoClient{C: c}
	memos, err := mc.List(context.Background(), 0, "", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(memos) != 2 {
		t.Errorf("expected 2 memos, got %d", len(memos))
	}
}

func TestMemoClient_List_WithFilter(t *testing.T) {
	var receivedURL string
	srv, c := setupMockMemos(func(w http.ResponseWriter, r *http.Request) {
		receivedURL = r.URL.String()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"memos":[]}`))
	})
	defer srv.Close()

	mc := &MemoClient{C: c}
	_, err := mc.List(context.Background(), 10, "", "tag=='work'", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(receivedURL, "filter=tag%3D%3D%27work%27") {
		t.Errorf("URL should contain filter, got: %s", receivedURL)
	}
}

func TestMemoClient_Update_Success(t *testing.T) {
	srv, c := setupMockMemos(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/memos/123" && r.Method == "PATCH" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"name":"memos/123","content":"updated"}`))
			return
		}
		w.WriteHeader(404)
	})
	defer srv.Close()

	mc := &MemoClient{C: c}
	memo, err := mc.Update(context.Background(), "123", "updated", "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if memo.Content != "updated" {
		t.Errorf("content = %q, want %q", memo.Content, "updated")
	}
}

func TestMemoClient_Delete_Success(t *testing.T) {
	srv, c := setupMockMemos(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/memos/123" && r.Method == "DELETE" {
			w.WriteHeader(200)
			return
		}
		w.WriteHeader(404)
	})
	defer srv.Close()

	mc := &MemoClient{C: c}
	err := mc.Delete(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMemoClient_Search_Success(t *testing.T) {
	srv, c := setupMockMemos(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"memos":[{"name":"memos/1","content":"hello world"}]}`))
	})
	defer srv.Close()

	mc := &MemoClient{C: c}
	memos, err := mc.Search(context.Background(), "hello", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(memos) != 1 {
		t.Errorf("expected 1 memo, got %d", len(memos))
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}