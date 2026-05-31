package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCapsuleLifecycle(t *testing.T) {
	handler := NewHandler()
	router := handler.InitRoutes()

	createBody := []byte(`{"author":"Includeoyasi","text":"First memory capsule","label":"demo"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/api/create", bytes.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRes := httptest.NewRecorder()

	router.ServeHTTP(createRes, createReq)

	if createRes.Code != http.StatusCreated {
		t.Fatalf("unexpected create status: got %d, want %d", createRes.Code, http.StatusCreated)
	}

	var capsule Capsule
	if err := json.Unmarshal(createRes.Body.Bytes(), &capsule); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if capsule.ID != 1 {
		t.Fatalf("unexpected capsule id: got %d, want 1", capsule.ID)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/capsule?id=1", nil)
	getRes := httptest.NewRecorder()

	router.ServeHTTP(getRes, getReq)

	if getRes.Code != http.StatusOK {
		t.Fatalf("unexpected get status: got %d, want %d", getRes.Code, http.StatusOK)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/delete?id=1", nil)
	deleteRes := httptest.NewRecorder()

	router.ServeHTTP(deleteRes, deleteReq)

	if deleteRes.Code != http.StatusNoContent {
		t.Fatalf("unexpected delete status: got %d, want %d", deleteRes.Code, http.StatusNoContent)
	}
}

func TestGetCapsuleReturnsNotFound(t *testing.T) {
	handler := NewHandler()
	router := handler.InitRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/capsule?id=404", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("unexpected status: got %d, want %d", res.Code, http.StatusNotFound)
	}
}
