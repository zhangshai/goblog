package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/app/models/category"
	"goblog/app/policies"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"goblog/pkg/view"
	"net/http"
)

type ArticlesController struct {
	BaseController
}

func (ac *ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	_article, err := article.Get(id)
	// 3. 如果出现错误
	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {

		view.Render(w, view.D{"Article": _article, "IsOwner": policies.CanModifyArticle(_article)}, "articles.show", "articles._article_meta")
	}

}

// Index 文章列表页
func (ac *ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	//获取结果集
	articles, pagerData, err := article.GetAll(r, 2)
	if err != nil {
		ac.ResponseForSQLError(w, err)

	} else {
		view.Render(w, view.D{"Articles": articles, "PagerData": pagerData}, "articles.index", "articles._article_meta")
	}

}

//创建文章
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {

	Category, err := category.GetByUserID(auth.UidToString())
	logger.LogError(err)
	data := view.D{
		"Title":    "",
		"Body":     "",
		"Article":  article.Article{},
		"Category": Category,
		"Errors":   make(map[string]string),
	}
	view.Render(w, data, "articles.create", "articles._form_field")
}

// Store 文章创建页面
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

	currentUser := auth.User()
	// 1. 初始化数据
	_article := article.Article{
		Title:      r.PostFormValue("title"),
		Body:       r.PostFormValue("body"),
		CategoryID: types.StringToUint64(r.PostFormValue("category_id")),
		UserID:     currentUser.ID,
	}

	errors := requests.ValidateArticleForm(_article)
	if len(errors) != 0 {
		Category, err := category.GetByUserID(auth.UidToString())
		logger.LogError(err)
		data := view.D{
			"Article":  _article,
			"Errors":   errors,
			"Category": Category,
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

func (ac *ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {

	id := route.GetRouteVariable("id", r)
	//读取对应文章
	_article, err := article.Get(id)

	if err != nil {

		ac.ResponseForSQLError(w, err)
	} else {
		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			Category, err := category.GetByUserID(auth.UidToString())
			logger.LogError(err)
			view.Render(w, view.D{
				"Article":  _article,
				"Category": Category,
				"Errors":   view.D{},
			}, "articles.edit", "articles._form_field")
		}
	}
}

func (ac *ArticlesController) Update(w http.ResponseWriter, r *http.Request) {

	id := route.GetRouteVariable("id", r)
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	_article, err := article.Get(id)

	if err != nil {
		ac.ResponseForSQLError(w, err)
	} else {

		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		}
		_article.Title = r.PostFormValue("title")
		_article.Body = r.PostFormValue("body")
		_article.CategoryID = types.StringToUint64(r.PostFormValue("category_id"))
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
			Category, err := category.GetByUserID(auth.UidToString())
			logger.LogError(err)
			view.Render(w, view.D{
				"Article":  _article,
				"Errors":   errors,
				"Category": Category,
			}, "articles.edit", "articles._form_field")
		}
	}
}

func (ac *ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {

	id := route.GetRouteVariable("id", r)
	_article, err := article.Get(id)
	// 检查权限
	if !policies.CanModifyArticle(_article) {
		ac.ResponseForUnauthorized(w, r)
	} else {
		if err != nil {

			ac.ResponseForSQLError(w, err)
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

}
