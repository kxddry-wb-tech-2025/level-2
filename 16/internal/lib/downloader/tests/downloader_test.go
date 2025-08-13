package downloader_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"wget/internal/lib/downloader"
)

func TestDownloadOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	d := downloader.New(0, "ua")
	u, _ := url.Parse(srv.URL)
	body, ct, err := d.Download(context.Background(), u)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = body.Close() }()
	if ct != "text/html" {
		t.Errorf("got %q", ct)
	}
}
