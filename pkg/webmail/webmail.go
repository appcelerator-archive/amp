package webmail

import (
	"fmt"
	"net/http"
	"strings"

	conf "github.com/appcelerator/amp/pkg/config"
	//"github.com/appcelerator/amp/api/rpc/account
	"github.com/gorilla/mux"
)

func StartListener() {
	fmt.Printf("init web\n")
	config := conf.GetRegularConfig(false)
	port := config.WebMailServerPort
	go func() {
		router := mux.NewRouter().StrictSlash(true)
		router.HandleFunc("/v1/ampaccount/{account}/confirm/{token}", confirmAccount)
		router.HandleFunc("/v1/ampaccount/{account}/resetpassword/{token}", resetPassword)
		fmt.Printf("Start web server on port %s\n", port)
		http.ListenAndServe(":"+port, router)
	}()
}

func confirmAccount(w http.ResponseWriter, req *http.Request) {
	param := mux.Vars(req)
	account := param["account"]
	token := param["token"]
	fmt.Printf("received confirmAccount request for account=%s token=%s\n", account, token)
	var err error
	//err=account.Verification(account, token)
	if err != nil {
		ret := strings.Replace(accountVerificationBodyKo, "{accountName}", account, -1)
		ret = strings.Replace(ret, "{error}", err.Error(), -1)
		fmt.Fprintf(w, ret)
	}
	ret := strings.Replace(accountVerificationBodyOk, "{accountName}", account, -1)
	fmt.Fprintf(w, ret)
}

func resetPassword(w http.ResponseWriter, req *http.Request) {
	param := mux.Vars(req)
	account := param["account"]
	token := param["token"]
	fmt.Printf("received resetPassword request for account=%s token=%s\n", account, token)
	var err error
	//err=account.ResetPwd(account, totken)
	if err != nil {
		ret := strings.Replace(accountResetPasswordKo, "{accountName}", account, -1)
		ret = strings.Replace(ret, "{error}", err.Error(), -1)
		fmt.Fprintf(w, ret)
	}
	ret := strings.Replace(accountResetPasswordOk, "{accountName}", account, -1)
	fmt.Fprintf(w, ret)
}
