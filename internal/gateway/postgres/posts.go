package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/slinkeres/ozontask/graph/models"
)

type PostsPostgres struct {
	db *sqlx.DB
}

func NewPostsPostgres(db *sqlx.DB) *PostsPostgres {
	return &PostsPostgres{db: db}
}

func (p PostsPostgres) CreatePost(post models.Post) (models.Post, error) {

	query := `INSERT INTO Posts (name, content, author, comments_allowed) 
				VALUES ($1, $2, $3, $4)
				RETURNING id, created_at`

	tx, err := p.db.Begin()
	if err != nil {
		return models.Post{}, err
	}

	row := tx.QueryRow(query, post.Name, post.Content, post.Author, post.CommentsAllowed)
	if err := row.Scan(&post.ID, &post.CreatedAt); err != nil {
		tx.Rollback()
		return models.Post{}, err
	}

	return post, tx.Commit()
}

func (p PostsPostgres) GetPostById(id int) (models.Post, error) {

	query := `SELECT * FROM posts WHERE id = $1`

	var post models.Post

	if err := p.db.Get(&post, query, id); err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (p PostsPostgres) GetAllPosts(limit, offset int) ([]models.Post, error) {

	query := "SELECT * FROM posts ORDER BY created_at OFFSET $1"
	args := []interface{}{offset}

	if limit > 0 {
		query += " LIMIT $2"
		args = append(args, limit)
	}

	var posts []models.Post

	if err := p.db.Select(&posts, query, args...); err != nil {
		return nil, err
	}

	return posts, nil
}