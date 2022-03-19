package category

import (
	"goblog/pkg/logger"
	"goblog/pkg/model"
	"goblog/pkg/types"
)

func (categoey *Category) Create() error {

	if err := model.DB.Create(&categoey).Error; err != nil {

		logger.LogError(err)
		return err
	}
	return nil
}

func All() ([]Category, error) {

	var categories []Category
	if err := model.DB.Find(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}
func Get(idstr string) (Category, error) {

	var category Category
	id := types.StringToUint64(idstr)

	if err := model.DB.First(&category, id).Error; err != nil {
		return category, err
	}
	return category, nil
}
func GetByUserID(uid string) ([]Category, error) {
	var categories []Category
	if err := model.DB.Where("user_id = ?", uid).Preload("User").Find(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}
