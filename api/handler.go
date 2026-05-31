package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	store *CapsuleStore
}

type Capsule struct {
	ID        int64     `json:"id"`
	Author    string    `json:"author"`
	Text      string    `json:"text"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateCapsuleRequest struct {
	Author string `json:"author" binding:"required"`
	Text   string `json:"text" binding:"required"`
	Label  string `json:"label"`
}

type authRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CapsuleStore struct {
	mu       sync.RWMutex
	nextID   int64
	capsules map[int64]Capsule
}

var ErrCapsuleNotFound = errors.New("capsule not found")

func NewHandler() *Handler {
	return &Handler{
		store: NewCapsuleStore(),
	}
}

func NewCapsuleStore() *CapsuleStore {
	return &CapsuleStore{
		nextID:   1,
		capsules: make(map[int64]Capsule),
	}
}

func (s *CapsuleStore) Create(req CreateCapsuleRequest) Capsule {
	s.mu.Lock()
	defer s.mu.Unlock()

	capsule := Capsule{
		ID:        s.nextID,
		Author:    strings.TrimSpace(req.Author),
		Text:      strings.TrimSpace(req.Text),
		Label:     strings.TrimSpace(req.Label),
		CreatedAt: time.Now().UTC(),
	}

	s.capsules[capsule.ID] = capsule
	s.nextID++

	return capsule
}

func (s *CapsuleStore) Get(id int64) (Capsule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	capsule, ok := s.capsules[id]
	if !ok {
		return Capsule{}, ErrCapsuleNotFound
	}

	return capsule, nil
}

func (s *CapsuleStore) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.capsules[id]; !ok {
		return ErrCapsuleNotFound
	}

	delete(s.capsules, id)

	return nil
}

func (h Handler) singIn(ctx *gin.Context) {
	var req authRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": "demo-token",
	})
}

func (h Handler) singUp(ctx *gin.Context) {
	var req authRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"username": req.Username,
	})
}

func (h Handler) create(ctx *gin.Context) {
	var req CreateCapsuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	capsule := h.store.Create(req)
	ctx.JSON(http.StatusCreated, capsule)
}

func (h Handler) delete(ctx *gin.Context) {
	id, ok := capsuleID(ctx)
	if !ok {
		return
	}

	if err := h.store.Delete(id); errors.Is(err, ErrCapsuleNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": ErrCapsuleNotFound.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h Handler) capsule(ctx *gin.Context) {
	id, ok := capsuleID(ctx)
	if !ok {
		return
	}

	capsule, err := h.store.Get(id)
	if errors.Is(err, ErrCapsuleNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": ErrCapsuleNotFound.Error()})
		return
	}

	ctx.JSON(http.StatusOK, capsule)
}

func (h Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	auth := router.Group("/auth")
	{
		auth.POST("/sing-in", h.singIn)
		auth.POST("/sing-up", h.singUp)
	}

	api := router.Group("/api")
	{
		api.POST("/create", h.create)
		api.DELETE("/delete", h.delete)
		api.GET("/capsule", h.capsule)
	}

	return router
}

func capsuleID(ctx *gin.Context) (int64, bool) {
	rawID := ctx.Query("id")
	if rawID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id query parameter is required"})
		return 0, false
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return 0, false
	}

	return id, true
}
