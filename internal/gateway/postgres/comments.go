package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/slinkeres/ozontask/graph/models"
)

type CommentsPostgres struct {
	db *sqlx.DB
}

func NewCommentsPostgres(db *sqlx.DB) *CommentsPostgres {
	return &CommentsPostgres{db: db}
}

func (c CommentsPostgres) CreateComment(comment models.Comment) (models.Comment, error) {

	tx, err := c.db.Begin()
	if err != nil {
		return models.Comment{}, err
	}

	query := `INSERT INTO comments (content, author, post, reply_to) 
				VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	row := tx.QueryRow(query, comment.Content, comment.Author, comment.Post, comment.ReplyTo)
	if err := row.Scan(&comment.ID, &comment.CreatedAt); err != nil {
		tx.Rollback()
		return models.Comment{}, err
	}

	return comment, tx.Commit()

}

func (c CommentsPostgres) GetCommentsByPost(postId, limit, offset int) ([]*models.Comment, error) {

	query := `SELECT * FROM comments 
         WHERE post = $1 AND reply_to IS NULL 
         ORDER BY created_at 
         OFFSET $2`

	args := []interface{}{postId, offset}

	if limit >= 0 {
		query += " LIMIT $3"
		args = append(args, limit)
	}

	var comments []*models.Comment

	if err := c.db.Select(&comments, query, args...); err != nil {
		return nil, err
	}

	return comments, nil
}

func (c CommentsPostgres) GetRepliesOfComment(commentId int) ([]*models.Comment, error) {

	query := `SELECT * FROM comments WHERE reply_to = $1`

	var comments []*models.Comment

	if err := c.db.Select(&comments, query, commentId); err != nil {
		return nil, err
	}

	return comments, nil
}