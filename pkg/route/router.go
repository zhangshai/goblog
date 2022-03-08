package route

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var Route *mux.Router

func SetRoute(r *mux.Router) {
	Route = r
}

// RouteName2URL 通过路由名称来获取 URL
func Name2URL(routeName string, pairs ...string) string {
	// Router 路由对象
	fmt.Println(pairs)
	url, err := Route.Get(routeName).URL(pairs...)

	if err != nil {
		//logger.LogError(err)
		return ""
	}

	return url.String()
}

// 获取路由参数
func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
