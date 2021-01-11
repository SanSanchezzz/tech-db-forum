package repository

import (
	"database/sql"
	"fmt"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
	"github.com/SanSanchezzz/tech-db-forum/internal/post"
	"github.com/lib/pq"
	"strconv"
	"time"
)

type PostRepository struct {
	dbConn *sql.DB
}

func NewPostRepository(conn *sql.DB) post.PostRepository{
	return &PostRepository{
		dbConn: conn,
	}
}

func (pr *PostRepository) GetPost(id int) (*models.Post, error) {
	post := &models.Post{}
	err := pr.dbConn.QueryRow(
		`select id, author, message, is_edited, parent,  forum, thread, created from posts where id = $1`,
		id,
		).Scan(
		&post.ID,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Parent,
		&post.Forum,
		&post.Thread,
		&post.Created,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}


func (pr *PostRepository) GetPosts(limit, since *string, id uint32) ([]*models.Post, error) {
	lim, _ := strconv.Atoi(*limit)
	query := `select id, message, is_edited, created, parent, author, forum, thread from posts where thread = $1 `

	if *since != "" {
		query += fmt.Sprintf("and id > %s ", *since)
	}

	query += "order by created, id  "

	if lim > 0 {
		query += fmt.Sprintf("limit %s;", *limit)
	}

	rows, err := pr.dbConn.Query(query, id)
	if err != nil {
		return nil, err
	}

	posts := make([]*models.Post, 0)

	for rows.Next() {
		post := &models.Post{}
		if err := rows.Scan(
			&post.ID,
			&post.Message,
			&post.IsEdited,
			&post.Created,
			&post.Parent,
			&post.Author,
			&post.Forum,
			&post.Thread,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepository) GetPostsDesc(limit, since *string, id uint32) ([]*models.Post, error) {
	lim, _ := strconv.Atoi(*limit)
	query := `select id, message, is_edited, created, parent, author, forum, thread from posts where thread = $1 `

	if *since != "" {
		query += fmt.Sprintf("and id < %s ", *since)
	}

	query += "order by created desc, id desc "

	if lim > 0 {
		query += fmt.Sprintf("limit %s;", *limit)
	}

	rows, err := pr.dbConn.Query(query, id)
	if err != nil {
		return nil, err
	}

	posts := make([]*models.Post, 0)

	for rows.Next() {
		post := &models.Post{}
		if err := rows.Scan(
			&post.ID,
			&post.Message,
			&post.IsEdited,
			&post.Created,
			&post.Parent,
			&post.Author,
			&post.Forum,
			&post.Thread,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepository) GetPostsTree(limit, since *string, id uint32) ([]*models.Post, error) {
	lim, _ := strconv.Atoi(*limit)
	query := `select id, message, is_edited, created, parent, author, forum, thread from posts where thread = $1 `

	if *since != "" {
		query += fmt.Sprintf("and path > (select path from posts where id = %s) ", *since)
	}

	query += "order by path, id  "

	if lim > 0 {
		query += fmt.Sprintf("limit %s;", *limit)
	}

	rows, err := pr.dbConn.Query(query, id)
	if err != nil {
		return nil, err
	}

	posts := make([]*models.Post, 0)

	for rows.Next() {
		post := &models.Post{}
		if err := rows.Scan(
			&post.ID,
			&post.Message,
			&post.IsEdited,
			&post.Created,
			&post.Parent,
			&post.Author,
			&post.Forum,
			&post.Thread,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepository) GetPostsDescTree(limit, since *string, id uint32) ([]*models.Post, error) {
	lim, _ := strconv.Atoi(*limit)
	query := `select id, message, is_edited, created, parent, author, forum, thread from posts where thread = $1 `

	if *since != "" {
		query += fmt.Sprintf("and path < (select path from posts where id = %s) ", *since)
	}

	query += "order by path desc, id desc "

	if lim > 0 {
		query += fmt.Sprintf("limit %s;", *limit)
	}

	rows, err := pr.dbConn.Query(query, id)
	if err != nil {
		return nil, err
	}

	posts := make([]*models.Post, 0)

	for rows.Next() {
		post := &models.Post{}
		if err := rows.Scan(
			&post.ID,
			&post.Message,
			&post.IsEdited,
			&post.Created,
			&post.Parent,
			&post.Author,
			&post.Forum,
			&post.Thread,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepository) GetPostsDescParentTree(limit, since *string, id uint32) ([]*models.Post, error) {
	lim, _ := strconv.Atoi(*limit)
	query := `select id, message, is_edited, created, parent, author, forum, thread from posts where thread = $1 `

	if *since != "" {
		query += fmt.Sprintf("and path[1] in (select distinct path[1] from posts where path[1] < (select path[1] from posts where id = %s) and array_length(path, 1) = 1 and thread = $1 order by path desc ", *since)
	} else {
		query += "and path[1] in (select distinct path[1] from posts where array_length(path, 1) = 1 and thread = $1 order by path[1] desc "

	}

	if lim > 0 {
		query += fmt.Sprintf("limit %s ) ", *limit)
	}

	query += "order by path[1] desc, path[2:], created, id;"


	rows, err := pr.dbConn.Query(query, id)
	if err != nil {
		return nil, err
	}

	posts := make([]*models.Post, 0)

	for rows.Next() {
		post := &models.Post{}
		if err := rows.Scan(
			&post.ID,
			&post.Message,
			&post.IsEdited,
			&post.Created,
			&post.Parent,
			&post.Author,
			&post.Forum,
			&post.Thread,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepository) GetPostsParentTree(limit, since *string, id uint32) ([]*models.Post, error) {
	lim, _ := strconv.Atoi(*limit)
	query := `select id, message, is_edited, created, parent, author, forum, thread from posts where thread = $1 `

	if *since != "" {
		query += fmt.Sprintf("and path[1] in (select distinct path[1] from posts where path[1] > (select path[1] from posts where id = %s) and array_length(path, 1) = 1 and thread = $1 order by path ", *since)
	} else {
		query += "and path[1] in (select distinct path[1] from posts where array_length(path, 1) = 1 and thread = $1 order by path[1] "
	}

	if lim > 0 {
		query += fmt.Sprintf("limit %s ) ", *limit)
	}

	query += "order by path[1], path[2:], created, id; "

	rows, err := pr.dbConn.Query(query, id)
	if err != nil {
		return nil, err
	}

	posts := make([]*models.Post, 0)

	for rows.Next() {
		post := &models.Post{}
		if err := rows.Scan(
			&post.ID,
			&post.Message,
			&post.IsEdited,
			&post.Created,
			&post.Parent,
			&post.Author,
			&post.Forum,
			&post.Thread,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepository) Insert(posts []*models.Post, thread *models.Thread, now *time.Time) error {
	tx, err := pr.dbConn.Begin()

	for _, post := range posts {
		err = tx.QueryRow(`insert into posts(id, message, is_edited, created, parent, author, thread, forum, path)
			values(default, $1, default, $2, $3, $4, $5, $6, $7) returning id`,
			post.Message,
			now,
			post.Parent,
			post.Author,
			thread.ID,
			thread.Forum,
			pq.Array([]int64{int64(post.Parent)})).
			Scan(
				&post.ID,
			)

		if err != nil {
			tx.Rollback()
			return err
		}

		post.Created = *now
		post.Thread = thread.ID
		post.Forum = thread.Forum
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (pr *PostRepository) Update(post *models.Post) error {
	_, err := pr.dbConn.Exec(
		`update posts set message = $1, is_edited = true where id = $2`,
		post.Message,
		post.ID,
		)

	if err != nil {
		return err
	}

	post.IsEdited = true

	return nil
}
