package route

import (
	"html/template"
	"log"
	"net/http"
)

//旧的注册模板方法，不建议用了
//func loginView(w http.ResponseWriter, r *http.Request) {
//	//解析
//	tpl,err:=template.ParseFiles("view/user/login.html");
//	if nil!=err {
//		log.Fatal(err.Error());//Fatal模板出错了打印并直接退出
//	}
//	tpl.ExecuteTemplate(w,"/user/login.shtml",nil);
//}

//访问后端文件shtml(后端模板文件渲染)
//http.HandleFunc("/user/login.shtml", loginView);

//注册所有的模板页面
func RegisterView() {
	//提供指定目录的静态文件支持
	http.Handle("/asset/",http.FileServer(http.Dir(".")));
	http.Handle("/mnt/",http.FileServer(http.Dir(".")))
	//解析
	tpl,err:=template.ParseGlob("view/**/*");
	if nil!=err {
		log.Fatal(err.Error());//Fatal模板出错了打印并直接退出
	}
	for _,v:= range tpl.Templates() {
		tplname := v.Name();
		//循环注册拿到的模板
		http.HandleFunc(tplname, func(writer http.ResponseWriter, request *http.Request) {
			tpl.ExecuteTemplate(writer,tplname,nil);
		})
	}

}
