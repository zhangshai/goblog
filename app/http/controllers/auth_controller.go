package controllers

import (
	"goblog/app/models/user"
	"goblog/app/requests"
	"goblog/pkg/auth"
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
	_user := user.User{
		Name:            r.PostFormValue("name"),
		Email:           r.PostFormValue("email"),
		Password:        r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirm"),
	}
	errs := requests.ValidateRegistrationForm(_user)
	if len(errs) > 0 {
		view.RenderSimple(w, view.D{"Errors": errs, "User": _user}, "auth.register")

	} else {

		//2. 验证通过 —— 入库，并跳转到首页

		_user.Create()

		if _user.ID > 0 {

		} else {
			w.WriteHeader(http.StatusInternalServerError)

		}
	}

	// 3. 表单不通过 —— 重新显示表单
}

func (*AuthController) Login(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.login")

}

func (*AuthController) Dologin(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	if err := auth.Attempt(email, password); err == nil {

		http.Redirect(w, r, "/", http.StatusFound)
	} else {

		view.RenderSimple(w,
			view.D{
				"Error":    err.Error(),
				"Email":    email,
				"Password": password,
			}, "auth.login")
	}

}
