package client

import (
	"fmt"
	"iris-app/app/lib/nacos"
	"iris-app/app/model/dto"
	"iris-app/app/model/request"
	"net/rpc"
	"net/rpc/jsonrpc"
)

const UserServiceName = "user"

type UserService interface {
	Get(userId int, user *dto.User) error
	Create(data request.RegisterUser, userId *int) error
}

type UserRpcClient struct {
	*rpc.Client
}

func NewUserClient() *UserRpcClient {
	ins := nacos.GetOneHealthyInstance(UserServiceName, "rpc-service")
	client, err := jsonrpc.Dial("tcp", fmt.Sprintf("%s:%d", ins.Ip, ins.Port))
	if err != nil {
		fmt.Println("create UserRpcClient err", err)
		return nil
	}
	return &UserRpcClient{client}
}

func (uc *UserRpcClient) Get(userId int, userRes *dto.User) error {
	return uc.Client.Call(UserServiceName+".Get", userId, userRes)
}

func (uc *UserRpcClient) Create(data request.RegisterUser, userId *int) error {
	return uc.Client.Call(UserServiceName+".Create", data, userId)
}
