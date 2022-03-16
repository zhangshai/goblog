package middlewares

import (
	"goblog/pkg/auth"
	"goblog/pkg/flash"
	"net/http"
)

func Auth(next HttpHandlerFunc) HttpHandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {
		if !auth.Check() {
			flash.Warning("登录用户才能访问此页面")
			http.Redirect(rw, r, "/", http.StatusFound)
			return
		}
		next(rw, r)
	}
}
