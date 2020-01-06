package controllers

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

const (
	CMD_SINGLE_MSG = 10 //点对点单聊，dstid是用户id
	CMD_ROOM_MSG   = 11//群聊，dstid是群组id
	CMD_HEART      = 0//心跳消息，不处理
)

const (
	//文本样式
	MEDIA_TYPE_TEXT = 1
	//新闻样式类比图文消息
	MEDIA_TYPE_News = 2
	//语音样式
	MEDIA_TYPE_VOICE = 3
	//图片样式
	MEDIA_TYPE_IMG = 4
	//红包样式
	MEDIA_TYPE_REDPACKAGR = 5
	//emoj 表情样式
	MEDIA_TYPE_EMOJ = 6
	//emoj 表情样式
	MEDIA_TYPE_LINK = 7
	//emoj 表情样式
	MEDIA_TYPE_VIDEO = 8
	//emoj 表情样式
	MEDIA_TYPE_CONCAT = 9
	//emoj 表情样式
	MEDIA_TYPE_UDEF = 100
)

type Message struct {
	Id 		int64  `json:"id,omitempty" form:"id"`
	Userid 	int64  `json:"userid,omitempty" form:"userid"`
	Cmd 	int	   `json:"cmd,omitempty" form:"cmd"`
	Dstid 	int64  `json:"dstid,omitempty" form:"dstid"`
	Media	int	   `json:"media,omitempty" form:"media"`
	Content string `json:"content,omitempty" form:"content"`
	Pic		string `json:"pic,omitempty" form:"pic"`
	Url     string `json:"url,omitempty" form:"url"`
	Memo 	string `json:"memo,omitempty" form:"memo"`
	Amount	int	   `json:"amount,omitempty" form:"amount"`
}

/**
消息发送结构体
谁发的：userid，要发给谁：dstid，这个消息有什么用：cmd(单聊还是群聊),消息怎么展示：media，消息内容是什么：（url,pic,content..）
1、MEDIA_TYPE_TEXT
{id:1,userid:2,dstid:3,cmd:10,media:1,content:"hello"}
2、MEDIA_TYPE_News
{id:1,userid:2,dstid:3,cmd:10,media:2,content:"标题",pic:"http://www.baidu.com/a/log,jpg",url:"http://www.a,com/dsturl","memo":"这是描述"}
3、MEDIA_TYPE_VOICE，amount单位秒
{id:1,userid:2,dstid:3,cmd:10,media:3,url:"http://www.a,com/dsturl.mp3",anount:40}
4、MEDIA_TYPE_IMG
{id:1,userid:2,dstid:3,cmd:10,media:4,url:"http://www.baidu.com/a/log,jpg"}
5、MEDIA_TYPE_REDPACKAGR //红包amount 单位分
{id:1,userid:2,dstid:3,cmd:10,media:5,url:"http://www.baidu.com/a/b/c/redpackageaddress?id=100000","amount":300,"memo":"恭喜发财"}
6、MEDIA_TYPE_EMOJ 6
{id:1,userid:2,dstid:3,cmd:10,media:6,"content":"cry"}
7、MEDIA_TYPE_Link 6
{id:1,userid:2,dstid:3,cmd:10,media:7,"url":"http://www.a,com/dsturl.html"}

7、MEDIA_TYPE_Link 6
{id:1,userid:2,dstid:3,cmd:10,media:7,"url":"http://www.a,com/dsturl.html"}

8、MEDIA_TYPE_VIDEO 8
{id:1,userid:2,dstid:3,cmd:10,media:8,pic:"http://www.baidu.com/a/log,jpg",url:"http://www.a,com/a.mp4"}

9、MEDIA_TYPE_CONTACT 9
{id:1,userid:2,dstid:3,cmd:10,media:9,"content":"10086","pic":"http://www.baidu.com/a/avatar,jpg","memo":"胡大力"}

*/

//本科线在于形成userid河node的映射关系
type Node struct {
	Conn *websocket.Conn
	//并行转串行
	DataQueue chan []byte
	GroupSets set.Interface
}
//映射关系表
var clientMap map[int64] *Node = make(map[int64]*Node,0);
//读写锁
var rwlocker sync.RWMutex;




