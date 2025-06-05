package service

import (
	"database/sql"
	"errors"

	"github.com/slinkeres/ozontask/graph/models"
	"github.com/slinkeres/ozontask/internal/consts"
	"github.com/slinkeres/ozontask/internal/gateway"
	"github.com/slinkeres/ozontask/internal/logger"
	"github.com/slinkeres/ozontask/internal/pagination"
	ce "github.com/slinkeres/ozontask/internal/custom_err"
)

type PostsService struct {
	repo   gateway.Posts
	logger *logger.Logger
}

func NewPostsService(repo gateway.Posts, logger *logger.Logger) *PostsService {
	return &PostsService{repo: repo, logger: logger}
}

func (p PostsService) CreatePost(post models.Post) (models.Post, error) {

	if len(post.Author) == 0 {
		p.logger.Err.Println("empty author error")
		return models.Post{}, ce.CustomError{
			Message: "empty author error",
			Type:    "Bad Request",
		}
	}

	if len(post.Content) >= consts.MaxContentLength {
		p.logger.Err.Println("content is too long", len(post.Content))
		return models.Post{}, ce.CustomError{
			Message: "content is too long",
			Type:    "Bad Request",
		}
	}

	newPost, err := p.repo.CreatePost(post)
	if err != nil {
		p.logger.Err.Println("error with creating new post", err.Error())
		return models.Post{}, ce.CustomError{
			Message: "error with creating new post",
			Type:    "Internal Server Error",
		}
	}

	return newPost, nil

}

func (p PostsService) GetPostById(postId int) (models.Post, error) {

	if postId <= 0 {
		p.logger.Err.Println("wrong id error", postId)
		return models.Post{}, ce.CustomError{
			Message: "wrong id error",
			Type:    "Bad Request",
		}
	}

	post, err := p.repo.GetPostById(postId)
	if err != nil {

		p.logger.Err.Println("error with getting post", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return models.Post{}, ce.CustomError{
				Message: "post not foud",
				Type:    "Not Found Error",
			}
		}
		return models.Post{}, ce.CustomError{
			Message: "error with getting post",
			Type:    "Internal Server Error",
		}
	}

	return post, nil
}

func (p PostsService) GetAllPosts(page, pageSize *int) ([]models.Post, error) {

	if page != nil && *page < 0 {
		p.logger.Err.Println("wrong page number error", *page)
		return nil, ce.CustomError{
			Message: "wrong page number error",
			Type:    "Bad Request",
		}
	}

	if pageSize != nil && *pageSize < 0 {
		p.logger.Err.Println("wrong page size error", *pageSize)
		return nil, ce.CustomError{
			Message: "wrong page size error",
			Type:    "Bad Request",
		}
	}

	offset, limit := pagination.GetOffsetAndLimit(page, pageSize)

	posts, err := p.repo.GetAllPosts(limit, offset)
	if err != nil {
		p.logger.Err.Println("error with getting post", err.Error())
		return nil, ce.CustomError{
			Message: "error with getting post",
			Type:    "Internal Server Error",
		}
	}

	return posts, nil
}