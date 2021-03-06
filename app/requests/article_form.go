package requests

import (
	"goblog/app/models/article"

	"github.com/thedevsaddam/govalidator"
)

func ValidateArticleForm(data article.Article) map[string][]string {

	rules := govalidator.MapData{
		"title":      []string{"required", "min_cn:3", "max_cn:40"},
		"body":       []string{"required", "min_cn:10"},
		"categoryid": []string{"required"},
	}

	// 2. 定制错误消息
	messages := govalidator.MapData{
		"title": []string{
			"required:标题为必填项",
			"min:标题长度需大于 3",
			"max:标题长度需小于 40",
		},
		"body": []string{
			"required:文章内容为必填项",
			"min:长度需大于 10",
		},
		"categoryid": []string{
			"required:请选择分类",
			"gt: 请选择分类",
		},
	}

	// 3. 配置初始化
	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "valid", // 模型中的 Struct 标签标识符
		Messages:      messages,
	}

	// 4. 开始验证
	return govalidator.New(opts).ValidateStruct()
}
