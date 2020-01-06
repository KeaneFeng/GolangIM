package model;

import "time"

const (
	SEX_WOMEN="W"//女
	SEX_MAN="M"//南
	SEX_UNKNOW="U"//未知
);
//use modelName.SEX_WOMEN
type User struct {
	Id         	int64     `xorm:"pk autoincr bigint(64)" form:"id" json:"id"`//用户id
	Mobile   	string 		`xorm:"varchar(20)" form:"mobile" json:"mobile"`//用户手机号
   	Passwd    	string	`xorm:"varchar(40)" form:"passwd" json:"-"`   // 用户密码=f(plainpwd_salt),md5
	Avatar	   	string 		`xorm:"varchar(150)" form:"avatar" json:"avatar"`//头像
	Sex        	string	`xorm:"varchar(2)" form:"sex" json:"sex"`   // 性别
	Nickname    string	`xorm:"varchar(20)" form:"nickname" json:"nickname"`   // 别名
	Salt      	string	`xorm:"varchar(10)" form:"salt" json:"-"`   // 随机数
	Online     	int	`xorm:"int(10)" form:"online" json:"online"`   //是否在线
	Token      	string	`xorm:"varchar(40)" form:"token" json:"token"`   // token令牌 chat?id=1&token=...
	Memo      	string	`xorm:"varchar(140)" form:"memo" json:"memo"`   // 统计用户增加量
	Createat   	time.Time	`xorm:"datetime" form:"createat" json:"createat"`   // 创建时间
}