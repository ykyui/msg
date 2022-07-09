package controller

import (
	"encoding/json"
	"io/ioutil"
	"msg/auth"
	"msg/mongodb"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	countInit.Add(1)
	go func() {
		addPublicApi("/signup", signup, []string{http.MethodPost})
		addPublicApi("/login", login, []string{http.MethodPost})
		countInit.Done()
	}()
}

func signup(rw http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var userInfo mongodb.UserInfo
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if !userInfo.SignupVerify() {
		return nil, userSignupFormatError, errorUserSignupFormat
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), 14)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	userInfo.Password = string(bytes)
	id, err := userInfo.Insert()
	return id, http.StatusOK, err
}

func login(rw http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var userInfo mongodb.UserInfo
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	providedPassword := userInfo.Password
	if err != nil {
		return nil, http.StatusForbidden, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(providedPassword)); err != nil {
		return nil, http.StatusForbidden, err
	}
	jwt, err := auth.GenerateToken(userInfo.Username)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return struct {
		Token string `json:"token"`
	}{jwt}, http.StatusOK, nil
}
