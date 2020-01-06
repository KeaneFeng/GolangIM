package service

import (
	"errors"
	"../model"
	"time"
)

type ContactService struct {

}
//自动添加好友
func (service *ContactService)AddFriend(
	userid,
	dstid int64) error{
	//如果加自己
	if userid == dstid {
		return errors.New("不能添加自己");
	}
	//判断是否已经加了好友
	//条件链式操作
	tmp :=model.Contact{};
	DbEngin.Where("ownerid = ?",userid).And("dstid = ?",dstid).
		And("cate = ?",model.CONCAT_CATE_USER).Get(&tmp);
	//如果存在记录说明已经是好友了不用再次添加
	if tmp.Id > 0 {
		return errors.New("该用户已经是您的好友");
	}
	//启动事务
	session := DbEngin.NewSession();
	session.Begin();
	//插入自己的数据
	_, e2 := session.InsertOne(model.Contact{
		Ownerid:userid,
		Dstobj:dstid,
		Cate:model.CONCAT_CATE_USER,
		Createat:time.Now(),
	});
	//插入对方的数据
	_, e3 := session.InsertOne(model.Contact{
		Ownerid:dstid,
		Dstobj:userid,
		Cate:model.CONCAT_CATE_USER,
		Createat:time.Now(),
	});
	//两个都创建成功了才完成否则回滚
	if e2 == nil && e3 == nil {
		//提交
		session.Commit();
		return nil;
	}else{
		//事务回滚
		session.Rollback();
		if e2!=nil {
			return e2;
		}else {
			return e3;
		}
	}
}

//搜索当前用户归属的群组
func (service *ContactService)SearchComunity(userId int64)([]model.Community){
	conconts := make([]model.Contact,0);//好友列表容器
	comIds := make([]int64,0);//存储对端id容器
	//查找我归属的群组
	DbEngin.Where("ownerid = ? and cate = ?",userId,model.CONCAT_CATE_COMUNITY).Find(&conconts);
	for _, v := range conconts {
		comIds = append(comIds,v.Dstobj);
	}
	coms := make([]model.Community,0);//群组表容器
	//判断我归属的群组是否为空
	if len(comIds)==0 {
		return coms
	}
	//返回归属的群组列表
	DbEngin.In("id",comIds).Find(&coms);
	return coms;
}
//获取用户全部群ID
func (service *ContactService)SearchComunityIds(userId int64) (comIds []int64){
	//todo 获取用户全部群ID
	conconts := make([]model.Contact,0);
	comIds = make([]int64, 0);

	DbEngin.Where("ownerid = ? and cate = ?",userId,model.CONCAT_CATE_COMUNITY).Find(&conconts);
	for _,v :=range conconts{
		comIds = append(comIds,v.Dstobj);
	}
	return comIds;
}

//加群
func (service *ContactService) JoinCommunity(userId,comId int64)  error{
	//拼接函数
	cot := model.Contact{
		Ownerid:userId,
		Dstobj:comId,
		Cate:model.CONCAT_CATE_COMUNITY,
	};
	//判断该数据是否存在
	DbEngin.Get(&cot);
	if cot.Id == 0 {
		//设置创建时间并进行储存
		cot.Createat = time.Now();
		_, err := DbEngin.InsertOne(&cot);
		return err;
	}else{
		return nil;
	}
}

//建群
func (service *ContactService) CreateCommunity(comm model.Community) (ret model.Community,err error) {
	//字段校验
	if len(comm.Name)==0 {
		err = errors.New("缺少群名称");
		return ret,err;
	}
	if comm.Ownerid==0 {
		err = errors.New("请先登录");
		return ret,err;
	}
	//查询已经创建的群组
	com := model.Community{
		Ownerid:comm.Ownerid,
	};
	num,err:=DbEngin.Count(&com);
	if num>=5 {
		err = errors.New("一个用户只能创5个群");
		return com,err;
	}else {
		//符合创建条件的执行创建的事务操作
		comm.Createat = time.Now();//创建时间
		session := DbEngin.NewSession();
		session.Begin();
		_, err = session.InsertOne(&comm);
		if err != nil {
			session.Rollback();
			return com,err;
		}
		_, err = session.InsertOne(//把创建人加入到群组中
			model.Contact{
				Ownerid:  comm.Ownerid,
				Dstobj:   comm.Id,
				Cate:     model.CONCAT_CATE_COMUNITY,
				Createat: time.Now(),
			});
		if err != nil {
			session.Rollback();
		}else {
			session.Commit();
		}
		return com,err;
	}
}

//查找好友
func (service *ContactService)SearchFriend(userId int64) ([]model.User){
	conconts := make([]model.Contact,0);
	objIds := make([]int64,0);
	//查找我的所有好友
	DbEngin.Where("ownerid = ? and cate = ?",userId,model.CONCAT_CATE_USER).Find(&conconts);
	for _,v :=range conconts{
		objIds = append(objIds,v.Dstobj);
	}
	coms := make([]model.User,0)
	if len(objIds)==0{
		return coms
	}
	//遍历获取所有的好友数据
	DbEngin.In("id",objIds).Find(&coms);
	return coms;
}
