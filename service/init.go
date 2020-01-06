package service

import (
	"../model"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

var DbEngin *xorm.Engine;
func init()  {
	drivename := "mysql";
	DsName := "root:123456@(127.0.0.1:3306)/chat?charset=utf8";//链接数据库
	err := errors.New("");
	DbEngin,err = xorm.NewEngine(drivename,DsName);
	if  (nil != err && err.Error()!="") {
		log.Fatal(err.Error());
	}
	DbEngin.ShowSQL(true);//是否显示sql语句
	DbEngin.SetMaxOpenConns(2);//设置数据库连接数
	DbEngin.Sync2(
	new(model.User),
		new(model.Community),
		new(model.Contact),
	);//自动建表
	fmt.Println("init database success");
}
