package api

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	catpb "github.com/Includeoyasi/cat/pkg/cat"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	catClient catpb.CatClient
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

var ErrCapsuleNotFound = errors.New("capsule not found")

func NewHandler(catClient catpb.CatClient) *Handler {
	return &Handler{
		catClient: catClient,
	}
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

	if h.catClient == nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "cat service is unavailable"})
		return
	}

	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 3*time.Second)
	defer cancel()

	res, err := h.catClient.CreateCapsule(reqCtx, &catpb.CreateCapsuleRequest{
		Author: req.Author,
		Text:   req.Text,
		Label:  req.Label,
	})
	if status.Code(err) == codes.InvalidArgument {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "cat service is unavailable"})
		return
	}

	ctx.JSON(http.StatusCreated, capsuleFromProto(res.GetResult()))
}

func (h Handler) delete(ctx *gin.Context) {
	id, ok := capsuleID(ctx)
	if !ok {
		return
	}

	if h.catClient == nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "cat service is unavailable"})
		return
	}

	if id > math.MaxInt32 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is out of range"})
		return
	}

	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 3*time.Second)
	defer cancel()

	_, err := h.catClient.DeleteCapsule(reqCtx, &catpb.DeleteCapsuleRequest{Id: int32(id)})
	if status.Code(err) == codes.NotFound {
		ctx.JSON(http.StatusNotFound, gin.H{"error": ErrCapsuleNotFound.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "cat service is unavailable"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h Handler) capsule(ctx *gin.Context) {
	id, ok := capsuleID(ctx)
	if !ok {
		return
	}

	if h.catClient == nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "cat service is unavailable"})
		return
	}

	h.getCapsuleFromCat(ctx, id)
}

func (h Handler) getCapsuleFromCat(ctx *gin.Context, id int64) {
	if id > math.MaxInt32 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is out of range"})
		return
	}

	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 3*time.Second)
	defer cancel()

	res, err := h.catClient.GetCapsule(reqCtx, &catpb.GetCapsuleRequest{Id: int32(id)})
	if status.Code(err) == codes.NotFound {
		ctx.JSON(http.StatusNotFound, gin.H{"error": ErrCapsuleNotFound.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "cat service is unavailable"})
		return
	}

	ctx.JSON(http.StatusOK, capsuleFromProto(res.GetResult()))
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

func capsuleFromProto(c *catpb.Capsule) Capsule {
	if c == nil {
		return Capsule{}
	}

	var createdAt time.Time
	if c.GetCreatedAt() != nil {
		createdAt = c.GetCreatedAt().AsTime()
	}

	return Capsule{
		ID:        int64(c.GetId()),
		Author:    c.GetAuthor(),
		Text:      c.GetText(),
		Label:     c.GetLabel(),
		CreatedAt: createdAt,
	}
}
