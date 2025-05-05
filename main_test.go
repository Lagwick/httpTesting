package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList["moscow"])},
	}

	for _, tc := range requests {
		url := "/cafe?city=moscow&count=" + strconv.Itoa(tc.count)

		req := httptest.NewRequest("GET", url, nil)
		resp := httptest.NewRecorder()

		handler.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code)

		result := strings.TrimSpace(resp.Body.String())
		if result == "" {
			assert.Equal(t, 0, tc.want)
			continue
		}

		cafes := strings.Split(result, ",")
		assert.Equal(t, tc.want, len(cafes))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, tc := range requests {
		req := httptest.NewRequest("GET", "/cafe?city=moscow&search="+tc.search, nil)
		resp := httptest.NewRecorder()

		handler.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code)

		body := strings.TrimSpace(resp.Body.String())
		if body == "" {
			assert.Equal(t, 0, tc.wantCount)
			continue
		}

		cafes := strings.Split(body, ",")
		assert.Equal(t, tc.wantCount, len(cafes))

		searchLower := strings.ToLower(tc.search)
		for _, cafe := range cafes {
			assert.Contains(t, strings.ToLower(cafe), searchLower)
		}
	}
}
