package main

import (
	"database/sql"
	"fmt"
	"goblog/bootstrap"
	"goblog/pkg/database"
	"goblog/pkg/logger"
	"html/template"

	"net/http"
	"net/url"
	"strconv"
	"strings"

	"unicode/utf8"

	"github.com/gorilla/mux"
)

//实例化db
var db *sql.DB
var router *mux.Router

type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "请求页面不存在")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "<h1>Hello, 这里是 goblog</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

type Article struct {
	Title, Body string
	ID          int64
}

func (a Article) Delete() (rowsAffected int64, err error) {
	rs, err := db.Exec("delete from articles where id=" + strconv.FormatInt(a.ID, 10))
	if err != nil {
		return 0, err
	}
	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}
	return 0, nil
}

func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	//读取对应文章
	article := Article{}
	query := "select * from articles where id =?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "文章不存在")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "服务器内部错误")
		}
	} else {
		updateURL, _ := router.Get("articles.update").URL("id", id)
		data := ArticlesFormData{
			Title:  article.Title,
			Body:   article.Body,
			URL:    updateURL,
			Errors: nil,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
		logger.LogError(err)
		err = tmpl.Execute(w, data)
		logger.LogError(err)
	}
}
func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	_, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")

		}
	} else {
		errors := validateArticleFormData(title, body)
		if len(errors) == 0 {
			query := "update articles set title=?,body=? where id=?"
			stmt, err := db.Prepare(query)
			defer stmt.Close()
			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			}
			rs, err := stmt.Exec(title, body, id)
			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			}
			if n, _ := rs.RowsAffected(); n > 0 {
				showURL, _ := router.Get("articles.show").URL("id", id)
				http.Redirect(w, r, showURL.String(), http.StatusFound)
			} else {
				fmt.Fprint(w, "您没有做任何更改！")
			}
		} else {
			updateURL, _ := router.Get("articles.update").URL("id", id)
			data := ArticlesFormData{
				Title:  title,
				Body:   body,
				URL:    updateURL,
				Errors: errors,
			}
			tmpl, err := template.ParseFiles("resources/views/articles/edit.gohtml")
			logger.LogError(err)
			err = tmpl.Execute(w, data)
			logger.LogError(err)
		}
	}

}
func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

func articlesDelHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)
	rows, err := getArticleByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	}
	rs, err := rows.Delete()
	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 服务器内部错误")
	} else {
		if rs > 0 {
			indexURL, _ := router.Get("articles.index").URL()
			http.Redirect(w, r, indexURL.String(), http.StatusFound)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		}
	}

}

//表单验证
func validateArticleFormData(title, body string) map[string]string {
	errors := make(map[string]string)
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}

	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}
	return errors
}

// 获取某一篇文章
func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "select * from articles where id =?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err

}

//定义路由中间件函数
func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(rw, r)
	})
}

//修正错误路由
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		next.ServeHTTP(rw, r)
	})
}

func main() {
	database.Initialize()
	db = database.DB
	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()

	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDelHandler).Methods("POST").Name("articles.del")

	//添加路由中间件
	router.Use(forceHTMLMiddleware)
	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
