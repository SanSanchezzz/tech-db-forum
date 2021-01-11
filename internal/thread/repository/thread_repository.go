package repository

import (
	"database/sql"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
	"github.com/SanSanchezzz/tech-db-forum/internal/thread"
	"strconv"
)

type ThreadRepository struct {
	dbConn *sql.DB
}

func NewThreadRepository(conn *sql.DB) thread.ThreadRepository{
	return &ThreadRepository{
		dbConn: conn,
	}
}

func (tr *ThreadRepository) SelectBySlug(slug string) (*models.Thread, error) { // TODO: передавать по ссылке
	thread := &models.Thread{}

	err := tr.dbConn.QueryRow(
		`select id, title, nickname, forum, message, votes, slug, created from threads where slug = $1`,
		slug,
	).Scan(
		&thread.ID,
		&thread.Title,
		&thread.Nickname,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created,
	)

	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (tr *ThreadRepository) SelectByID(id uint32) (*models.Thread, error) { // TODO: передавать по ссылке
	thread := &models.Thread{}

	err := tr.dbConn.QueryRow(
		`select id, title, nickname, forum, message, votes, slug, created from threads where id = $1`,
		id,
	).Scan(
		&thread.ID,
		&thread.Title,
		&thread.Nickname,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created,
	)

	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (tr *ThreadRepository) Insert(thread *models.Thread) (uint32, error) {
	var id uint32
	err := tr.dbConn.QueryRow(
		`insert into threads (id, title, nickname, forum, message, votes, slug, created) values (default, $1, $2, $3, $4, default, $5, $6) returning id`,
		thread.Title, thread.Nickname, thread.Forum, thread.Message, thread.Slug, thread.Created,
	).Scan(
		&id,
		)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (tr *ThreadRepository) Update(thread *models.Thread) error {
	_, err := tr.dbConn.Exec(
		`update threads set title = $1, message = $2 where id = $3`,
		thread.Title, thread.Message, thread.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (tr *ThreadRepository) SelectThreadsByParam(limit int, desc bool, slug, since *string) ([]*models.Thread, error) {
	query := `select id, title, nickname, forum, message, votes, slug, created from threads where forum = $1 `
	//if since != nil {
	//	query += `and created >= $2 `
	//}

	if desc {
		if *since != "" {
			query += `and created <= $2 `
		}
		query += `order by created desc `
	} else {
		if *since != "" {
			query += `and created >= $2 `
		}
		query += `order by created `
	}

	if *since != "" {
		query += `limit $3`
	} else {
		query += `limit $2`
	}

	var rows *sql.Rows
	var err error

	if *since != "" {
		rows, err = tr.dbConn.Query(
			query,
			*slug, *since, limit,
		)
	} else {
		rows, err = tr.dbConn.Query(
			query,
			*slug, limit,
		)
	}

	if err != nil {
		return nil, err
	}

	threads := make([]*models.Thread, 0)
	for rows.Next() {
		thread := &models.Thread{}
		err := rows.Scan(
			&thread.ID,
			&thread.Title,
			&thread.Nickname,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
		)

		if err != nil {
			return nil, err
		}

		threads = append(threads, thread)
	}

	return threads, nil
}

func (tr *ThreadRepository) InsertVote(voice int, nickname, slugOrID string) (*models.Thread, error) {
	queryStart := `insert into votes (nickname, thread, voice) `
	queryEnd := `on conflict (nickname, thread) DO UPDATE SET voice = $3;`
	queryID := `values ($1, $2, $3) `
	querySlug := `select $1, id, $3 from threads where slug = $2  `
	querySelect := `select id, nickname, created, forum, message, slug, title, votes from threads where `

	if _, err := strconv.Atoi(slugOrID); err == nil {
		queryStart += queryID
		querySelect += `id = $1`
	} else {
		queryStart += querySlug
		querySelect += `slug = $1`
	}

	queryStart += queryEnd

	_, err := tr.dbConn.Exec(queryStart, nickname, slugOrID, voice)
	if err != nil {
		return nil, err
	}

	thread := &models.Thread{}

	err = tr.dbConn.QueryRow(querySelect,
		slugOrID).Scan(&thread.ID, &thread.Nickname, &thread.Created, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	if err != nil {
		return nil, err
	}

	return thread, nil
}
