package main

import (
	"./route"
	"net/http"
)

/**
程序主入口
 */
func main() {
	route.ApiRoute();//加载api路由
	route.RegisterView();//执行注册模板方法
	//启动项目服务器
	http.ListenAndServe(":8800",nil);
	
}
