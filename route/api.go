package route

import (
	"../controllers"
	"net/http"
)

func ApiRoute(){
	//用户登录接口
	http.HandleFunc("/user/login", controllers.UserLogin);
	//用户注册接口
	http.HandleFunc("/user/register",controllers.UserRegister);
	//显示全部好友，参数userid
	http.HandleFunc("/contact/loadfriend", controllers.LoadFriend);
	//建群，头像pic,名称name，备注memo
	http.HandleFunc("/contact/createcommunity", controllers.CreateCommunity);
	//显示全部群，参数userid
	http.HandleFunc("/contact/loadcommunity", controllers.LoadCommunity);
	//加群，参数uerid，dstid
	http.HandleFunc("/contact/joincommunity", controllers.JoinCommunity);
	//自动添加好友
	http.HandleFunc("/contact/addfriend", controllers.Addfriend);
	//ws路由
	http.HandleFunc("/chat", controllers.Chat);
	//文件上傳
	http.HandleFunc("/attach/upload", controllers.Upload);
}