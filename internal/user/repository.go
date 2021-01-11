package user

import "github.com/SanSanchezzz/tech-db-forum/internal/models"

type UserRepository interface {
	GetByNicknameOrEmail(nickname, email string) ([]*models.User, error)
	Insert(nickname, fullname, about, email string) error
	SelectByNickname(nickname string) (*models.User, error)
	Update(user *models.User) error
	SelectUserByEmail(email string) (*models.User, error)
	SelectUsersByForum(limit int, desc bool, slug, since *string) ([]*models.User, error)

}
