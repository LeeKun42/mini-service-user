package middleware

import (
	"errors"
	"github.com/kataras/iris/v12"
	"strings"
	"user/app/service/jwt"
	"user/app/service/user"
)

func JwtAuthCheck(ctx iris.Context) {
	authorization := ctx.GetHeader("Authorization")
	if authorization == "" {
		ctx.StopWithError(401, errors.New("token无效"))
	}
	authArr := strings.Split(authorization, " ")
	if len(authArr) != 2 {
		ctx.StopWithError(401, errors.New("token无效"))
	}
	tokenString := authArr[1]
	claims, err := jwt.NewService().Check(tokenString)
	if err != nil {
		ctx.StopWithError(401, err)
	}
	ctx.Values().Set("user_id", claims.UserId)
	ctx.Values().Set("jwt_token", tokenString)
	ctx.Next()
}

func UserRoleCheck(ctx iris.Context) {
	userId, _ := ctx.Values().GetInt("user_id")
	userInfo, _ := user.NewService().Get(userId)
	if userInfo.ID != 0 {
		ctx.StopWithError(403, errors.New("无权限"))
	}
	ctx.Next()
}
