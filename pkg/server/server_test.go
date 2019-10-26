package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zerok/servethis/pkg/server"
)

// When the user request an index.html file, it should be served if it
// exists and the user SHOULD NOT be redirected to /.
func TestNoIndexRedirect(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/index.html", nil)
	require.NoError(t, err)
	ctx := context.Background()
	rootDir := filepath.Join("..", "..", "testdata")
	rootDir, _ = filepath.Abs(rootDir)
	s := server.New(ctx, rootDir)
	s.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
}

// When requesting / then a folder listing should be returned and not
// the content of an index.html file.
func TestAutoIndex(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)
	ctx := context.Background()
	rootDir := filepath.Join("..", "..", "testdata")
	rootDir, _ = filepath.Abs(rootDir)
	s := server.New(ctx, rootDir)
	s.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	require.NotContains(t, body, "Index file")
	require.Contains(t, body, "index.html")
}
