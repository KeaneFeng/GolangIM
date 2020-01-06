package service;

import (
	"../model"
	"errors"
	"fmt"
	"math/rand"
	"../util"
	"time"
)

type UserService struct {

}

func (s *UserService)Register(
	mobile,//手机号
	plainpwd,//明文密码
	nickname,//昵称
	avatar,//头像地址
	sex string,//性别
)(user model.User,err error){
	//检测手机号码是否存在
	tmp:=model.User{};
	_, err = DbEngin.Where("mobile=?", mobile).Get(&tmp);
	if nil!=err {
		return tmp,err
	}
	//如果存在则返回已经注册
	if tmp.Id>0 {
		return tmp,errors.New("该手机号码已经注册");
	}
	//否则拼接插入数据
	tmp.Mobile = mobile;
	tmp.Avatar = avatar;
	tmp.Nickname = nickname;
	tmp.Sex = sex;
	tmp.Salt = fmt.Sprintf("%06d",rand.Int31n(10000));
	tmp.Passwd = util.MakePasswd(plainpwd,tmp.Salt);
	tmp.Createat = time.Now();
	//token 可以是一个随机书
	tmp.Token = fmt.Sprintf("%08d",rand.Int31());
	_, err = DbEngin.InsertOne(&tmp);//1前端恶意插入特殊字符，2数据库链接操作失败
	//返回新用户信息
	return tmp,err;
}

func (s *UserService)Login(
	mobile,//手机号
	plainpwd string,//明文密码
)(user model.User,err error){
	//通过手机号码查询用户
	tmp:=model.User{};
	DbEngin.Where("mobile = ?",mobile).Get(&tmp);
	//根据查询到的数据校验密码
	if tmp.Id == 0 {
		return tmp,errors.New("用户不存在");
	}
	//比对密码
	if !util.ValidatePasswd(plainpwd,tmp.Salt,tmp.Passwd){
		return tmp,errors.New("密码不正确");
	}
	//刷新token,安全机制
	str := fmt.Sprintf("%d",time.Now().Unix());
	token := util.MD5Encode(str);
	tmp.Token=token;
	DbEngin.ID(tmp.Id).Cols("token").Update(&tmp);
	//返回数据
	return tmp,nil;
}
//查询某个用户
func (s *UserService)Find(userId int64)(user model.User){
	tmp :=model.User{};
	DbEngin.ID(userId).Get(&tmp);
	return tmp;
}