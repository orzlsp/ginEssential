package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(10);not null;unique"`
	Password  string `gorm:"size(255);not null"`
}

func main() {
	db := initDB()
	defer db.Close() //关闭延迟
	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		//获取参数
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")
		//数据验证
		if len(telephone) != 11 {
			ctx.JSON(http.StatusUnprocessableEntity,
				gin.H{
					"code": 422,
					"msg":  "手机号必须得11位",
				})
			return
		}
		if len(password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity,
				gin.H{
					"code": 422,
					"msg":  "密码不能少于6位",
				})
			return
		}
		if len(name) == 0 {
			name = Randomstring(10)
		}
		log.Println(name, telephone, password)
		//判断手机号是否存在
		if isTelephoneExist(db, telephone) {
			ctx.JSON(http.StatusUnprocessableEntity,
				gin.H{
					"code": 422,
					"msg":  "用户已存在",
				})
			return
		}
		//创建用户
		newUser := User{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}
		db.Create(&newUser)

		//返回结果
		ctx.JSON(200, gin.H{
			"msg": "注册成功",
		})
	})
	panic(r.Run()) // listen and serve on 0.0.0.0:8080
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	//// SELECT * FROM user WHERE name = 'jinzhu' limit 1;
	db.Where("telephone=?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func Randomstring(n int) string {
	var letters = []byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP")
	result := make([]byte, n)
	//设置随机数种子，加上这行代码，可以保证每次随机都是随机的
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func initDB() *gorm.DB {
	driverName := "mysql"
	host := "192.168.38.191"
	port := "3307"
	database := "ginessentail"
	username := "linan"
	password := "linan020"
	charset := "utf8mb4"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database,err: " + err.Error())
	}
	db.AutoMigrate(&User{})
	return db

}
