package service

import (
	"user/app/model/dto"
	"user/app/service/jwt"
	"user/app/service/user"
)

type UserService interface {
	Get(userId int, user *dto.User) error
	CheckJwtToken(token string, userId *int) error
}
type UserRpcService struct{}

func (us *UserRpcService) Get(userId int, reply *dto.User) error {
	user, err := user.NewService().Get(userId)
	if err != nil {
		return err
	}
	*reply = user
	return nil
}

func (us *UserRpcService) CheckJwtToken(token string, userId *int) error {
	claims, err := jwt.NewService().Check(token)
	if err != nil {
		return err
	}
	*userId = claims.UserId
	return nil
}
