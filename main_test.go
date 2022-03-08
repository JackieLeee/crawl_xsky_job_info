package main

import "testing"

/**
 * @Author  Flagship
 * @Date  2022/3/8 22:03
 * @Description
 */

func TestGetToken(t *testing.T) {
	cookies, err := GetCookies()
	if err != nil {
		t.Errorf("get cookie failed, err:%s", err.Error())
		return
	}
	token, err := ParseCsrfToken(cookies)
	if err != nil {
		t.Errorf("parse token failed, err:%s", err.Error())
		return
	}
	t.Logf("test success, token: %s", token)
}

func TestGetJobListFromApi(t *testing.T) {
	jobInfoList, err := GetJobListFromApi()
	if err != nil {
		t.Errorf("get cookie failed, err:%s", err.Error())
		return
	}
	t.Logf("test success, job info list len: %d", len(jobInfoList))
}
