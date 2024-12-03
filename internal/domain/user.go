package domain

import "time"

//User 领域对象 DDD entity

type User struct {
	Id       int64
	Email    string
	PassWord string
	Phone    string
	NickName string

	WechatInfo WechatInfo
	//
	Ctime time.Time
	//todo WeChatInfo
}
