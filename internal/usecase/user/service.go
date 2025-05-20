package authusecase

import (
	"time"

	"github.com/user2083251241/ebidsystem/internal/domain/entity"
	"github.com/user2083251241/ebidsystem/internal/interfaces/http/dto"
)

type UseCase interface {
	Register(req dto.RegisterRequest) (*entity.User, error)
	Login(req dto.LoginRequest) (*entity.User, error)
	Logout(userID uint, token string, expiration time.Duration) error
}
