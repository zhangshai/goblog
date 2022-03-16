package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/app/requests"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"

	"gorm.io/gorm"
)

type ArticlesController struct {
}

func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	_article, err := article.Get(id)
	// 3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {

		view.Render(w, view.D{"Article": _article}, "articles.show", "articles._article_meta")
	}

}

// Index 文章列表页
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	//获取结果集
	articles, err := article.GetAll()
	if err != nil {
		logger.LogError(err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		view.Render(w, view.D{"Articles": articles}, "articles.index", "articles._article_meta")
	}

}

//创建文章
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	data := view.D{
		"Title":   "",
		"Body":    "",
		"Article": article.Article{},
		"Errors":  make(map[string]string),
	}
	view.Render(w, data, "articles.create", "articles._form_field")
}

// Store 文章创建页面
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {
	// 1. 初始化数据
	_article := article.Article{
		Title: r.PostFormValue("title"),
		Body:  r.PostFormValue("body"),
	}

	errors := requests.ValidateArticleForm(_article)
	if len(errors) != 0 {

		data := view.D{
			"Article": _article,
			"Errors":  errors,
		}
		view.Render(w, data, "articles.create", "articles._form_field")

	} else {

		_article.Create()
		if _article.ID > 0 {
			indexURL := route.Name2URL("articles.show", "id", _article.GetStringID())
			http.Redirect(w, r, indexURL, http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	}
}

func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {

	id := route.GetRouteVariable("id", r)
	//读取对应文章
	_article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章不存在")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {

		view.Render(w, view.D{
			"Article": _article,
			"Errors":  view.D{},
		}, "articles.edit", "articles._form_field")
	}
}

func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {

	id := route.GetRouteVariable("id", r)
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	_article, err := article.Get(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")

		}
	} else {
		_article.Title = r.PostFormValue("title")
		_article.Body = r.PostFormValue("body")

		errors := requests.ValidateArticleForm(_article)
		if len(errors) == 0 {
			_article.Title = title
			_article.Body = body
			rowsAffected, err := _article.Update()
			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			}
			if rowsAffected > 0 {
				showURL := route.Name2URL("articles.show", "id", id)
				http.Redirect(w, r, showURL, http.StatusFound)
			} else {
				fmt.Fprint(w, "您没有做任何更改！")
			}
		} else {
			view.Render(w, view.D{
				"Article": _article,
				"Errors":  errors,
			}, "articles.edit", "articles._form_field")
		}
	}
}

func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {

	id := route.GetRouteVariable("id", r)
	_article, err := article.Get(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	}
	rowsAffected, err := _article.Delete()
	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		if rowsAffected > 0 {
			indexURL := route.Name2URL("articles.index")
			http.Redirect(w, r, indexURL, http.StatusFound)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		}
	}

}
