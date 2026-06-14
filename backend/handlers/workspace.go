package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
)

type WorkspaceReq struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

func ListWorkspaces(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	showArchived := c.Query("archived") == "true"

	var workspaces []models.Workspace
	query := database.DB.Where("user_id = ?", userID)
	if !showArchived {
		query = query.Where("is_archived = ?", false)
	}

	if err := query.Order("sort_order ASC, id ASC").Find(&workspaces).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถโหลดข้อมูลพื้นที่งานได้",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  workspaces,
		"total": len(workspaces),
	})
}

func CreateWorkspace(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var req WorkspaceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ข้อมูลไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	// Generate a slug based on Name + UserID to ensure uniqueIndex:idx_user_slug matches
	slugBase := strings.ToLower(strings.TrimSpace(req.Name))
	slugBase = strings.ReplaceAll(slugBase, " ", "-")
	slug := fmt.Sprintf("%s-%d", slugBase, userID)

	// Determine next SortOrder
	var maxSort int
	database.DB.Model(&models.Workspace{}).Where("user_id = ?", userID).Select("COALESCE(MAX(sort_order), 0)").Row().Scan(&maxSort)

	ws := models.Workspace{
		Name:      req.Name,
		Slug:      slug,
		Color:     req.Color,
		Icon:      req.Icon,
		UserID:    userID,
		SortOrder: maxSort + 1,
	}

	if ws.Color == "" {
		ws.Color = "#171717"
	}
	if ws.Icon == "" {
		ws.Icon = "briefcase"
	}

	if err := database.DB.Create(&ws).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถสร้างพื้นที่งานได้: " + err.Error(),
			"code":  "DATABASE_ERROR",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": ws,
	})
}

func UpdateWorkspace(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "รหัสพื้นที่งานไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	var ws models.Workspace
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&ws).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ไม่พบพื้นที่งานที่ต้องการแก้ไข",
			"code":  "WORKSPACE_NOT_FOUND",
		})
		return
	}

	type WorkspaceUpdateReq struct {
		Name       *string `json:"name"`
		Color      *string `json:"color"`
		Icon       *string `json:"icon"`
		IsArchived *bool   `json:"is_archived"`
		SortOrder  *int    `json:"sort_order"`
	}

	var req WorkspaceUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ข้อมูลไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	if req.Name != nil {
		ws.Name = *req.Name
		// Re-generate slug
		slugBase := strings.ToLower(strings.TrimSpace(ws.Name))
		slugBase = strings.ReplaceAll(slugBase, " ", "-")
		ws.Slug = fmt.Sprintf("%s-%d", slugBase, userID)
	}
	if req.Color != nil {
		ws.Color = *req.Color
	}
	if req.Icon != nil {
		ws.Icon = *req.Icon
	}
	if req.IsArchived != nil {
		ws.IsArchived = *req.IsArchived
	}
	if req.SortOrder != nil {
		ws.SortOrder = *req.SortOrder
	}

	if err := database.DB.Save(&ws).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถบันทึกการแก้ไขพื้นที่งานได้: " + err.Error(),
			"code":  "DATABASE_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ws,
	})
}

func DeleteWorkspace(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "รหัสพื้นที่งานไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	var ws models.Workspace
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&ws).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ไม่พบพื้นที่งานที่ต้องการลบ",
			"code":  "WORKSPACE_NOT_FOUND",
		})
		return
	}

	// Determine if we soft-delete (archive) or hard-delete
	hardDelete := c.Query("hard") == "true"

	if !hardDelete {
		// Soft delete: set IsArchived = true
		ws.IsArchived = true
		if err := database.DB.Save(&ws).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ไม่สามารถปิดใช้งานพื้นที่งานได้",
				"code":  "DATABASE_ERROR",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"message": "ปิดใช้งานพื้นที่งานเรียบร้อยแล้ว",
				"id":      ws.ID,
			},
		})
	} else {
		// Hard delete: delete workspace and all nested tasks Cascade
		if err := database.DB.Delete(&ws).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ไม่สามารถลบพื้นที่งานได้: " + err.Error(),
				"code":  "DATABASE_ERROR",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"message": "ลบพื้นที่งานเรียบร้อยแล้ว",
				"id":      id,
			},
		})
	}
}

func ReorderWorkspaces(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	type ReorderReq struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	var req ReorderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ข้อมูลไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	tx := database.DB.Begin()
	for idx, id := range req.IDs {
		if err := tx.Model(&models.Workspace{}).Where("id = ? AND user_id = ?", id, userID).Update("sort_order", idx).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "การเรียงลำดับล้มเหลว",
				"code":  "DATABASE_ERROR",
			})
			return
		}
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"message": "จัดเรียงลำดับสำเร็จแล้ว",
		},
	})
}
