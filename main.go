package main

import (
	"database/sql"

	"goblog/bootstrap"
	"goblog/pkg/database"
	"net/http"

	"strings"

	"github.com/gorilla/mux"
)

//实例化db
var db *sql.DB
var router *mux.Router

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
	//添加路由中间件
	router.Use(forceHTMLMiddleware)
	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
