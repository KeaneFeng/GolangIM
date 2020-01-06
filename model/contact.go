package model

import "time"

const (
	CONCAT_CATE_USER = 0x01  //用户
	CONCAT_CATE_COMUNITY = 0x02 //群组
)

//好友群组列表
type Contact struct {
	Id         int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`
	Ownerid       int64	`xorm:"bigint(20)" form:"ownerid" json:"ownerid"`   // 谁的ID
	Dstobj       int64	`xorm:"bigint(20)" form:"dstobj" json:"dstobj"`   // 对端ID
	Cate      int	`xorm:"int(11)" form:"cate" json:"cate"`   // 用户/群组类型
	Memo    string	`xorm:"varchar(120)" form:"memo" json:"memo"`   // 备注
	Createat   time.Time	`xorm:"datetime" form:"createat" json:"createat"`   // 创建时间
}
