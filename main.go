package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Mã giảm giá mẫu
type Discount struct {
	Code     string    `json:"code"`
	Type     string    `json:"type"` // "percent" hoặc "cash"
	Value    float64   `json:"value"`
	MinOrder float64   `json:"min_order"`
	ExpireAt time.Time `json:"expire_at"`
}

// Danh sách mã giả lập
var discountStore = map[string]Discount{}

// Tạo mã mới
func createDiscount(c *gin.Context) {
	var d Discount
	if err := c.BindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	// Kiểm tra mã đã tồn tại
	if _, exists := discountStore[d.Code]; exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã đã tồn tại"})
		return
	}

	// Kiểm tra dữ liệu hợp lệ
	if d.Type != "percent" && d.Type != "cash" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Loại mã không hợp lệ"})
		return
	}
	if d.Value <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Giá trị giảm giá phải lớn hơn 0"})
		return
	}
	if d.MinOrder < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Đơn tối thiểu không hợp lệ"})
		return
	}
	if d.ExpireAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ngày hết hạn không hợp lệ"})
		return
	}

	discountStore[d.Code] = d
	c.JSON(http.StatusOK, gin.H{"message": "Tạo mã thành công", "data": d})
}

// Lấy danh sách mã
func listDiscounts(c *gin.Context) {
	list := []Discount{}
	for _, v := range discountStore {
		list = append(list, v)
	}
	c.JSON(http.StatusOK, list)
}

// Áp dụng mã vào đơn hàng
func applyDiscount(c *gin.Context) {
	type Request struct {
		Code       string  `json:"code"`
		OrderTotal float64 `json:"order_total"`
	}

	var req Request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	discount, ok := discountStore[req.Code]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mã không tồn tại"})
		return
	}

	if time.Now().After(discount.ExpireAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mã đã hết hạn"})
		return
	}

	if req.OrderTotal < discount.MinOrder {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không đủ điều kiện đơn tối thiểu"})
		return
	}

	var discountAmount float64
	if discount.Type == "percent" {
		discountAmount = req.OrderTotal * discount.Value / 100
	} else {
		discountAmount = discount.Value
	}

	c.JSON(http.StatusOK, gin.H{
		"original_total": req.OrderTotal,
		"discount":       discountAmount,
		"final_total":    req.OrderTotal - discountAmount,
	})
}

func main() {
	r := gin.Default()

	r.POST("/discounts", createDiscount)
	r.GET("/discounts", listDiscounts)
	r.POST("/apply-discount", applyDiscount)

	r.Run("localhost:8080")
}
