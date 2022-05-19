package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

const (
	userName = "root"
	password = "!QAZ1qaz"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "login"
)

var DB *sql.DB

func main() {
	connDB()
	router := gin.Default()
	router.POST("/v1/signup", func(c *gin.Context) { signup(c) })
	router.POST("/v1/signin", func(c *gin.Context) { signin(c) })
	router.GET("/v1/profile", Authorization(func(c *gin.Context, email string, password string) { profile(c, email, password) }))
	router.POST("/v1/profile/update", Authorization(func(c *gin.Context, email string, password string) { update(c, email, password) }))
	router.Run(":8000")
}

/**
 * signup
 */
func signup(c *gin.Context) {
	var json UserType
	if err := c.Bind(&json); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	var jsonParameter UserType
	code := checkUserType(&json, &jsonParameter)
	if code == OK {
		//加密处理
		hash, err := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"code": 500, "message": err.Error(), "data": ""})
		} else {
			json.Password = string(hash)
			fmt.Println(json.Password)
			err = insertUsersDB(json)
			if err != nil {
				c.JSON(500, gin.H{"code": 500, "message": err.Error(), "data": ""})
			} else {
				var cNandEType NandEType
				cNandEType.FirstName = json.FirstName
				cNandEType.LastName = json.LastName
				cNandEType.Email = json.Email
				c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK", "data": cNandEType})
			}
		}
	} else if code == DB_ERROR {
		c.JSON(500, gin.H{"code": 500, "message": "DB Error", "data": ""})
	} else if code == FORMAT_ERROR {
		c.JSON(400, gin.H{"code": 400, "message": jsonParameter, "data": ""})
	}
}

/**
 * signin
 */
func signin(c *gin.Context) {
	var json EandPType
	if err := c.Bind(&json); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	if !checkEandP(c, json.Email, json.Password) {
		return
	}

	var user UserType
	if !checkEffective(c, json.Email, json.Password, &user) {
		return
	}

	// token生成
	token, err := macke(&json)
	fmt.Println(token)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error(), "data": ""})
	} else {
		err = redisManager(json.Email, token)
		if err != nil {
			c.JSON(500, gin.H{"code": 500, "message": err.Error(), "data": ""})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK", "data": token})
		}
	}
}

/**
 * profile
 */
func profile(c *gin.Context, email string, password string) {
	var user UserType
	if !checkEffective(c, email, password, &user) {
		return
	}

	var cNandEType NandEType
	cNandEType.FirstName = user.FirstName
	cNandEType.LastName = user.LastName
	cNandEType.Email = user.Email
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK", "data": cNandEType})
}

/**
 * update
 */
func update(c *gin.Context, email string, password string) {
	var user UserType
	if !checkEffective(c, email, password, &user) {
		return
	}

	var json NameType
	if err := c.Bind(&json); err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	var fl NameType
	if !checkFirstName(json.FirstName) {
		fl.FirstName = "FirstName error"
	}

	if !checkLastName(json.LastName) {
		fl.LastName = "LastName error"
	}
	if fl.FirstName != "" || fl.LastName != "" {
		c.JSON(400, gin.H{"code": 400, "message": fl, "data": ""})
		return
	}

	if user.FirstName == json.FirstName && user.LastName == json.LastName {
		c.JSON(400, gin.H{"code": 400, "message": "The name isn't changed", "data": ""})
		return
	}

	err := updateUsersDB(json.FirstName, json.LastName, email)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": err.Error(), "data": ""})
	} else {
		var cNandEType NandEType
		cNandEType.FirstName = json.FirstName
		cNandEType.LastName = json.LastName
		cNandEType.Email = email
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "OK", "data": cNandEType})
	}
}

/**
 * Authorization
 */
func Authorization(f func(c *gin.Context, email string, password string)) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenKey := c.Request.Header.Get("Authorization")
		info, err := parseToken(tokenKey)
		if err != nil {
			c.JSON(400, gin.H{"code": 400, "message": err.Error(), "data": ""})
		} else {
			f(c, info.Email, info.Password)
		}
	}
}

/**
 * redis
 */
func redisManager(email string, token string) error {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword("trechina"))

	if conn != nil {
		defer conn.Close()
	}

	if err == nil {
		_, err = conn.Do("Set", email, token)
		if err != nil {
			fmt.Println("hset err1 = ", err)
		}
	} else {
		fmt.Println("hset err2 = ", err)
	}

	// r, err := redis.String(conn.Do("Get", email))
	// if err != nil {
	// 	fmt.Println("set err3 = ", err)
	// 	return err
	// }
	// fmt.Println("操作 ok", r)

	return err
}

/**
 * 检查email与password格式
 */
func checkEandP(c *gin.Context, email string, password string) bool {
	var ep EandPType
	count, code := checkEmail(email)
	if code == DB_ERROR {
		c.JSON(500, gin.H{"code": 500, "message": "DB Error", "data": ""})
		return false
	}
	if code == FORMAT_ERROR {
		ep.Email = "Email error(****** Unclear requirements ******)"
	} else if count == 0 {
		ep.Email = "Email not found"
	}

	if !checkPassword(password) {
		ep.Password = "Password error"
	}

	if ep.Email != "" || ep.Password != "" {
		c.JSON(400, gin.H{"code": 400, "message": ep, "data": ""})
		return false
	}
	return true
}

/**
 * 检查有效性
 */
func checkEffective(c *gin.Context, email string, password string, user *UserType) bool {

	// 数据库查询(checkEmail判断过是否存在Email)
	err := selectPasswordByEmailDB(email, user)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error(), "data": ""})
		return false
	}

	// 密码验证
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error(), "data": ""})
		return false
	}

	return true
}
