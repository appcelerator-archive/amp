package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/appcelerator/amp/data/accounts"
	"github.com/howeyc/gopass"
	"google.golang.org/grpc"
)

func getName() (in string) {
	fmt.Scanln(&in)
	username = strings.TrimSpace(in)
	err := accounts.CheckName(in)
	if err != nil {
		mgr.Warn(err.Error())
		return getName()
	}
	return
}

func getEmailAddress() (email string) {
	fmt.Print("email: ")
	fmt.Scanln(&email)
	email = strings.TrimSpace(email)
	_, err := accounts.CheckEmailAddress(email)
	if err != nil {
		mgr.Warn(err.Error())
		return getEmailAddress()
	}
	return
}

func getToken() (token string) {
	fmt.Print("token: ")
	fmt.Scanln(&token)
	token = strings.TrimSpace(token)
	return
}

func getPassword() (password string) {
	fmt.Print("password: ")
	pw, err := gopass.GetPasswd()
	if err != nil {
		mgr.Warn(err.Error())
		return getPassword()
	}
	password = string(pw)
	password = strings.TrimSpace(password)
	err = accounts.CheckPassword(password)
	if err != nil {
		mgr.Warn(grpc.ErrorDesc(err))
		return getPassword()
	}
	return
}

func convertTime(in int64) time.Time {
	intVal, err := strconv.ParseInt(strconv.FormatInt(in, 10), 10, 64)
	if err != nil {
		mgr.Warn(err.Error())
	}
	out := time.Unix(intVal, 0)
	return out
}
