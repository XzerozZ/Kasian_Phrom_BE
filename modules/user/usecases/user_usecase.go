package usecases

import (
	"errors"
	"time"
	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/user/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
)

type UserUseCase interface {
	Register(user *entities.User, roleName string) (*entities.User, error)
	Login(email, password string) (string, *entities.User, error)
}

type UserUseCaseImpl struct {
	userrepo 	repositories.UserRepository
	jwtSecret	string
}

func NewUserUseCase(userrepo repositories.UserRepository, config configs.JWT) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		userrepo: 	userrepo,
		jwtSecret:	config.Secret,
	}
}

func (u *UserUseCaseImpl) Register(user *entities.User, roleName string) (*entities.User, error) {
	role, err := u.userrepo.GetRoleByName(roleName)
	if err != nil {
		return nil, errors.New("role not found")
	}

	user.ID = uuid.New()
	user.RoleID = role.ID
	user.Role = role
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)
	return u.userrepo.CreateUser(user)
}

func (u *UserUseCaseImpl) Login(email, password string) (string, *entities.User, error) {
	user, err := u.userrepo.FindUserByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    user.Role.RoleName,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", nil, err
	}
	return tokenString, &user, nil
}