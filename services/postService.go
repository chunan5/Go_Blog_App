package services

import (
	"Go/dto"
	"Go/models"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type PostService struct {
	DB *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{DB: db}
}

func (pc *PostService) Create(createPostDto dto.CreatePostDto, userId int) error {

	var post models.Post

	post = models.Post{
		Title:       createPostDto.Title,
		Description: createPostDto.Description,
		UserId:      userId,
	}

	if err := pc.DB.Create(&post).Error; err != nil {
		return err
	}

	return nil
}

func (pc *PostService) View(postUuid string) (models.Post, error) {

	var post models.Post
	err := pc.DB.Where("uuid = ?", postUuid).First(&post)

	if err.Error != nil {
		return models.Post{}, err.Error
	}

	return post, nil
}

func (pc *PostService) Delete(postUuid string) error {
	var post models.Post
	err := pc.DB.Where("uuid = ?", postUuid).First(&post)

	if err.Error != nil {
		return err.Error
	}

	if err := pc.DB.Delete(&post).Error; err != nil {
		return err
	}

	return nil
}

func (pc *PostService) Update(updatePostDto dto.UpdatePostDto) error {

	var post models.Post

	err := pc.DB.Where("uuid=?", updatePostDto.UUID).First(&post)

	if err.Error != nil {
		return err.Error
	}

	post.Title = updatePostDto.Title
	post.Description = updatePostDto.Description

	if err := pc.DB.Save(&post).Error; err != nil {
		return err
	}
	return nil
}

func (pc *PostService) List(listPostDto dto.ListPostDto) ([]models.Post, error) {

	var posts []models.Post

	query := pc.DB.Model(&models.Post{}).Joins("JOIN users ON users.id = posts.user_id").
		Select("posts.title, users.name as author, posts.created_at, posts.updated_at")

	if listPostDto.Title != nil && *listPostDto.Title != "" {
		query = query.Where("title LIKE ?", "%"+*listPostDto.Title+"%")
	}

	if listPostDto.AuthorName != nil && *listPostDto.AuthorName != "" {
		query = query.Where("users.name = ?", *listPostDto.AuthorName)
	}

	if listPostDto.CreatedAt != nil {
		date, err := time.Parse("2006-01-02", *listPostDto.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid date format")
		}
		query = query.Where("DATE(posts.created_at) = ?", date.Format("2006-01-02"))
	}

	err := query.Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}
