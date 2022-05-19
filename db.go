package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

/**
 * 链接数据库
 */
func connDB() bool {
	path := strings.Join([]string{userName, ":", password, "@(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	DB, _ = sql.Open("mysql", path)
	DB.SetConnMaxIdleTime(100)
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail")
		return false
	} else {
		fmt.Println("connect success")
		return true
	}
}

/**
 * 通过email和password查找
 */
func selectUsersDB(email string, password string) (UserType, error) {
	var user UserType

	count := 0
	err := DB.QueryRow("SELECT COUNT(1) FROM user WHERE email = ? AND password = ?", email, password).Scan(&count)
	if err != nil {
		return user, err
	} else if count == 0 {
		return user, errors.New("no found!")
	}

	err1 := DB.QueryRow("SELECT * FROM user WHERE email = ? AND password = ?", email, password).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err1 != nil {
		fmt.Println("select error:", err1.Error())
		return user, err1
	} else {
		fmt.Println("select:", user)
		return user, nil
	}
}

/**
 * 通过Email查找password
 */
func selectPasswordByEmailDB(email string, user *UserType) error {
	err := DB.QueryRow("SELECT * FROM user WHERE email = ?", email).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Password)
	return err
}

/**
 * 查找email数量
 */
func emailCountUsersDB(email string) (int, error) {
	count := 0
	err := DB.QueryRow("SELECT COUNT(1) FROM user WHERE email = ?", email).Scan(&count)
	if err != nil {
		return count, err
	} else {
		return count, nil
	}
}

/**
 * 插入
 */
func insertUsersDB(user UserType) error {
	_, err := DB.Exec("INSERT INTO user(first_name, last_name, email, password) VALUES (?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, user.Password)
	return err
}

/**
 * 修改
 */
func updateUsersDB(firstName string, lastName string, email string) error {
	_, err := DB.Exec("UPDATE user SET first_name = ?, last_name = ? WHERE email = ?", firstName, lastName, email)
	return err
}
