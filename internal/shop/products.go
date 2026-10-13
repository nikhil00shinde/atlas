package shop 


import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Product struct {
	ID         int64  `gorm:"primaryKey" json:"id"`
	Name       string `gorm:"not null"   json:"name"`
	PriceCents int64  `gorm:"not null"   json:"price_cents"`
}

type Store struct {
	DB	*gorm.DB
}

func (s *Store) ListProducts(ctx context.Context) ([]Product, error) {
	products := []Product{}
	err := s.DB.WithContext(ctx).Order("id").Find(&products).Error 
	return products, err 
}



func (s *Store) Ping(ctx context.Context) error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err 
	}
	return sqlDB.PingContext(ctx)
}



type Handler struct {
	Store *Store
	Log		*slog.Logger 
}

func (h *Handler) ListProducts(c *gin.Context) {
	products, err := h.Store.ListProducts(c.Request.Context())
	if err != nil {
		h.Log.Error("list products", "err", err, "request_id", c.GetString("request_id"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, products)
}


func (h *Handler) Live(c *gin.Context) {
	c.Status(http.StatusOK)
}


func (h *Handler) Ready(c *gin.Context) {
	if err := h.Store.Ping(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db un reachable"})
		return 
	}

	c.Status(http.StatusOK)
}



