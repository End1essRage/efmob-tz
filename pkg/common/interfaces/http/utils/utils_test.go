package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	WriteJSON(w, http.StatusOK, map[string]string{"foo": "bar"})

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}

func TestParseQuery_StringAndPtr(t *testing.T) {
	type Query struct {
		Name    string  `schema:"name"`
		OptName *string `schema:"opt_name"`
	}

	r := httptest.NewRequest(http.MethodGet, "/?name=alice&opt_name=bob", nil)
	w := httptest.NewRecorder()

	var q Query
	ok := ParseQuery(w, r, &q)
	require.True(t, ok)
	require.Equal(t, "alice", q.Name)
	require.NotNil(t, q.OptName)
	require.Equal(t, "bob", *q.OptName)
}

func TestParseQuery_IntAndBool(t *testing.T) {
	type Query struct {
		Age  int  `schema:"age"`
		Flag bool `schema:"flag"`
	}

	r := httptest.NewRequest(http.MethodGet, "/?age=42&flag=true", nil)
	w := httptest.NewRecorder()

	var q Query
	ok := ParseQuery(w, r, &q)
	require.True(t, ok)
	require.Equal(t, 42, q.Age)
	require.True(t, q.Flag)
}

func TestParseQuery_BadInt(t *testing.T) {
	type Query struct {
		Age int `schema:"age"`
	}

	r := httptest.NewRequest(http.MethodGet, "/?age=notanint", nil)
	w := httptest.NewRecorder()

	var q Query
	ok := ParseQuery(w, r, &q)
	require.False(t, ok)
	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestParseQuery_NilDst(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	ok := ParseQuery(w, r, nil)
	require.False(t, ok)
	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func TestParseQuery_InvalidDst(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	x := 42
	ok := ParseQuery(w, r, &x)
	require.False(t, ok)
	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}
