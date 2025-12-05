// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"kk-nav/internal/models"
	"kk-nav/internal/utils"
	"gorm.io/gorm"
)

// UsersHandler 管理后台用户处理器
type UsersHandler struct {
	db *gorm.DB
}

// NewUsersHandler 创建用户处理器
func NewUsersHandler(db *gorm.DB) *UsersHandler {
	return &UsersHandler{db: db}
}

// Index 用户列表
func (h *UsersHandler) Index(c *gin.Context) {
	var users []models.User

	// 搜索
	query := h.db
	if search := c.Query("search"); search != "" {
		query = query.Where("email ILIKE ? OR username ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 角色筛选
	if role := c.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}

	// 激活状态筛选
	if active := c.Query("active"); active != "" {
		if active == "true" {
			query = query.Where("active = ?", true)
		} else if active == "false" {
			query = query.Where("active = ?", false)
		}
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var total int64
	query.Model(&models.User{}).Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users)

	// 移除密码哈希
	for i := range users {
		users[i].PasswordHash = ""
	}

	utils.Success(c, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (int(total) + pageSize - 1) / pageSize,
		},
	})
}

// Create 创建用户
func (h *UsersHandler) Create(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Username string `json:"username" binding:"required,min=3,max=100"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role"`
		Active   bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 检查邮箱是否已存在
	var existing models.User
	if err := h.db.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		utils.Error(c, 400, "Email already exists")
		return
	}

	// 检查用户名是否已存在
	if err := h.db.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		utils.Error(c, 400, "Username already exists")
		return
	}

	// 加密密码
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalServerError(c, "Failed to hash password")
		return
	}

	// 设置默认值
	role := req.Role
	if role == "" {
		role = "user"
	}

	user := models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: passwordHash,
		Role:         role,
		Active:       req.Active,
	}

	if err := h.db.Create(&user).Error; err != nil {
		utils.InternalServerError(c, "Failed to create user")
		return
	}

	user.PasswordHash = ""
	utils.SuccessWithMessage(c, "User created successfully", user)
}

// Update 更新用户
func (h *UsersHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFound(c, "User not found")
			return
		}
		utils.InternalServerError(c, "Database error")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
		Active   *bool  `json:"active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 更新字段
	if req.Email != "" {
		// 检查邮箱是否与其他用户冲突
		var existing models.User
		if err := h.db.Where("email = ? AND id != ?", req.Email, id).First(&existing).Error; err == nil {
			utils.Error(c, 400, "Email already exists")
			return
		}
		user.Email = req.Email
	}

	if req.Username != "" {
		// 检查用户名是否与其他用户冲突
		var existing models.User
		if err := h.db.Where("username = ? AND id != ?", req.Username, id).First(&existing).Error; err == nil {
			utils.Error(c, 400, "Username already exists")
			return
		}
		user.Username = req.Username
	}

	if req.Password != "" {
		// 更新密码
		passwordHash, err := utils.HashPassword(req.Password)
		if err != nil {
			utils.InternalServerError(c, "Failed to hash password")
			return
		}
		user.PasswordHash = passwordHash
	}

	if req.Role != "" {
		user.Role = req.Role
	}

	if req.Active != nil {
		user.Active = *req.Active
	}

	if err := h.db.Save(&user).Error; err != nil {
		utils.InternalServerError(c, "Failed to update user")
		return
	}

	user.PasswordHash = ""
	utils.SuccessWithMessage(c, "User updated successfully", user)
}

// Delete 删除用户
func (h *UsersHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	// 不能删除自己
	currentUserID, _ := c.Get("user_id")
	if currentUserID.(uint) == uint(id) {
		utils.Error(c, 400, "Cannot delete yourself")
		return
	}

	if err := h.db.Delete(&models.User{}, id).Error; err != nil {
		utils.InternalServerError(c, "Failed to delete user")
		return
	}

	utils.SuccessWithMessage(c, "User deleted successfully", nil)
}

