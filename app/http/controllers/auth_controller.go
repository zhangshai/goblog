package controllers

import (
	"goblog/app/models/user"
	"goblog/pkg/view"
	"net/http"
)

type AuthController struct {
}

func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")

}
func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {

	// 1. 表单验证
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	// 2. 验证通过 —— 入库，并跳转到首页

	_user := user.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	_user.Create()

	if _user.ID > 0 {

	} else {
		w.WriteHeader(http.StatusInternalServerError)

	}

	// 3. 表单不通过 —— 重新显示表单
}
