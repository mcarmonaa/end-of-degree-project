package auth

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
)

type Service struct {
	nonces *nonces
	db     *gorm.DB
}

func NewAuthSvc(db *gorm.DB) *Service {
	return &Service{
		nonces: newNonces(100),
		db:     db,
	}

}

func (s *Service) Register(context.Context, *RegisterRequest) (*RegisterReply, error) {

	return nil, nil
}

func (s *Service) GetSalt(context.Context, *SaltRequest) (*SaltReply, error) {
	return nil, nil
}

func (s *Service) Login(context.Context, *LoginRequest) (*LoginReply, error) {
	return nil, nil
}
