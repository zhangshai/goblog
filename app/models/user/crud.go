package user

import (
	"errors"
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/types"

	"gorm.io/gorm"
)

func (user *User) Create() (err error) {
	if err = model.DB.Create(&user).Error; err != nil {
		logger.LogError(err)
		return err

	}
	return nil
}

func Get(idstr string) (User, error) {

	var user User
	id := types.StringToUint64(idstr)
	if err := model.DB.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}

func GetByEmail(email string) (User, error) {
	var user User
	if err := model.DB.Where("email=?", email).First(&user).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return user, errors.New("账号不存在")
		} else {
			return user, errors.New("服务器内部错误")
		}
	}
	return user, nil
}
