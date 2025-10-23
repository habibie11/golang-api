package user

import "golang.org/x/crypto/bcrypt"

// type RegisterUserInput struct {
// 	Name       string
// 	Occupation string
// 	Email      string
// 	Password   string
// }

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Occupation = input.Occupation
	user.Email = input.Email

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	user.PasswordHash = string(passwordHash)

	user.Role = "user"

	newUser, err := s.repository.Save(user)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}