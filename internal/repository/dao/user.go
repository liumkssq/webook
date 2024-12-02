package dao

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

/*
	数据访问层
*/
//  直接对应数据的表 model PO(persistence object) 持久层对象
type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	Nickname string

	Phone sql.NullString `gorm:"unique"`

	// 微信的字段 todo
	//WechatUnionID sql.NullString `gorm:"type=varchar(1024)"`
	//WechatOpenID  sql.NullString `gorm:"type=varchar(1024);unique"`

	Ctime int64
	Utime int64
}

// 定义错误
var (
	ErrUserDuplicate = errors.New("邮箱冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

// Interface 数据访问接口
type UserDAO interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	Insert(ctx context.Context, u User) error
	//todo
	//FindByWechat(ctx context.Context, openID string) (User, error)
}

type DBProvider func() *gorm.DB

type GORMUserDAO struct {
	db *gorm.DB
	p  DBProvider
}

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	//TODO implement me
	panic("implement me")
}

func (dao GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	//TODO implement me
	panic("implement me")
}

func (dao GORMUserDAO) Insert(ctx context.Context, u User) error {
	//TODO implement me
	panic("implement me")
}

func NewGORMUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{db: db}
}
