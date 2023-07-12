package user

import (
	"errors"
	"gorm.io/gorm"
	"iris-app/app/lib/hash"
	"iris-app/app/model"
	"iris-app/app/model/dto"
	"iris-app/app/model/request"
	"iris-app/app/service/jwt"
	"time"
)

type userService struct{}

var Service = userService{}

const (
	// StatusDisabled 用户状态 禁用
	StatusDisabled int = 0
	// StatusEnable 用户状态 启用
	StatusEnable int = 1
)

func (us *userService) Register(input request.RegisterUser) (int, error) {
	var user dto.User
	model.Instance().Table(user.TableName()).Where("account", input.Account).First(&user)
	if user.ID != 0 {
		return 0, errors.New("注册失败：账号已存在")
	}
	user = dto.User{
		Account:   input.Account,
		NickName:  input.NickName,
		Passwd:    hash.Make(input.Passwd),
		CreatedAt: time.Now(),
		Status:    StatusEnable,
	}
	result := model.Instance().Create(&user)
	if result.Error != nil {
		return 0, errors.New("注册失败：账号已存在")
	}
	return user.ID, nil
}

func (us *userService) Login(input request.UserLogin) (string, error) {
	var user dto.User
	result := model.Instance().Table(user.TableName()).Where("account", input.Account).First(&user)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", errors.New("账号不存在")
	}
	if user.Status == StatusDisabled {
		return "", errors.New("账号已被禁用，请联系管理员")
	}
	if !hash.Check(input.Passwd, user.Passwd) { //密码不正确
		return "", errors.New("密码不正确")
	}
	token := jwt.Service.Create(user.ID, 0)
	return token, nil
}

func (us *userService) Disabled(userId int) error {
	var user dto.User
	tx := model.Instance().Table(user.TableName()).Where("id", userId).Update("status", StatusDisabled)
	return tx.Error
}

func (us *userService) Enable(userId int) error {
	var user dto.User
	tx := model.Instance().Table(user.TableName()).Where("id", userId).Update("status", StatusEnable)
	return tx.Error
}

func (us *userService) GetList(input request.UserSearch) ([]dto.User, error) {
	var users []dto.User
	offset := (input.PageIndex - 1) * input.PageSize
	result := model.Instance().Offset(offset).Limit(input.PageSize).Preload("Service").Find(&users)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return users, errors.New("用户不存在")
	}
	return users, nil
}

func (us *userService) Get(id int) (dto.User, error) {
	var user dto.User
	result := model.Instance().Preload("Service").Find(&user, id)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user, errors.New("用户不存在")
	}
	return user, nil
}
