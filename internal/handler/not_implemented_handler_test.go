package handler

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotImplementedHandler(t *testing.T) {
	t.Run("Should return 'not implemented'", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/update/test", nil)

		w := httptest.NewRecorder()
		h := http.HandlerFunc(NotImplementedHandler)

		h.ServeHTTP(w, request)
		res := w.Result()
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusNotImplemented)
	})
}
