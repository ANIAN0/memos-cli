package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ANIAN0/memos-cli/pkg/httpclient"
)

func setupMockComments(handler http.HandlerFunc) (*httptest.Server, *httpclient.Client) {
	srv := httptest.NewServer(handler)
	c := httpclient.New(srv.URL)
	return srv, c
}

func TestCommentClient_List_Success(t *testing.T) {
	srv, c := setupMockComments(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/memos/123/comments" && r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"comments":[{"name":"memos/123/comments/1","content":"great!"},{"name":"memos/123/comments/2","content":"thanks"}]}`))
			return
		}
		w.WriteHeader(404)
	})
	defer srv.Close()

	cc := &CommentClient{C: c}
	comments, err := cc.List(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(comments))
	}
	if comments[0].Content != "great!" {
		t.Errorf("comments[0].Content = %q, want %q", comments[0].Content, "great!")
	}
}

func TestCommentClient_List_WithPrefix(t *testing.T) {
	srv, c := setupMockComments(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/memos/456/comments" && r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"comments":[]}`))
			return
		}
		w.WriteHeader(404)
	})
	defer srv.Close()

	cc := &CommentClient{C: c}
	comments, err := cc.List(context.Background(), "memos/456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 0 {
		t.Errorf("expected 0 comments, got %d", len(comments))
	}
}

func TestCommentClient_List_404(t *testing.T) {
	srv, c := setupMockComments(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	defer srv.Close()

	cc := &CommentClient{C: c}
	_, err := cc.List(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error on 404")
	}
}

func TestCommentClient_Create_Success(t *testing.T) {
	srv, c := setupMockComments(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/memos/123/comments" && r.Method == "POST" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"name":"memos/123/comments/3","content":"nice work!"}`))
			return
		}
		w.WriteHeader(404)
	})
	defer srv.Close()

	cc := &CommentClient{C: c}
	comment, err := cc.Create(context.Background(), "123", "nice work!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.Content != "nice work!" {
		t.Errorf("comment.Content = %q, want %q", comment.Content, "nice work!")
	}
}

func TestCommentClient_Create_EmptyContent(t *testing.T) {
	srv, c := setupMockComments(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{}`))
	})
	defer srv.Close()

	cc := &CommentClient{C: c}
	comment, err := cc.Create(context.Background(), "123", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment == nil {
		t.Error("expected non-nil comment")
	}
}