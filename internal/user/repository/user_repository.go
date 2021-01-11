package repository

import (
	"database/sql"
	"fmt"
	"github.com/SanSanchezzz/tech-db-forum/internal/models"
	"github.com/SanSanchezzz/tech-db-forum/internal/user"
)

type UserRepository struct {
	dbConn *sql.DB
}

func NewUserRepository(conn *sql.DB) user.UserRepository{
	return &UserRepository{
		dbConn: conn,
	}
}

func (ur *UserRepository) SelectUsersByForum(limit int, desc bool, slug, since *string) ([]*models.User, error) {
	query := `select users.nickname, users.fullname, users.about, users.email from forum_user join users on forum_user.nickname = users.nickname WHERE forum_user.forum = $1 `

	if desc {
		if *since != "" {
			query += fmt.Sprintf(`AND users.nickname < '%s' `, *since)
		}
		query += `ORDER BY users.nickname desc `
	} else {
		if *since != "" {
			query += fmt.Sprintf(`AND users.nickname > '%s' `, *since)
		}
		query += `ORDER BY users.nickname `
	}

	query += fmt.Sprintf(`limit %d`, limit)

	rows, err := ur.dbConn.Query(query, slug)

	if err != nil {
		return nil, err
	}

	users :=  make([]*models.User, 0)
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) GetByNicknameOrEmail(nickname, email string) ([]*models.User, error) {
	rows, err := ur.dbConn.Query(
		`select nickname, fullname, about, email from users where nickname = $1 or email = $2`,
		nickname, email,
		)

	if err != nil {
		return nil, err
	}

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
			)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) Insert(nickname, fullname, about, email string) error {
	_, err := ur.dbConn.Exec(
		`insert into users (nickname, fullname, about, email) values ($1, $2, $3, $4)`,
		nickname, fullname, about, email,
		)

	if err != nil {
		return err
	}

	return nil
}

func (ur * UserRepository) SelectByNickname(nickname string) (*models.User, error) {
	user := &models.User{}

	err := ur.dbConn.QueryRow(
		`select nickname, fullname, about, email from users where nickname = $1`,
		nickname,
		).Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
			)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) Update(user *models.User) error {
	_, err := ur.dbConn.Exec(
		`update users set fullname = $1, about = $2, email = $3 where nickname = $4`,
		user.Fullname, user.About, user.Email, user.Nickname,
	)

	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) SelectUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	err := ur.dbConn.QueryRow(
		`select nickname, fullname, about, email from users where email = $1`,
		email,
		).Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
			)

	if err != nil {
		return nil, err
	}

	return user, nil
}
