package controllers

import (
	"../util"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func init()  {
	os.MkdirAll("./mnt",os.ModePerm);//允许文件自动创建
}
func Upload(w http.ResponseWriter,r *http.Request)  {
	UploadLocal(w,r);
}

func UploadLocal(writer http.ResponseWriter,request *http.Request)  {
	//todo 获得上传的源文件
	srcfile,head,err := request.FormFile("file");
	if err != nil {
		util.RespFail(writer,err.Error());
	}

	//todo 创建一个新文件
	suffix := ".png";//默认的文件后缀
	//如果前端文件名称包含后缀
	ofilename :=head.Filename;
	tmp := strings.Split(ofilename,".");//.号分割数组
	if len(tmp)>1 {
		suffix = "."+tmp[len(tmp)-1];//最后一个就是后缀名
	}
	//如果前端指定filetype
	//formdata.append(filetype,".png");
	filetype := request.FormValue("filetype");
	if len(filetype)>0 {
		suffix=filetype;//如果指定就使用前端指定的filetype
	}
	filename := fmt.Sprintf("%d%04d%s",time.Now().Unix(),rand.Int31(),suffix);//定义文件名称
	fmt.Println(filename);
	dstfile,err := os.Create("./mnt/"+filename);//创建新文件
	if err != nil {
		util.RespFail(writer,err.Error());
		return
	}
	//todo 将源文件内容copy到新文件
	_, err =io.Copy(dstfile,srcfile);
	if err != nil {
		util.RespFail(writer,err.Error());
		return
	}
	//todo 将新文件路径转成url地址
	url := "/mnt/"+filename;
	//todo 响应到前端
	util.RespOk(writer,url,"");
}
/**
oss 配置信息
 */
const (
	AccessKeyId="5p2RZKnrUanMuQw9"
	AccessKeySecret="bsNmjU8Au08axedV40TRPCS5XIFAkK"
	EndPoint="oss-cn-shenzhen.aliyuncs.com"
	Bucket="winliondev"
)
//权限设置为公共读状态
//需要安装
func UploadOss(writer http.ResponseWriter,request *http.Request)  {
	//todo 获得上传的文件
	srcfile,head,err:=request.FormFile("file")
	if err!=nil{
		util.RespFail(writer,err.Error())
		return
	}


	//todo 获得文件后缀.png/.mp3

	suffix := ".png"
	//如果前端文件名称包含后缀 xx.xx.png
	ofilename := head.Filename
	tmp := strings.Split(ofilename,".")
	if len(tmp)>1{
		suffix = "."+tmp[len(tmp)-1]
	}
	//如果前端指定filetype
	//formdata.append("filetype",".png")
	filetype := request.FormValue("filetype")
	if len(filetype)>0{
		suffix = filetype
	}
	//todo 初始化ossclient
	client,err:=oss.New(EndPoint,AccessKeyId,AccessKeySecret)
	if err!=nil{
		util.RespFail(writer,err.Error())
		return
	}
	//todo 获得bucket
	bucket,err := client.Bucket(Bucket)
	if err!=nil{
		util.RespFail(writer,err.Error())
		return
	}
	//todo 设置文件名称
	//time.Now().Unix()
	filename := fmt.Sprintf("mnt/%d%04d%s",
		time.Now().Unix(), rand.Int31(),
		suffix)
	//todo 通过bucket上传
	err=bucket.PutObject(filename,srcfile)
	if err!=nil{
		util.RespFail(writer,err.Error())
		return
	}
	//todo 获得url地址
	url := "http://"+Bucket+"."+EndPoint+"/"+filename

	//todo 响应到前端
	util.RespOk(writer,url,"")
}