package main

import (
	"regexp"
)

/**
 * 检查参数正确性
 */
func checkUserType(json *UserType, user *UserType) int {

	count, code := checkEmail(json.Email)
	if code == DB_ERROR {
		return code
	}
	if count > 0 {
		user.Email = "Email used"
	} else if code == FORMAT_ERROR {
		user.Email = "Email error(****** Unclear requirements ******)"
	}

	if !checkFirstName(json.FirstName) {
		user.FirstName = "FirstName error"
	}

	if !checkLastName(json.LastName) {
		user.LastName = "LastName error"
	}

	if !checkPassword(json.Password) {
		user.Password = "Password error"
	}

	if user.FirstName != "" || user.LastName != "" || user.Email != "" || user.Password != "" {
		return FORMAT_ERROR
	} else {
		return OK
	}
}

/**
 * 检查firstName正确性
 */
func checkFirstName(firstName string) (isOk bool) {
	isOk = true
	if len(firstName) <= 0 || len(firstName) > 64 {
		isOk = false
	}
	return
}

/**
 * 检查lastName正确性
 */
func checkLastName(lastName string) (isOk bool) {
	isOk = true
	if len(lastName) <= 0 || len(lastName) > 64 {
		isOk = false
	}
	return
}

/**
 * 检查password正确性
 */
func checkPassword(password string) (isOk bool) {
	isOk = true
	if len(password) >= 6 && len(password) <= 16 {
		match, _ := regexp.MatchString(`^[0-9A-Za-z]+$`, password)
		match1, _ := regexp.MatchString(`[A-Za-z]+`, password)
		match2, _ := regexp.MatchString(`[0-9]+`, password)
		if !match || !match1 || !match2 {
			isOk = false
		}
	} else {
		isOk = false
	}
	return
}

/**
 * 检查email正确性
 */
func checkEmail(email string) (int, int) {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9][\w\.-]*[a-zA-Z0-9]@[a-zA-Z0-9][\w\.-]*[a-zA-Z0-9]\.[a-zA-Z][a-zA-Z\.]*[a-zA-Z]$`, email)
	if !match {
		return 0, FORMAT_ERROR
	}
	count, err := emailCountUsersDB(email)
	if err != nil {
		return count, DB_ERROR
	} else {
		return count, OK
	}
}