//ws://127.0.0.1/chat?id=1&token=xxx
func Chat(writer http.ResponseWriter,request *http.Request)  {
	//todo 检验是否合法
	query := request.URL.Query();
	id:= query.Get("id");
	token:=query.Get("token");
	userId, _ := strconv.ParseInt(id, 10, 64); //将字符串转为int64
	isvalida:=checkToken(userId,token);
	//isvalida==true 继续执行否则终止，用ws的Upgrader来处理
	conn,err:=(&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer,request,nil);
	if err != nil {
		log.Println(err.Error());
		return;
	}
	//todo 获得conn
	node :=&Node{
		Conn:conn,
		DataQueue:make(chan []byte,50),
		GroupSets:set.New(set.ThreadSafe),//当用户就得时候初始化groupset有用户加入的时候需要刷新
	}
	//获取用户全部群ID
	comIds:= contactService.SearchComunityIds(userId);
	for _,v:=range comIds {
		node.GroupSets.Add(v);
	}
	log.Println(node);
	//TODO todo userid和node形成绑定关系
	rwlocker.Lock();//读写锁执行锁操作
	clientMap[userId] = node;//赋值
	rwlocker.Unlock();//赋值完以后再解锁
	//todo 完成发送逻辑，com
	go sendproc(node);//调用发送协程
	//todo 完成接收逻辑
	go recvproc(node);

	sendMsg(userId,[]byte("test websoket!"));
	
}
//todo 添加新的群ID到用户的fourpset中
func AddGroupId(userId,gid int64)  {
	//取得node
	rwlocker.Lock();
	node,ok := clientMap[userId];
	if ok {
		node.GroupSets.Add(gid);//添加gid到set
	}
	rwlocker.Unlock();
}

//发送逻辑（发送协程）
func sendproc(node *Node){
	for{
		select {
		case data:= <-node.DataQueue://从管道取出数据
			err := node.Conn.WriteMessage(websocket.TextMessage, data);
			if err != nil {
				log.Println(err.Error());
				return
			}

		}
	}
}

//接收协程
func recvproc(node *Node)  {
	for{
		_,data,err := node.Conn.ReadMessage();
		if err != nil {
			log.Println(err.Error());
			return
		}
		//dispatch(data)//适合非分布式
		// 分布式部署用此方法	把消息广播到局域网
		broadMsg(data);
		log.Println("[ws]<=%s\n",data);
	}
}

func init(){
	go udpsendproc()
	go udprecvproc()
}

//用来存放发送的要广播的数据
var  udpsendchan chan []byte=make(chan []byte,1024)
//todo 将消息广播到局域网
func broadMsg(data []byte){
	udpsendchan<-data //从data通道拿消息放到chan里面去
}
//todo 完成udp数据的发送协程
func udpsendproc(){
	log.Println("start udpsendproc")
	//todo 使用udp协议拨号
	con,err:=net.DialUDP("udp",nil,
		&net.UDPAddr{
			IP:net.IPv4(192,168,6,255),
			Port:3000,
		})
	defer con.Close()
	if err!=nil{
		log.Println(err.Error())
		return
	}
	//todo 通过的到的con发送消息
	//con.Write()
	for{
		select {
		case data := <- udpsendchan:
			_,err=con.Write(data)
			if err!=nil{
				log.Println(err.Error())
				return
			}
		}
	}
}
//todo 完成upd接收并处理功能
func udprecvproc(){
	log.Println("start udprecvproc")
	//todo 监听udp广播端口
	con,err:=net.ListenUDP("udp",&net.UDPAddr{
		IP:net.IPv4zero,
		Port:3000,
	})
	defer con.Close()
	if err!=nil{log.Println(err.Error())}
	//TODO 处理端口发过来的数据
	for{
		var buf [512]byte
		n,err:=con.Read(buf[0:])
		if err!=nil{
			log.Println(err.Error())
			return
		}
		//直接数据处理
		dispatch(buf[0:n])
	}
	log.Println("stop updrecvproc")
}

//后端调度逻辑处理
func dispatch(data[]byte) {
	//todo 解析data为message
	log.Println("this is dispatch")
	msg := Message{};
	err := json.Unmarshal(data,&msg);
	if err != nil {
		log.Println(err.Error());
		return
	}
	//todo 根据cmd对逻辑进行处理
	switch msg.Cmd {
	case CMD_SINGLE_MSG:
		sendMsg(msg.Dstid,data);
	case CMD_ROOM_MSG:
		//todo 群聊分发逻辑
		for _, v := range clientMap{
			if v.GroupSets.Has(msg.Dstid) {
				v.DataQueue<-data;
			}
		}
	case CMD_HEART:
		//todo 心跳什么都不需要做
	}
}

//发送消息
func sendMsg(userId int64,msg []byte)  {
	rwlocker.RLock();//读写锁保证并发安全性
	node,ok := clientMap[userId];
	rwlocker.RUnlock();
	if ok{
		node.DataQueue <- msg;
	}
}
/**
检测是否有效
 */
func checkToken(userId int64,token string)bool{
	//从数据库里面查询并比对
	user := userService.Find(userId)
	return user.Token==token
}