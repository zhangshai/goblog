package route

import "github.com/gorilla/mux"

// Router 路由对象
var Route *mux.Router

// Initialize 初始化路由
func Initialize() {
	Route = mux.NewRouter()
}

// RouteName2URL 通过路由名称来获取 URL
func Name2URL(routeName string, pairs ...string) string {
	url, err := Route.Get(routeName).URL(pairs...)
	if err != nil {
		// checkError(err)
		return ""
	}

	return url.String()
}
