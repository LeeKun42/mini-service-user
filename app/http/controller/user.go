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

type userController struct{}

var UserController = userController{}

func (uc *userController) Hello(context iris.Context) {
	/** rpc 调用测试*/
	var user dto.User
	err := client.NewUserClient().Get(1, &user)
	if err != nil {
		fmt.Println("rpc call error", err)
	}
	response.Success(context, iris.Map{"message": "Hello Iris!v1.0.1", "time": time.Now().Format("2006-01-02 13:01:02"), "user info": user})
}

func (uc *userController) Register(context iris.Context) {
	var input request.RegisterUser
	context.ReadJSON(&input)
	id, err := user.Service.Register(input)
	if err != nil {
		response.Fail(context, err.Error())
	} else {
		response.Success(context, iris.Map{"user_id": id})
	}
}

func (uc *userController) Login(context iris.Context) {
	var input request.UserLogin
	context.ReadQuery(&input)
	token, err := user.Service.Login(input)
	if err != nil {
		response.Fail(context, err.Error())
	} else {
		response.Success(context, iris.Map{"token": token})
	}
}

func (uc *userController) RefreshToken(context iris.Context) {
	token, err := jwt.Service.Refresh(context.Values().GetString("jwt_token"))
	if err != nil {
		response.Error(context, 401, err.Error())
	} else {
		response.Success(context, iris.Map{"token": token})
	}
}

func (uc *userController) Logout(context iris.Context) {
	jwt.Service.Invalidate(context.Values().GetString("jwt_token"))
	response.Success(context, iris.Map{})
}

func (uc *userController) Info(context iris.Context) {
	loginUserId, _ := context.Values().GetInt("user_id")
	userInfo, err := user.Service.Get(loginUserId)
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

func (uc *userController) GetList(context iris.Context) {
	var input request.UserSearch
	context.ReadQuery(&input)
	if input.PageIndex == 0 {
		input.PageIndex = 1
	}
	if input.PageSize == 0 {
		input.PageSize = 10
	}
	users, err := user.Service.GetList(input)
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
