package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ANIAN0/memos-cli/pkg/httpclient"
)

func setupMockAttachments(handler http.HandlerFunc) (*httptest.Server, *httpclient.Client) {
	srv := httptest.NewServer(handler)
	c := httpclient.New(srv.URL)
	return srv, c
}

func TestAttachmentClient_Upload_Success(t *testing.T) {
	srv, c := setupMockAttachments(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/attachments" && r.Method == "POST" {
			// Verify it's multipart
			contentType := r.Header.Get("Content-Type")
			if contentType == "" || !contains(contentType, "multipart/form-data") {
				t.Errorf("expected multipart/form-data, got %q", contentType)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"name":"attachments/123","filename":"test.txt","size":11,"type":"text/plain"}`))
			return
		}
		w.WriteHeader(404)
	})
	defer srv.Close()

	// Create temp file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("hello world"), 0644)

	ac := &AttachmentClient{C: c}
	att, err := ac.Upload(context.Background(), testFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if att.Name != "attachments/123" {
		t.Errorf("name = %q, want %q", att.Name, "attachments/123")
	}
}

func TestAttachmentClient_List_Success(t *testing.T) {
	srv, c := setupMockAttachments(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/attachments" && r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[{"name":"attachments/1","filename":"a.txt","size":100},{"name":"attachments/2","filename":"b.png","size":200}]`))
			return
		}
		w.WriteHeader(404)
	})
	defer srv.Close()

	ac := &AttachmentClient{C: c}
	atts, err := ac.List(context.Background(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(atts) != 2 {
		t.Errorf("expected 2 attachments, got %d", len(atts))
	}
}

func TestAttachmentClient_List_WithPageSize(t *testing.T) {
	var receivedURL string
	srv, c := setupMockAttachments(func(w http.ResponseWriter, r *http.Request) {
		receivedURL = r.URL.String()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[]`))
	})
	defer srv.Close()

	ac := &AttachmentClient{C: c}
	_, err := ac.List(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(receivedURL, "pageSize=10") {
		t.Errorf("URL should contain pageSize=10, got: %s", receivedURL)
	}
}

func TestAttachmentClient_Get_Success(t *testing.T) {
	fileContent := []byte("attachment content")
	srv, c := setupMockAttachments(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/attachments/123" && r.Method == "GET" {
			w.Write(fileContent)
			return
		}
		w.WriteHeader(404)
	})
	defer srv.Close()

	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "downloaded.txt")

	ac := &AttachmentClient{C: c}
	err := ac.Get(context.Background(), "123", outputFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	savedContent, _ := os.ReadFile(outputFile)
	if string(savedContent) != string(fileContent) {
		t.Errorf("content mismatch: got %q, want %q", string(savedContent), string(fileContent))
	}
}

func TestAttachmentClient_Delete_Success(t *testing.T) {
	srv, c := setupMockAttachments(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/attachments/123" && r.Method == "DELETE" {
			w.WriteHeader(200)
			return
		}
		w.WriteHeader(404)
	})
	defer srv.Close()

	ac := &AttachmentClient{C: c}
	err := ac.Delete(context.Background(), "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAttachmentClient_Delete_WithPrefix(t *testing.T) {
	srv, c := setupMockAttachments(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/attachments/456" && r.Method == "DELETE" {
			w.WriteHeader(200)
			return
		}
		w.WriteHeader(404)
	})
	defer srv.Close()

	ac := &AttachmentClient{C: c}
	err := ac.Delete(context.Background(), "attachments/456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}



// Helper to read response body for testing
func readBody(t *testing.T, r *http.Request) []byte {
	t.Helper()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	return body
}