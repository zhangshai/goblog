package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, 这里是 goblog</h1>")
	} else if r.URL.Path == "/about" {
		fmt.Fprint(w, "<h1>此博客是用以记录编程笔记，如您有反馈或建议</h1><a href='https://www.baidu.com'>baidu</a>")
	} else {
		fmt.Fprint(w, "<h1>无法找到页面</h1>")
	}

}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)
}
