package controllers

import (
	"../args"
	"../model"
	"../service"
	"../util"
	"net/http"
)

var contactService  service.ContactService;

//加载好友列表
func LoadFriend (w http.ResponseWriter,req *http.Request)  {
	var arg args.ContactArg;
	util.Bind(req,&arg);
	users := contactService.SearchFriend(arg.Userid);
	util.RespOkList(w,users,len(users))
}
//获取我的所有群组
func LoadCommunity(w http.ResponseWriter,req *http.Request)  {
	var arg args.ContactArg;
	util.Bind(req,&arg);
	comunitys := contactService.SearchComunity(arg.Userid);
	util.RespOkList(w,comunitys,len(comunitys));
}
//加入群组
func JoinCommunity(w http.ResponseWriter,req *http.Request)  {
	var arg args.ContactArg;
	util.Bind(req,&arg);
	err := contactService.JoinCommunity(arg.Userid, arg.Dstid);
	//todo 刷新用户的群组信息
	AddGroupId(arg.Userid,arg.Dstid);
	if err != nil {
		util.RespFail(w,err.Error());
	}
	util.RespOk(w,nil,"");
}
//创建群组
func CreateCommunity(w http.ResponseWriter,req *http.Request)  {
	var arg model.Community;
	util.Bind(req,&arg);
	com, err := contactService.CreateCommunity(arg);
	if err != nil {
		util.RespFail(w,err.Error());
	}else {
		util.RespOk(w,com,"");
	}

}
//自动添加好友（无需通过）
func Addfriend(w http.ResponseWriter,req *http.Request)  {
	var arg args.ContactArg;
	//参数的对象绑定
	util.Bind(req,&arg);
	err := contactService.AddFriend(arg.Userid,arg.Dstid);

	if err!=nil{
		util.RespFail(w,err.Error());
	}else{
		util.RespOk(w,nil,"好友添加成功");
	}
}