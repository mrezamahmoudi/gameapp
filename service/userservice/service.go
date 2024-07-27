package userservice

import (
	"GameApp/entity"
	"GameApp/pkg/phonenumber"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error)
	GetUserByID(userID uint) (entity.User, error)
}
type Service struct {
	Repo Repository
}

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	User entity.User
}

func New(repo Repository) Service {
	return Service{repo}
}

func (s Service) Register(req *RegisterRequest) (RegisterResponse, error) {
	// todo - verify phonenumber by verification code

	// validate phonenumber
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("phoneNumber is not valid.")
	}

	// check uniqueness of phonenumber
	if isUnique, err := s.Repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return RegisterResponse{}, fmt.Errorf("Unexpected error %w", err)
		}
		if !isUnique {
			return RegisterResponse{}, fmt.Errorf("phoneNumber is not unique.")
		}
	}

	// validate name.
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("name length should be greater than 3.")
	}

	// TODO - check the password regex pattern
	// validate password.

	//bcrypt.GenerateFromPassword(pass, 0)
	user := entity.User{
		ID:          0,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Password:    GetMD5Hash(req.Password),
	}

	// create new user in storage
	createdUser, err := s.Repo.Register(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("Unexpected error %w", err)
	}
	// return created user
	return RegisterResponse{
		User: createdUser,
	}, nil
}

// Login Section
type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
}

func (s Service) Login(req *LoginRequest) (LoginResponse, error) {
	user, exist, err := s.Repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("Unexpected error %w", err)
	}
	if !exist {
		return LoginResponse{}, fmt.Errorf("username or password is not correct.")
	}

	if user.Password != GetMD5Hash(req.Password) {
		return LoginResponse{}, fmt.Errorf("username or password is not correct.")
	}

	// generate random session id
	// save session id in db
	
	// return session id to user

	return LoginResponse{}, nil
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// Get Profile section
type ProfileRequest struct {
	UserID uint
}
type ProfileResponse struct {
	Name string `json:"name"`
}

// all request inputs for interactor/services should be sanitized
func (s Service) Profile(req ProfileRequest) (res ProfileResponse, err error) {
	// get user by id
	user, err := s.Repo.GetUserByID(req.UserID)
	if err != nil {
		// I do not expect the repository call return "record not  found error" because i assume the interactor input is sanitized
		return ProfileResponse{}, fmt.Errorf("Unexpected error %w", err)
	}
	return ProfileResponse{user.Name}, nil
}

//func (s Service) GetUserByID(userID uint) (entity.User, error) {

//}
