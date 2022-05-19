package main

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

/**
 * 生成token
 */
func macke(info *EandPType) (token string, err error) {
	claims := jwt.MapClaims{
		"email":    info.Email,
		"password": info.Password,
	}
	then := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(then)
	token, err = then.SignedString([]byte("gettoken"))
	return
}

/**
 * 按照这样的规则解析
 */
func secret() jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		return []byte("gettoken"), nil
	}
}

/**
 * 解析token
 */
func parseToken(token string) (info *EandPType, err error) {
	info = &EandPType{}
	tokn, err1 := jwt.Parse(token, secret())
	if err1 != nil {
		err = err1
		return
	}

	claim, ok := tokn.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("token analysis error")
		return
	}
	if !tokn.Valid {
		err = errors.New("token error")
		return
	}

	//强行转换为string类型
	info.Email = claim["email"].(string)
	//强行转换为string类型
	info.Password = claim["password"].(string)
	return
}
