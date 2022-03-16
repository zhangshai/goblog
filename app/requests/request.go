package requests

import (
	"errors"
	"fmt"
	"goblog/pkg/model"
	"strings"

	"github.com/thedevsaddam/govalidator"
)

func init() {
	// not_exists:users,email
	govalidator.AddCustomRule("not_exists", func(field, rule, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "not_exists:"), ",")

		tableName := rng[0]
		dbFiled := rng[1]
		val := value.(string)
		var count int64
		model.DB.Table(tableName).Where(dbFiled+" =?", val).Count(&count)
		if count != 0 {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("%v 已被占用", val)
		}
		return nil

	})
	govalidator.AddCustomRule("must_exists", func(field, rule, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "must_exists:"), ",")

		tableName := rng[0]
		dbFiled := rng[1]
		val := value.(string)
		var count int64
		model.DB.Table(tableName).Where(dbFiled+" =?", val).Count(&count)
		if count == 0 {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("%v 不存在", val)
		}
		return nil

	})
}