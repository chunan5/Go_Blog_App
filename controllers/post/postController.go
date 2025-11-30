package post

import (
	"Go/dto"
	"Go/models"
	"Go/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type PostController struct {
	PostService *services.PostService
}

func NewPostController(postService *services.PostService) *PostController {
	return &PostController{PostService: postService}
}

func (pc *PostController) Create(c *gin.Context) {
	var postDto dto.CreatePostDto

	if err := c.ShouldBindJSON(&postDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User type invalid"})
		return
	}

	err := pc.PostService.Create(postDto, int(user.ID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Post Created!!"})

	}

}

func (pc *PostController) View(c *gin.Context) {
	uuid := c.Param("uuid")

	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid cannot be empty"})
		return
	}

	post, err := pc.PostService.View(uuid)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Post fetched Successfully!", "data": post})
	}
}

func (pc *PostController) Delete(c *gin.Context) {
	uuid := c.Param("uuid")

	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid cannot be empty"})
		return
	}

	err := pc.PostService.Delete(uuid)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Post Deleted Successfully!"})
	}

}

func (pc *PostController) Update(c *gin.Context) {
	var postDto dto.UpdatePostDto

	if err := c.ShouldBindJSON(&postDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := pc.PostService.Update(postDto)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Post Updated Successfully!"})
	}

}

func (pc *PostController) List(c *gin.Context) {
	var listPostDto dto.ListPostDto

	if err := c.ShouldBindJSON(&listPostDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	posts, err := pc.PostService.List(listPostDto)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "posts not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Posts List fetched Successfully!", "data": posts})
	}

}
