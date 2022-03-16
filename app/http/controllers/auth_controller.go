package controllers

import (
	"goblog/app/models/user"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/flash"
	"goblog/pkg/mail"
	"goblog/pkg/types"
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

			flash.Success("恭喜您注册成功！")
			// 登录用户并跳转到首页
			auth.Login(_user)
			http.Redirect(w, r, "/", http.StatusFound)
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

		flash.Success("欢迎回来！")
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

func (*AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	auth.Logout()
	flash.Success("您已退出登录")
	http.Redirect(w, r, "/", http.StatusFound)
}

func (*AuthController) FindPass(w http.ResponseWriter, r *http.Request) {

	view.RenderSimple(w, view.D{}, "auth.findpass")
}

func (*AuthController) DoFindPass(w http.ResponseWriter, r *http.Request) {

	email := r.PostFormValue("email")
	err := requests.ValidateMail(email)

	if len(err) > 0 {

		view.RenderSimple(w,
			view.D{
				"Error": err,
				"Email": email,
			}, "auth.findpass")
	} else {

		newpasswd := types.RandStr(6)

		subject := "找回密码邮件"
		body := "你的新密码为:" + newpasswd
		_errors := mail.SendMail(email, subject, body)

		e := []string{_errors.Error()}
		data := map[string][]string{
			"email": e,
		}

		if err != nil {
			view.RenderSimple(w,
				view.D{
					"Error": data,
					"Email": email,
				}, "auth.findpass")
		} else {
			flash.Success("邮件已发送")
			http.Redirect(w, r, "/", http.StatusFound)
		}

	}

}
