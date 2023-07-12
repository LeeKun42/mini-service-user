package service

import (
	"iris-app/app/model/dto"
	"iris-app/app/model/request"
	"iris-app/app/service/user"
)

type UserService interface {
	Get(userId int, user *dto.User) error
	Create(data request.RegisterUser, userId *int) error
}
type UserRpcService struct{}

func (us *UserRpcService) Get(userId int, reply *dto.User) error {
	user, err := user.Service.Get(userId)
	if err != nil {
		return err
	}
	*reply = user
	return nil
}

func (us *UserRpcService) Create(data request.RegisterUser, userId *int) error {
	id, err := user.Service.Register(data)
	if err != nil {
		return err
	}
	*userId = id
	return nil
}
