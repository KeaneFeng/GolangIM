package controllers

import (
	"../model"
	"../service"
	"../util"
	"net/http"
)

//用户注册方法
var userService service.UserService;
func UserRegister(writer http.ResponseWriter, request *http.Request) {
	//数据库操作
	//逻辑处理
	//restapi 返回json
	//获取当前前端参数 mobile passwd
	request.ParseForm();
	mobile := request.PostForm.Get("mobile");
	plainpwd := request.PostForm.Get("passwd");
	nickname := request.PostForm.Get("nickname");
	avatar := request.PostForm.Get("avatar");
	sex := model.SEX_UNKNOW;
	user, err := userService.Register(mobile, plainpwd, nickname, avatar, sex);
	//判断返回信息
	if nil!=err{
		util.RespFail(writer,err.Error());
	}else{
		util.RespOk(writer,user,"");
	}
}
//用户登录方法
func UserLogin(writer http.ResponseWriter, request *http.Request) {
	//数据库操作
	//逻辑处理
	//restapi 返回json
	//获取当前前端参数 mobile passwd
	request.ParseForm();
	mobile := request.PostForm.Get("mobile");
	passwd := request.PostForm.Get("passwd");

	user, err := userService.Login(mobile, passwd);
	if nil!=err {
		util.RespFail(writer,err.Error());
	}else{
		util.RespOk(writer,user,"")
	}
}
