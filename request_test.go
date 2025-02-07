package httpclientutils_test

import (
	_ "crypto/tls"
	"encoding/json"
	"github.com/InheritxSolution/httpclientutils"
	"net/http"
	"net/http/httptest"
	"testing"
	_ "time"

	"github.com/stretchr/testify/assert"
)

func TestMakeHTTPRequest_GetRequest(t *testing.T) {
	mockResponse := map[string]string{"message": "success"}
	mockResponseBody, _ := json.Marshal(mockResponse)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBody)
	}))
	defer ts.Close()

	var result map[string]string
	status, headers, body, err := httpclientutils.MakeHTTPRequest(
		httpclientutils.WithMethod(http.MethodGet),
		httpclientutils.WithURL(ts.URL),
		httpclientutils.WithResolveResponse(&result),
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, true, len(body) > 0)
	assert.Equal(t, "application/json", headers.Get("Content-Type"))
	assert.Equal(t, mockResponse, result)
}

func TestMakeHTTPRequest_PostRequest(t *testing.T) {
	mockResponse := map[string]string{"message": "created"}
	mockResponseBody, _ := json.Marshal(mockResponse)
	mockRequestBody := map[string]string{"name": "test"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		var reqBody map[string]string
		json.NewDecoder(r.Body).Decode(&reqBody)
		assert.Equal(t, "test", reqBody["name"])
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(mockResponseBody)
	}))
	defer ts.Close()

	var result map[string]string
	status, _, _, err := httpclientutils.MakeHTTPRequest(
		httpclientutils.WithMethod(http.MethodPost),
		httpclientutils.WithURL(ts.URL),
		httpclientutils.WithBody(mockRequestBody),
		httpclientutils.WithResolveResponse(&result),
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)
	assert.Equal(t, mockResponse, result)
}

func TestMakeHTTPRequest_PutRequest(t *testing.T) {
	mockResponse := map[string]string{"message": "updated"}
	mockResponseBody, _ := json.Marshal(mockResponse)
	mockRequestBody := map[string]string{"name": "test"}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		var reqBody map[string]string
		json.NewDecoder(r.Body).Decode(&reqBody)
		assert.Equal(t, "test", reqBody["name"])
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBody)
	}))
	defer ts.Close()

	var result map[string]string
	status, _, _, err := httpclientutils.MakeHTTPRequest(
		httpclientutils.WithMethod(http.MethodPut),
		httpclientutils.WithURL(ts.URL),
		httpclientutils.WithBody(mockRequestBody),
		httpclientutils.WithResolveResponse(&result),
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, mockResponse, result)
}

func TestMakeHTTPRequest_DeleteRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	status, _, _, err := httpclientutils.MakeHTTPRequest(
		httpclientutils.WithMethod(http.MethodDelete),
		httpclientutils.WithURL(ts.URL),
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, status)
}

func TestMakeHTTPRequest_ResolveJSONResponse(t *testing.T) {
	mockResponse := map[string]string{"message": "resolved"}
	mockResponseBody, _ := json.Marshal(mockResponse)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponseBody)
	}))
	defer ts.Close()

	var result map[string]string
	status, _, _, err := httpclientutils.MakeHTTPRequest(
		httpclientutils.WithMethod(http.MethodGet),
		httpclientutils.WithURL(ts.URL),
		httpclientutils.WithResolveResponse(&result),
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, mockResponse, result)
}
