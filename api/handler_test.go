package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	catpb "github.com/Includeoyasi/cat/pkg/cat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCapsuleLifecycle(t *testing.T) {
	handler := NewHandler(newFakeCatClient())
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
	if capsule.ID != 2 {
		t.Fatalf("unexpected capsule id: got %d, want 2", capsule.ID)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/capsule?id=2", nil)
	getRes := httptest.NewRecorder()

	router.ServeHTTP(getRes, getReq)

	if getRes.Code != http.StatusOK {
		t.Fatalf("unexpected get status: got %d, want %d", getRes.Code, http.StatusOK)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/delete?id=2", nil)
	deleteRes := httptest.NewRecorder()

	router.ServeHTTP(deleteRes, deleteReq)

	if deleteRes.Code != http.StatusNoContent {
		t.Fatalf("unexpected delete status: got %d, want %d", deleteRes.Code, http.StatusNoContent)
	}
}

func TestGetCapsuleReturnsNotFound(t *testing.T) {
	handler := NewHandler(newFakeCatClient())
	router := handler.InitRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/capsule?id=404", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("unexpected status: got %d, want %d", res.Code, http.StatusNotFound)
	}
}

func TestGetCapsuleReadsFromCatService(t *testing.T) {
	handler := NewHandler(newFakeCatClient())
	router := handler.InitRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/capsule?id=7", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d, want %d", res.Code, http.StatusOK)
	}

	var capsule Capsule
	if err := json.Unmarshal(res.Body.Bytes(), &capsule); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if capsule.Text != "Capsule from cat" {
		t.Fatalf("unexpected capsule text: got %q", capsule.Text)
	}
}

type fakeCatClient struct {
	nextID   int32
	capsules map[int32]*catpb.Capsule
	err      error
}

func newFakeCatClient() *fakeCatClient {
	return &fakeCatClient{
		nextID: 2,
		capsules: map[int32]*catpb.Capsule{
			7: {
				Id:        7,
				Author:    "Includeoyasi",
				Text:      "Capsule from cat",
				Label:     "grpc",
				CreatedAt: timestamppb.Now(),
			},
		},
	}
}

func (f fakeCatClient) GetCapsule(ctx context.Context, in *catpb.GetCapsuleRequest, opts ...grpc.CallOption) (*catpb.GetCapsuleResponce, error) {
	if f.err != nil {
		return nil, f.err
	}

	capsule, ok := f.capsules[in.GetId()]
	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &catpb.GetCapsuleResponce{Result: capsule}, nil
}

func (f *fakeCatClient) CreateCapsule(ctx context.Context, in *catpb.CreateCapsuleRequest, opts ...grpc.CallOption) (*catpb.CreateCapsuleResponce, error) {
	if f.err != nil {
		return nil, f.err
	}

	capsule := &catpb.Capsule{
		Id:        f.nextID,
		Author:    in.GetAuthor(),
		Text:      in.GetText(),
		Label:     in.GetLabel(),
		CreatedAt: timestamppb.Now(),
	}
	f.capsules[capsule.GetId()] = capsule
	f.nextID++

	return &catpb.CreateCapsuleResponce{Result: capsule}, nil
}

func (f *fakeCatClient) DeleteCapsule(ctx context.Context, in *catpb.DeleteCapsuleRequest, opts ...grpc.CallOption) (*catpb.DeleteCapsuleResponce, error) {
	if f.err != nil {
		return nil, f.err
	}

	if _, ok := f.capsules[in.GetId()]; !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}

	delete(f.capsules, in.GetId())

	return &catpb.DeleteCapsuleResponce{Deleted: true}, nil
}
