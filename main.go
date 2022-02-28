package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

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
	fmt.Fprint(w, "创建新的文章")

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
	route := mux.NewRouter()
	route.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	route.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")
	route.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	route.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	route.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("article.store")
	route.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	//添加路由中间件
	route.Use(forceHTMLMiddleware)
	homeURL, _ := route.Get("home").URL()
	fmt.Println("homeURL=", homeURL)

	articleURL, _ := route.Get("articles.show").URL("id", "23")
	fmt.Println(articleURL)
	http.ListenAndServe(":3000", removeTrailingSlash(route))
}
