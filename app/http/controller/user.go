package controller

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"time"
	"user/app/http/response"
	"user/app/model/dto"
	"user/app/model/request"
	response2 "user/app/model/resp"
	"user/app/rpc/client"
	"user/app/service/jwt"
	"user/app/service/user"
)

type UserController struct {
	UserService *user.Service
	JwtService  *jwt.Service
}

func NewUserController() *UserController {
	return &UserController{
		UserService: user.NewService(),
		JwtService:  jwt.NewService(),
	}
}

func (uc *UserController) Hello(context iris.Context) {
	/** rpc 调用测试*/
	var user dto.User
	err := client.NewUserClient().Get(1, &user)
	if err != nil {
		fmt.Println("rpc call error", err)
	}
	response.Success(context, iris.Map{"message": "Hello Iris!v1.0.1", "time": time.Now().Format("2006-01-02 13:01:02"), "user info": user})
}

func (uc *UserController) Register(context iris.Context) {
	var input request.RegisterUser
	context.ReadJSON(&input)
	id, err := uc.UserService.Register(input)
	if err != nil {
		response.Fail(context, err.Error())
	} else {
		response.Success(context, iris.Map{"user_id": id})
	}
}

func (uc *UserController) Login(context iris.Context) {
	var input request.UserLogin
	context.ReadQuery(&input)
	token, err := uc.UserService.Login(input)
	if err != nil {
		response.Fail(context, err.Error())
	} else {
		response.Success(context, iris.Map{"token": token})
	}
}

func (uc *UserController) RefreshToken(context iris.Context) {
	token, err := uc.JwtService.Refresh(context.Values().GetString("jwt_token"))
	if err != nil {
		response.Error(context, 401, err.Error())
	} else {
		response.Success(context, iris.Map{"token": token})
	}
}

func (uc *UserController) Logout(context iris.Context) {
	uc.JwtService.Invalidate(context.Values().GetString("jwt_token"))
	response.Success(context, iris.Map{})
}

func (uc *UserController) Info(context iris.Context) {
	loginUserId, _ := context.Values().GetInt("user_id")
	userInfo, err := uc.UserService.Get(loginUserId)
	if err != nil {
		response.Fail(context, err.Error())
	} else {
		result := response2.UserInfo{
			ID:       userInfo.ID,
			NickName: userInfo.NickName,
			Account:  userInfo.Account,
			Status:   userInfo.Status,
		}

		response.Success(context, result)
	}
}

func (uc *UserController) GetList(context iris.Context) {
	var input request.UserSearch
	context.ReadQuery(&input)
	if input.PageIndex == 0 {
		input.PageIndex = 1
	}
	if input.PageSize == 0 {
		input.PageSize = 10
	}
	users, err := uc.UserService.GetList(input)
	if err != nil {
		response.Fail(context, err.Error())
	} else {
		var res []response2.UserInfo
		for _, userInfo := range users {
			result := response2.UserInfo{
				ID:       userInfo.ID,
				NickName: userInfo.NickName,
				Account:  userInfo.Account,
				Status:   userInfo.Status,
			}
			res = append(res, result)
		}

		response.Success(context, res)
	}
}
