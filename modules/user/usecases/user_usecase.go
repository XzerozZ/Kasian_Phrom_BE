package usecases

import (
	"errors"
	"time"
	"mime/multipart"
	"os"
	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/user/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/gofiber/fiber/v2"
)

type UserUseCase interface {
	Register(user *entities.User, roleName string) (*entities.User, error)
	Login(email, password string) (string, *entities.User, error)
	UpdateUserByID(id string, user entities.User, files multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error)
}

type UserUseCaseImpl struct {
	userrepo 	repositories.UserRepository
	jwtSecret	string
	supa		configs.Supabase
}

func NewUserUseCase(userrepo repositories.UserRepository, jwt configs.JWT, supa configs.Supabase) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		userrepo: 	userrepo,
		jwtSecret:	jwt.Secret,
		supa:  		supa,
	}
}

func (u *UserUseCaseImpl) Register(user *entities.User, roleName string) (*entities.User, error) {
	role, err := u.userrepo.GetRoleByName(roleName)
	if err != nil {
		return nil, errors.New("role not found")
	}

	user.ID = uuid.New().String()
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
		"user_id": user.ID,
		"role":    user.Role.RoleName,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", nil, err
	}
	
	return tokenString, &user, nil
}

func (u *UserUseCaseImpl) UpdateUserByID(id string, user entities.User, file multipart.FileHeader, ctx *fiber.Ctx) (*entities.User, error) {
	existingUser, err := u.userrepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	existingUser.Firstname = user.Firstname
	existingUser.Lastname = user.Lastname
	existingUser.Username = user.Username
	existingUser.Email = user.Email
	fileName := uuid.New().String() + ".jpg"
	if err := ctx.SaveFile(&file, "./uploads/"+fileName); err != nil {
		return nil, err
	}

	imageUrl, err := utils.UploadImage(fileName, "", u.supa)
	if err != nil {
		os.Remove("./uploads/" + fileName)
		return nil, err
	}

	if err := os.Remove("./uploads/" + fileName); err != nil {
		return nil, err
	}

	existingUser.ImageLink = imageUrl
	var updatedUser *entities.User
	updatedUser, err = u.userrepo.UpdateUserByID(existingUser)
    if err != nil {
        return nil, err
    }

	return updatedUser, nil
}