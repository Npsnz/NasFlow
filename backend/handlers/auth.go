package handlers

import (
	"errors"
	"net/http"
	"time"

	"backend/database"
	"backend/middleware"
	"backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Helper to set cookie
func setTokenCookies(c *gin.Context, accessToken, refreshToken string, rememberMe bool) {
	// Set access token cookie (valid for 15 mins)
	c.SetCookie("access_token", accessToken, 15*60, "/", "", false, true)

	// Set refresh token cookie (valid for 30 days)
	refreshMaxAge := 30 * 24 * 3600
	if !rememberMe {
		refreshMaxAge = 0 // Session-only cookie
	}
	c.SetCookie("refresh_token", refreshToken, refreshMaxAge, "/api/auth", "", false, true)
}

// Helper to clear cookies
func clearTokenCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/api/auth", "", false, true)
}

func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ข้อมูลไม่ถูกต้อง: " + err.Error(),
			"code":  "BAD_REQUEST",
		})
		return
	}

	// Check if user already exists
	var existingUser models.User
	err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "อีเมลนี้ถูกใช้งานแล้ว",
			"code":  "EMAIL_ALREADY_EXISTS",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ระบบไม่สามารถเข้ารหัสรหัสผ่านได้",
			"code":  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// Save to DB
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถสร้างบัญชีผู้ใช้ได้: " + err.Error(),
			"code":  "DATABASE_ERROR",
		})
		return
	}

	// Seed workspaces for the new user
	if err := database.SeedDefaultWorkspaces(user.ID); err != nil {
		// Log error, but don't fail registration
		println("Failed to seed default workspaces: ", err.Error())
	}

	// Generate tokens
	accessToken, err := middleware.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถออกโทเค็นล็อกอินได้",
			"code":  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถออกรีเฟรชโทเค็นได้",
			"code":  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	setTokenCookies(c, accessToken, refreshToken, true)

	c.JSON(http.StatusCreated, gin.H{
		"data":  user,
		"token": accessToken,
	})
}

func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "อีเมลหรือรหัสผ่านไม่ถูกต้องตามเงื่อนไข",
			"code":  "BAD_REQUEST",
		})
		return
	}

	var user models.User
	err := database.DB.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "อีเมลหรือรหัสผ่านไม่ถูกต้อง",
				"code":  "INVALID_CREDENTIALS",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "เกิดข้อผิดพลาดในการตรวจสอบข้อมูล",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "อีเมลหรือรหัสผ่านไม่ถูกต้อง",
			"code":  "INVALID_CREDENTIALS",
		})
		return
	}

	// Generate tokens
	accessToken, err := middleware.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถออกโทเค็นล็อกอินได้",
			"code":  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถออกรีเฟรชโทเค็นได้",
			"code":  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	// Optional check for Remember Me
	rememberMe := true // default to true, can read from req body or querystring if needed

	setTokenCookies(c, accessToken, refreshToken, rememberMe)

	c.JSON(http.StatusOK, gin.H{
		"data":  user,
		"token": accessToken,
	})
}

func Logout(c *gin.Context) {
	clearTokenCookies(c)
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"message": "ออกจากระบบสำเร็จ",
		},
	})
}

func GetMe(c *gin.Context) {
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ไม่ได้เข้าสู่ระบบ",
			"code":  "UNAUTHORIZED",
		})
		return
	}
	user := userVal.(models.User)
	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func RefreshToken(c *gin.Context) {
	// Read refresh token from cookie
	refreshCookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ไม่พบรีเฟรชโทเค็น",
			"code":  "REFRESH_TOKEN_NOT_FOUND",
		})
		return
	}

	claims, err := middleware.ValidateToken(refreshCookie)
	if err != nil || claims.TokenType != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "รีเฟรชโทเค็นหมดอายุหรือทำงานไม่ถูกต้อง",
			"code":  "INVALID_REFRESH_TOKEN",
		})
		return
	}

	// Fetch user
	var user models.User
	if err := database.DB.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ไม่พบบัญชีผู้ใช้",
			"code":  "USER_NOT_FOUND",
		})
		return
	}

	// Generate new tokens
	accessToken, err := middleware.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถออกโทเค็นใหม่ได้",
			"code":  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถออกรีเฟรชโทเค็นใหม่ได้",
			"code":  "INTERNAL_SERVER_ERROR",
		})
		return
	}

	// Keep cookies updated
	setTokenCookies(c, accessToken, refreshToken, true)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"access_token": accessToken,
			"token":        accessToken,
			"expires_at":   time.Now().Add(15 * time.Minute).Unix(),
		},
	})
}

// UpdateProfile allows changing profile details (name, email, avatar base64, password)
func UpdateProfile(c *gin.Context) {
	userVal, _ := c.Get("user")
	user := userVal.(models.User)

	type UpdateReq struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		AvatarURL       string `json:"avatar_url"`
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	var req UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ข้อมูลไม่ถูกต้อง",
			"code":  "BAD_REQUEST",
		})
		return
	}

	// If email changes, check uniqueness
	if req.Email != "" && req.Email != user.Email {
		var checkUser models.User
		if err := database.DB.Where("email = ?", req.Email).First(&checkUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": "อีเมลนี้ถูกใช้งานโดยผู้อื่นแล้ว",
				"code":  "EMAIL_ALREADY_EXISTS",
			})
			return
		}
		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	// Check if updating password
	if req.NewPassword != "" {
		if req.CurrentPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "กรุณากรอกรหัสผ่านปัจจุบันเพื่อเปลี่ยนรหัสผ่านใหม่",
				"code":  "PASSWORD_REQUIRED",
			})
			return
		}
		// Validate current password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "รหัสผ่านปัจจุบันไม่ถูกต้อง",
				"code":  "INVALID_PASSWORD",
			})
			return
		}
		// Hash new password
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ไม่สามารถเข้ารหัสผ่านใหม่ได้",
				"code":  "INTERNAL_SERVER_ERROR",
			})
			return
		}
		user.Password = string(hashed)
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถบันทึกข้อมูลได้",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// DeleteAccount deletes all user data from DB
func DeleteAccount(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	// Wipe all workspaces, tasks, tags, comments, and finally the user
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// GORM Cascade constraints will handle tasks/comments on Workspaces / Tasks,
	// but let's delete explicitly to be sure.
	tx.Where("user_id = ?", userID).Delete(&models.Comment{})
	tx.Where("user_id = ?", userID).Delete(&models.Tag{})
	tx.Where("user_id = ?", userID).Delete(&models.Task{})
	tx.Where("user_id = ?", userID).Delete(&models.Workspace{})
	if err := tx.Delete(&models.User{}, userID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถลบบัญชีผู้ใช้ได้: " + err.Error(),
			"code":  "DATABASE_ERROR",
		})
		return
	}

	tx.Commit()
	clearTokenCookies(c)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"message": "ลบบัญชีผู้ใช้เรียบร้อยแล้ว",
		},
	})
}
