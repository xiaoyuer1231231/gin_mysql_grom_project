package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiaoyuer1231231/gin_mysql_grom_project/models"
	"gorm.io/gorm"
)

type CommentController struct {
	DB *gorm.DB
}

func NewCommentController(db *gorm.DB) *CommentController {
	return &CommentController{DB: db}
}

type CommentCreateRequest struct {
	PostId  uint   `json:"post_id" binding:"required"` // ✅ 正确：binding:"required"
	Content string `json:"content" binding:"required"`
}

func (cc *CommentController) CreateComment(c *gin.Context) {
	userId := c.GetUint("user_id")
	var reg CommentCreateRequest
	if err := c.ShouldBindJSON(&reg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": http.StatusInternalServerError})
		return
	}
	// 查询文章是否存在
	var post models.Post
	resultPost := cc.DB.Model(&models.Post{}).Where("id = ?", reg.PostId).Find(&post)
	if resultPost.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章数据不存在无法评论"})
		return
	}
	fmt.Println("1111", reg.PostId)
	var comment models.Comment
	comment.Content = reg.Content
	comment.PostID = reg.PostId
	comment.UserID = userId
	resultCommentPost := cc.DB.Model(&models.Comment{}).Create(&comment)
	if resultCommentPost.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论失败", "code": http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully",
		"code": http.StatusOK})
}

func (cc *CommentController) QueryComment(c *gin.Context) {
	//userId := c.GetUint("user_id")
	postId := c.Query("postId")
	fmt.Printf("postId:%s\n", postId)
	var comments []models.Comment
	resultMsg := cc.DB.Debug().Model(&models.Comment{}).Where("post_id =? ", postId).Find(&comments)
	if resultMsg.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "查询数据失败", "code": http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully",
		"code": http.StatusOK,
		"data": comments})
}
