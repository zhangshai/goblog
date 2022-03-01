package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

//定义路由
var route = mux.NewRouter()

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
func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "文章ID"+id)
}
func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "访问文章列表")
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	if error := r.ParseForm(); error != nil {
		fmt.Fprint(w, "请提供正确的数据")
		return
	}
	title := r.PostForm.Get("title")
	fmt.Fprintf(w, "POST PostForm: %v <br>", r.PostForm)
	fmt.Fprintf(w, "POST Form: %v <br>", r.Form)
	fmt.Fprintf(w, "title 的值为: %v", title)

}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<title>创建文章 —— 我的技术博客</title>
	</head>
	<body>
		<form action="%s" method="post">
			<p><input type="text" name="title"></p>
			<p><textarea name="body" cols="30" rows="10"></textarea></p>
			<p><button type="submit">提交</button></p>
		</form>
	</body>
	</html>
	`
	storeURL, _ := route.Get("articles.store").URL()
	fmt.Fprintf(w, html, storeURL)
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
	route.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	route.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")
	route.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	route.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	route.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("article.store")
	route.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("article.create")
	route.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	//添加路由中间件
	route.Use(forceHTMLMiddleware)
	http.ListenAndServe(":3000", removeTrailingSlash(route))
}
