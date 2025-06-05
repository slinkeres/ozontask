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

type CommentsService struct {
	repo       gateway.Comments
	logger     *logger.Logger
	PostGetter PostGetter
}

type PostGetter interface {
	GetPostById(id int) (models.Post, error)
}

func NewCommentsService(repo gateway.Comments, logger *logger.Logger, getter PostGetter) *CommentsService {
	return &CommentsService{repo: repo, logger: logger, PostGetter: getter}
}

func (c CommentsService) CreateComment(comment models.Comment) (models.Comment, error) {
	if len(comment.Author) == 0 {
		c.logger.Err.Println("empty author error")
		return models.Comment{}, ce.CustomError{
			Message: "empty author error",
			Type:    "Bad Request",
		}
	}

	if len(comment.Content) >= consts.MaxContentLength {
		c.logger.Err.Println("content is too long", len(comment.Content))
		return models.Comment{}, ce.CustomError{
			Message: "content is too long",
			Type:    "Bad Request",
		}
	}

	if comment.Post <= 0 {
		c.logger.Err.Println("wrong id error", comment.Post)
		return models.Comment{}, ce.CustomError{
			Message: "wrong id error",
			Type:    "Bad Request",
		}
	}

	post, err := c.PostGetter.GetPostById(comment.Post)
	if err != nil {
		c.logger.Err.Println("error with getting post", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return models.Comment{}, ce.CustomError{
				Message: "post not found",
				Type:    "Not Found Error",
			}
		}
	}

	if !post.CommentsAllowed {
		c.logger.Err.Println("comments not allowed for this post")
		return models.Comment{}, ce.CustomError{
			Message: "comments not allowed for this post",
			Type:    "Bad Request",
		}

	}

	newComment, err := c.repo.CreateComment(comment)
	if err != nil {
		c.logger.Err.Println("error with creating new comment", err.Error())
		return models.Comment{}, ce.CustomError{
			Message: "error with creating new comment",
			Type:    "Internal Server Error",
		}
	}

	return newComment, nil
}

func (c CommentsService) GetCommentsByPost(postId int, page *int, pageSize *int) ([]*models.Comment, error) {

	if postId <= 0 {
		c.logger.Err.Println("wrong id error", postId)
		return nil, ce.CustomError{
			Message: "wrong id error",
			Type:    "Bad Request",
		}
	}

	if page != nil && *page < 0 {
		c.logger.Err.Println("wrong page number error", *page)
		return nil, ce.CustomError{
			Message: "wrong page number error",
			Type:    "Bad Request",
		}
	}

	if pageSize != nil && *pageSize < 0 {
		c.logger.Err.Println("wrong page size error", *pageSize)
		return nil, ce.CustomError{
			Message: "wrong page size error",
			Type:    "Bad Request",
		}
	}

	offset, limit := pagination.GetOffsetAndLimit(page, pageSize)

	comments, err := c.repo.GetCommentsByPost(postId, limit, offset)
	if err != nil {
		c.logger.Err.Println("error with getting comments", postId, err.Error())
		return nil, ce.CustomError{
			Message: "error with getting comments",
			Type:    "Internal Server Error",
		}
	}

	return comments, nil
}

func (c CommentsService) GetRepliesOfComment(commentId int) ([]*models.Comment, error) {

	if commentId <= 0 {
		c.logger.Err.Println("wrong id error", commentId)
		return nil, ce.CustomError{
			Message: "wrong id error",
			Type:    "Bad Request",
		}
	}

	comments, err := c.repo.GetRepliesOfComment(commentId)
	if err != nil {
		c.logger.Err.Println("error with getting replies ti comment", commentId, err.Error())
		return nil, ce.CustomError{
			Message: "error with getting replies ti comment",
			Type:    "Internal Server Error",
		}
	}

	return comments, nil

}