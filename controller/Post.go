package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaoyuer1231231/gin_mysql_grom_project/models"
	"gorm.io/gorm"
)

type PostController struct {
	DB *gorm.DB
}

func NewPostController(db *gorm.DB) *PostController {
	return &PostController{DB: db}
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// CreatePost 创建文章
// @Summary      创建文章
// @Description  创建新的文章
// @Tags         文章管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreatePostRequest true "文章信息"
// @Success      200  {object} models.Response{data=models.Post}
// @Failure      400  {object} models.Response
// @Failure      401  {object} models.Response
// @Router       /post/createPost [post]
func (pc *PostController) CreatePost(c *gin.Context) {
	userID := c.GetUint("user_id")
	fmt.Println("CreatePost-----------------------------", userID)

	var CreatePostRequest struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&CreatePostRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": http.StatusInternalServerError})
		return
	}
	post := models.Post{
		UserID:  userID,
		Title:   CreatePostRequest.Title,
		Content: CreatePostRequest.Content,
	}
	result := pc.DB.Create(&post)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error(), "code": http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"code":    http.StatusOK,
	})
}

// 获取所有文章信息和单个文章的信息
func (pc *PostController) QueryPost(c *gin.Context) {
	var posts []models.Post
	postID := c.Query("id") // 返回 "123"
	fmt.Println("QueryPost------id:", postID)
	// 如果有 ID 参数，查询单篇文章
	if postID != "" {
		id, err := strconv.ParseUint(postID, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID", "code": http.StatusInternalServerError})
			return
		}

		result := pc.DB.Find(&posts, uint(id))
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败", "code": http.StatusInternalServerError})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"post": posts})
		return // 重要：这里要 return，避免执行下面的代码
	}

	// 如果没有 ID 参数，查询所有文章
	result := pc.DB.Order("created_at DESC").Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败", "code": http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}
func (pc *PostController) UptDateById(c *gin.Context) {
	userId := c.GetUint("user_id")
	var updatePostRequest struct {
		PostId  string `json:"postId" binding:"required"`
		Title   string `json:"title" `
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&updatePostRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": http.StatusInternalServerError})
		return
	}

	//查询数据是否存在
	resultQuery := pc.DB.Model(&models.Post{}).Where("id = ? and user_id =?", updatePostRequest.PostId, userId).First(&models.Post{})
	if resultQuery.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据不存在", "code": http.StatusInternalServerError})
		return
	}
	fmt.Println("QueryPost------id:")
	updates := make(map[string]interface{})

	if updatePostRequest.Content != "" {
		updates["content"] = updatePostRequest.Content
	}
	if updatePostRequest.Title != "" {
		updates["title"] = updatePostRequest.Title
	}
	resultUpdate := pc.DB.Debug().Model(&models.Post{}).Where("id = ?", updatePostRequest.PostId).Updates(updates)
	if resultUpdate.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "", "code": http.StatusInternalServerError})
		return

	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"code":    http.StatusOK,
	})
}

func (pc *PostController) DeleteById(c *gin.Context) {
	userId := c.GetUint("user_id")
	PostId := c.Query("id")

	//查询数据是否存在
	var post models.Post
	resultQuery := pc.DB.Debug().Model(&models.Post{}).Where("id = ? and user_id =?", PostId, userId).First(&post)
	if resultQuery.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数据不存在", "code": http.StatusInternalServerError})
		return
	}
	resultMsg := pc.DB.Model(&post).Delete(&post)
	if resultMsg.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post", "code": http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
