//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	memos "github.com/ANIAN0/memos-cli/internal/client"
	"github.com/ANIAN0/memos-cli/pkg/httpclient"
)

// skipIfNoMemos skips the test if MEMOS_TEST_* env vars are not set.
func skipIfNoMemos(t *testing.T) (url, token string) {
	t.Helper()
	url = os.Getenv("MEMOS_TEST_URL")
	token = os.Getenv("MEMOS_TEST_TOKEN")
	if url == "" || token == "" {
		t.Skip("set MEMOS_TEST_URL, MEMOS_TEST_TOKEN to run integration tests")
	}
	return url, token
}

// newClient creates a new HTTP client for testing.
func newClient(t *testing.T, url, token string) *httpclient.Client {
	t.Helper()
	return httpclient.New(url,
		httpclient.WithTimeout(60*time.Second),
		httpclient.WithToken(token),
		httpclient.WithAuthHeader("Authorization"),
		httpclient.WithVerbose(true),
	)
}

// TestE2E_MemoCreateListSearchUpdateDelete tests the full memo workflow:
// create -> get -> list -> search -> update -> delete
func TestE2E_MemoCreateListSearchUpdateDelete(t *testing.T) {
	url, token := skipIfNoMemos(t)
	ctx := context.Background()

	c := newClient(t, url, token)
	mc := &memos.MemoClient{C: c}

	// Create unique content for search
	unique := fmt.Sprintf("integration-test-%d", time.Now().UnixNano())
	content := unique + " content for search"

	// Create memo
	memo, err := mc.Create(ctx, content, "PRIVATE", []string{"integration", "test"})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if memo.Name == "" {
		t.Fatal("create returned empty memo name")
	}
	t.Logf("created memo: %s", memo.Name)

	// Cleanup: delete memo after test
	t.Cleanup(func() {
		ctx := context.Background()
		_ = mc.Delete(ctx, memo.Name)
		t.Logf("cleanup: deleted memo %s", memo.Name)
	})

	// Get memo
	got, err := mc.Get(ctx, memo.Name)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.Content != content {
		t.Errorf("get content = %q, want %q", got.Content, content)
	}
	t.Logf("get memo: content matches")

	// List memos
	list, err := mc.List(ctx, 100, "", "", "")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	found := false
	for _, m := range list {
		if m.Name == memo.Name {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created memo not found in list (got %d memos)", len(list))
	} else {
		t.Logf("list: found created memo")
	}

	// Search via content.contains
	searchFilter := fmt.Sprintf(`content.contains("%s")`, unique)
	searchList, err := mc.List(ctx, 100, "", searchFilter, "")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	found = false
	for _, m := range searchList {
		if m.Name == memo.Name {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("memo not found via content.contains search (got %d results)", len(searchList))
	} else {
		t.Logf("search: found memo via content.contains")
	}

	// Update memo
	updatedContent := content + " (updated)"
	updated, err := mc.Update(ctx, memo.Name, updatedContent, "", nil)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if !strings.Contains(updated.Content, "(updated)") {
		t.Errorf("update did not change content: %q", updated.Content)
	}
	t.Logf("update: content updated successfully")

	// Delete memo
	if err := mc.Delete(ctx, memo.Name); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	t.Logf("delete: memo deleted successfully")
}

// TestE2E_CommentCreateList tests the comment workflow:
// create memo -> create comment -> list comments -> cleanup
func TestE2E_CommentCreateList(t *testing.T) {
	url, token := skipIfNoMemos(t)
	ctx := context.Background()

	c := newClient(t, url, token)
	mc := &memos.MemoClient{C: c}
	cc := &memos.CommentClient{C: c}

	// Create memo for comments
	memo, err := mc.Create(ctx, "Memo for comments", "PRIVATE", nil)
	if err != nil {
		t.Fatalf("create memo failed: %v", err)
	}
	t.Cleanup(func() { mc.Delete(ctx, memo.Name) })

	// Create comment
	commentContent := fmt.Sprintf("Test comment %d", time.Now().UnixNano())
	comment, err := cc.Create(ctx, memo.Name, commentContent)
	if err != nil {
		t.Fatalf("create comment failed: %v", err)
	}
	if comment.Name == "" {
		t.Fatal("create comment returned empty name")
	}
	t.Logf("created comment: %s", comment.Name)

	// List comments
	comments, err := cc.List(ctx, memo.Name)
	if err != nil {
		t.Fatalf("list comments failed: %v", err)
	}
	found := false
	for _, c := range comments {
		if c.Name == comment.Name {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created comment not found in list (got %d comments)", len(comments))
	} else {
		t.Logf("list: found created comment")
	}
}

// TestE2E_AttachmentUploadGet tests the attachment workflow:
// upload -> get -> verify content -> delete
func TestE2E_AttachmentUploadGet(t *testing.T) {
	url, token := skipIfNoMemos(t)
	ctx := context.Background()

	c := newClient(t, url, token)
	ac := &memos.AttachmentClient{C: c}

	// Create test file
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "test-attachment.txt")
	testContent := fmt.Sprintf("attachment integration test %d", time.Now().UnixNano())
	if err := os.WriteFile(src, []byte(testContent), 0644); err != nil {
		t.Fatalf("create test file failed: %v", err)
	}

	// Upload
	att, err := ac.Upload(ctx, src)
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}
	if att.Name == "" {
		t.Fatal("upload returned empty attachment name")
	}
	t.Logf("uploaded attachment: %s", att.Name)

	// Cleanup: delete attachment after test
	t.Cleanup(func() {
		ctx := context.Background()
		_ = ac.Delete(ctx, att.Name)
		t.Logf("cleanup: deleted attachment %s", att.Name)
	})

	// List attachments
	atts, err := ac.List(ctx, 100)
	if err != nil {
		t.Fatalf("list attachments failed: %v", err)
	}
	found := false
	for _, a := range atts {
		if a.Name == att.Name {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("uploaded attachment not found in list (got %d attachments)", len(atts))
	} else {
		t.Logf("list: found uploaded attachment")
	}

	// Get (download) attachment
	dst := filepath.Join(tmpDir, "downloaded.txt")
	if err := ac.Get(ctx, att.Name, dst); err != nil {
		t.Fatalf("get attachment failed: %v", err)
	}

	// Verify content
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read downloaded file failed: %v", err)
	}
	if string(got) != testContent {
		t.Errorf("downloaded content mismatch:\n  got:  %q\n  want: %q", string(got), testContent)
	}
	t.Logf("get: attachment content verified")
}