package model

import (
	"crypto/md5"
	"fmt"

	database "github.com/xunlbz/go-restful-template/pkg/database"
	"gorm.io/gorm"
)

// User 用户对象
type User struct {
	gorm.Model
	Name     string
	Password string
}

// UserDao 方法接口
type UserDao interface {
	GetAll(pageNo int, pageSize int) []User
	GetOne(username string) *User
	Insert(username string, password string) User
	Update(username string, password string) User
	Delete(userID uint) bool
	Exists(username string) bool
}

func init() {
	database.AddToMigration(&User{})
}

// NewUserDao 获取UserDao
func NewUserDao() UserDao {
	return &userDaoImpl{db: database.DEFAULTDB}
}

//UserDaoImpl  用户db操作接口
type userDaoImpl struct {
	db *gorm.DB
}

func (userDao *userDaoImpl) Insert(username string, password string) User {
	user := User{Name: username, Password: fmt.Sprintf("%x", md5.Sum([]byte(password)))}
	userDao.db.Create(&user)
	return user
}
func (userDao *userDaoImpl) Update(username string, password string) User {
	user := *userDao.GetOne(username)
	if user.ID != 0 {
		updateUser := User{}
		if username != "" {
			updateUser.Name = username
		}
		if password != "" {
			updateUser.Password = password
		}
		userDao.db.Model(&user).Updates(updateUser)
	}
	return user

}
func (userDao *userDaoImpl) Delete(userID uint) bool {
	userDao.db.Delete(new(User), userID)
	return true
}

func (userDao *userDaoImpl) GetOne(username string) *User {
	var user User
	userDao.db.First(&user, "name = ?", username)
	return &user
}

func (userDao *userDaoImpl) GetAll(pageNo int, pageSize int) (users []User) {
	userDao.db.Scopes(database.Paginate(pageNo, pageSize)).Find(&users)
	return users
}

func (userDao *userDaoImpl) Exists(username string) bool {
	user := userDao.GetOne(username)
	return user.ID > 0
}
