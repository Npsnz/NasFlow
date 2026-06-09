package handlers

import (
	"net/http"
	"strconv"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
)

type TagReq struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
}

func ListTags(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var tags []models.Tag
	if err := database.DB.Where("user_id = ?", userID).Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถโหลดข้อมูลแท็กได้",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tags,
		"total": len(tags),
	})
}

func CreateTag(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var req TagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ข้อมูลแท็กไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	tag := models.Tag{
		Name:   req.Name,
		Color:  req.Color,
		UserID: userID,
	}

	if err := database.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถสร้างแท็กได้",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": tag,
	})
}

func UpdateTag(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "รหัสแท็กไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	var tag models.Tag
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&tag).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ไม่พบแท็กที่ต้องการแก้ไข",
			"code":  "TAG_NOT_FOUND",
		})
		return
	}

	var req TagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ข้อมูลไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	tag.Name = req.Name
	tag.Color = req.Color

	if err := database.DB.Save(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถบันทึกแท็กได้",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tag,
	})
}

func DeleteTag(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "รหัสแท็กไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	var tag models.Tag
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&tag).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ไม่พบแท็กที่ต้องการลบ",
			"code":  "TAG_NOT_FOUND",
		})
		return
	}

	if err := database.DB.Delete(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถลบแท็กได้",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"message": "ลบแท็กสำเร็จแล้ว",
			"id":      id,
		},
	})
}
