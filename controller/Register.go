package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiaoyuer1231231/gin_mysql_grom_project/config"
	"github.com/xiaoyuer1231231/gin_mysql_grom_project/models"
	"github.com/xiaoyuer1231231/gin_mysql_grom_project/util"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB     *gorm.DB
	config *config.Config
}

func NewAuthController(db *gorm.DB, config *config.Config) *AuthHandler {
	return &AuthHandler{DB: db, config: config}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": http.StatusInternalServerError})
		return
	}
	//
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
	}
	// 查询人员是否存在
	fmt.Printf("---------------", user)
	var existingUser models.User
	resultFindUser := h.DB.Where("username = ? ", user.Username).First(&existingUser)
	if resultFindUser.RowsAffected >= 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户已存在", "code": http.StatusInternalServerError})
		return
	}
	// 加密密码
	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败", "code": http.StatusInternalServerError})
		return
	}
	user.Password = hashedPassword

	// 创建用户

	resultMsg := h.DB.Debug().Create(&user)
	if resultMsg.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": resultMsg.Error.Error(), "code": http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})

}

func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": http.StatusInternalServerError})
		return
	}
	// 查询账号是否存在
	var user models.User
	resultFindUser := h.DB.Debug().Where("Username = ? ", input.Username).First(&user)
	if resultFindUser.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户不存在", "code": http.StatusInternalServerError})
		return
	}

	hashPassWord := user.Password
	fmt.Println("数据库的密码", hashPassWord)
	if !util.CheckPasswordHash(input.Password, hashPassWord) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码不正确", "code": http.StatusInternalServerError})
		return
	}
	//获取token
	fmt.Println("数据库的密码config", h.config.JWT.Secret)

	token, err := util.GenerateToken(user.ID, h.config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"获取token失败": err.Error(), "code": http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
