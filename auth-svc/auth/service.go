package auth

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mcarmonaa/end-of-degree-project/auth-svc/nonces"
	"github.com/mcarmonaa/end-of-degree-project/message"

	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	svc             = "auth-svc: "
	alreadyExists   = svc + "already registered user "
	internal        = svc + "internal error"
	notFound        = svc + "not found user "
	unauthenticated = svc + "coulnd't authenticate user "
)

// Service represents the service exposes by auth-svc.
type Service struct {
	blackList *nonces.Nonces
	db        *gorm.DB
}

// NewAuthSvc creates a new Service connected against the given db creating a table "users" with the type Users as a model.
func NewAuthSvc(db *gorm.DB) (*Service, error) {
	if err := db.AutoMigrate(&User{}).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&User{}).AddUniqueIndex("idx_user_mail", "mail").Error; err != nil {
		return nil, err
	}

	return &Service{
		blackList: nonces.New(100),
		db:        db,
	}, nil

}

// Register creates a new user if the given mail doesn't exist in the database, otherwise it returns an error status.
func (s *Service) Register(c context.Context, r *RegisterRequest) (*RegisterReply, error) {
	mail := r.GetMail()
	if inDB(s.db, mail) {
		return nil, grpc.Errorf(codes.AlreadyExists, alreadyExists+"%q", mail)
	}

	salt, err := generateSalt()
	if err != nil {
		log.Println("register: couldn't generate salt: ", err)
		return nil, grpc.Errorf(codes.Internal, internal)
	}

	passKDF, err := generatePassword(r.GetPassword(), salt)
	if err != nil {
		log.Printf("register: couldn't create user %q: %v\n", mail, err)
		return nil, grpc.Errorf(codes.Internal, internal)
	}

	if err := s.db.Create(&User{Mail: mail, Password: passKDF, Salt: salt}).Error; err != nil {
		log.Printf("register: couldn't create user %q: %v\n", mail, err)
		return nil, grpc.Errorf(codes.Internal, internal)
	}

	return &RegisterReply{Salt: salt}, nil
}

// GetSalt returns the salt associated to the given mail.
func (s *Service) GetSalt(c context.Context, r *SaltRequest) (*SaltReply, error) {
	mail := r.GetMail()
	if !inDB(s.db, mail) {
		return nil, grpc.Errorf(codes.NotFound, notFound+"%q", mail)
	}

	var user User
	if err := s.db.Select("salt").Where("mail = ?", mail).First(&user).Error; err != nil {
		log.Printf("salt: couldn't get salt: %v\n", err)
		return nil, grpc.Errorf(codes.Internal, internal)
	}

	salt := user.Salt
	return &SaltReply{Salt: salt}, nil
}

// Login checks the given mail and the encrypted message(base64 encoded) with the user's kdf are correct and returns a JWE for the user in a
// encrypted message ciphered with the shared key send by the user.
//
// All encrypted data (with Ks or Ku) use AES256-GCM algorithm.
//
// C generates Ks randomly (An AES256 key).
//
// C request for its salt: C -> S GetSalt(mail)
// S reply with user's salt: S -> C salt
//
// C request for login: C -> S Login(mail, IV, Ku(TS, nonce, Ks) with Authenticated Data = mail) (the payload received would be plained in base64).
// if login success:
// S reply with a JWE: S -> C OK, IV, Ks(TS', nonce + 1, JWE) (IV and payload will be returns plained in base64).
// if login fails:
// S reply error: S -> C Error, Authentication failed
func (s *Service) Login(c context.Context, r *LoginRequest) (*LoginReply, error) {
	mail := r.GetMail()
	if !inDB(s.db, mail) {
		return nil, grpc.Errorf(codes.NotFound, notFound+"%q", mail)
	}

	var user User
	if err := s.db.Where("mail = ?", mail).First(&user).Error; err != nil {
		log.Printf("login: couldn't get password and salt: %v\n", err)
		return nil, grpc.Errorf(codes.Internal, internal)
	}

	widespread, err := decryptAndValidate(r.GetIv(), r.GetPayload(), user.Password, user.Mail, s.blackList)
	if err != nil {
		log.Printf("login: couldn't retrieve request data: %v\n", err)
		return nil, grpc.Errorf(codes.Unauthenticated, unauthenticated+"%q", mail)
	}

	var loginInfo message.Login
	if err := json.Unmarshal([]byte(widespread.Content), &loginInfo); err != nil {
		log.Printf("login: couldn't retrieve login message data: %v\n", err)
		return nil, grpc.Errorf(codes.Unauthenticated, unauthenticated+"%q", mail)
	}

	fmt.Printf("%T %#[1]v\n", loginInfo)

	return &LoginReply{}, nil
}
