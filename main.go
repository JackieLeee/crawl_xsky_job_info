package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/beego/beego/v2/core/logs"

	"github.com/flagship575/crawl_xsky_job_info/message"
)

/**
 * @Author  Flagship
 * @Date  2022/3/7 20:34
 * @Description
 */

const (
	KBaseUrl = "https://xskydata.jobs.feishu.cn"

	KCsrfTokenUrl = KBaseUrl + "/api/v1/csrf/token"
	KJobInfoUrl   = KBaseUrl + "/api/v1/search/job/posts"

	KHeaderAtsxCsrfToken = "atsx-csrf-token"
	KHeaderXCsrfToken    = "x-csrf-token"
	KHeaderWebsitePath   = "website-path"

	KSearchTypeSchool = "school"

	KFileName = "job_info.json"
)

var (
	TargetSearch = KSearchTypeSchool
)

func main() {
	jobList, err := GetJobListFromApi()
	if err != nil {
		logs.Warn("get job list from api failed, err:", err)
		return
	}

	//保存文件
	if err := SaveToFile(jobList); err != nil {
		logs.Warn("save job info to file failed, err:", err)
		return
	}
	logs.Debug("crawl finish, filename:", KFileName)
}

//GetJobListFromApi 从Api中获取职位列表
func GetJobListFromApi() ([]*message.JobInfo, error) {
	//获取cookie
	cookies, err := GetCookies()
	if err != nil {
		logs.Warn("get cookies failed, err:", err)
		return nil, err
	}
	logs.Debug("get cookie success")

	//首次获取职位信息，目的是拿到总数
	limit := 10
	jobInfo, err := GetJobInfo(cookies, 0, limit)
	if err != nil {
		logs.Warn("get job info failed, err:", err)
		return nil, err
	}
	logs.Debug("first get job info success")

	//记录总数
	totalCount := jobInfo.Data.Count
	logs.Debug("job info total count:", totalCount)

	errMap := make(map[int]error, 0)
	offset2JobInfo := make(map[int]*message.JobInfoResp, 1)
	offset2JobInfo[0] = jobInfo
	//多线程调用
	wg := sync.WaitGroup{}
	for offset := limit; offset+1 < totalCount; offset += limit {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			info, err := GetJobInfo(cookies, offset, limit)
			if err != nil {
				errMap[offset] = err
				logs.Warn("get job info failed, err:", err)
				return
			}
			offset2JobInfo[offset] = info
			logs.Debug("get job info success, start index: %d, len:%d", offset, len(info.Data.JobPostList))
		}(offset)
	}
	wg.Wait()

	if len(errMap) > 0 {
		for offset, err := range errMap {
			logs.Warn("err in go routine[offset:%d], err:%s", offset, err.Error())
		}
		return nil, errors.New("some err in go routines")
	}

	jobList := make([]*message.JobInfo, 0, totalCount)
	for _, job := range offset2JobInfo {
		jobList = append(jobList, job.Data.JobPostList...)
	}
	return jobList, nil
}

//SaveToFile 保存到文件
func SaveToFile(list []*message.JobInfo) error {
	//json序列化
	dataByte, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		logs.Warn("marshal job list failed, err:", err)
		return err
	}
	//写入文件并赋予读写权限
	if err := ioutil.WriteFile(KFileName, dataByte, 0666); err != nil {
		logs.Warn("write bytes to file failed, err:", err)
		return err
	}
	return nil
}

//GetJobInfo 获取职位信息
func GetJobInfo(cookies []*http.Cookie, offset, limit int) (*message.JobInfoResp, error) {
	//获取csrf_token
	csrfToken, err := ParseCsrfToken(cookies)
	if err != nil {
		logs.Warn("parse csrf_token from cookies failed, err:", err)
		return nil, err
	}

	//构造请求
	client := &http.Client{}
	params := message.JobInfoReq{
		Limit:  limit,
		Offset: offset,
	}
	b, err := json.Marshal(params)
	if err != nil {
		logs.Warn("marshal request params failed, err:", err)
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, KJobInfoUrl, bytes.NewBuffer(b))
	if err != nil {
		logs.Warn("create new request failed, err:", err)
		return nil, err
	}
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
	request.Header.Set(KHeaderXCsrfToken, csrfToken)
	request.Header.Set(KHeaderWebsitePath, TargetSearch)

	//执行请求
	response, err := client.Do(request)
	if err != nil {
		logs.Warn("post job info api failed, err:", err)
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("post job info api failed, err")
		logs.Warn("post job info api get status code:", response.StatusCode)
		return nil, err
	}

	//处理请求结果
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logs.Warn("read result from response failed, err:", err)
		return nil, err
	}
	postResult := &message.JobInfoResp{}
	if err := json.Unmarshal(body, postResult); err != nil {
		logs.Warn("unmarshal body failed, err: %s", err.Error())
		return nil, err
	}

	return postResult, nil
}

//ParseCsrfToken 解析出csrf_token
func ParseCsrfToken(cookies []*http.Cookie) (string, error) {
	csrfToken := ""
	for _, cookie := range cookies {
		if cookie.Name == KHeaderAtsxCsrfToken {
			atsxCsrfToken := cookie.Value
			csrfToken = atsxCsrfToken[:len(atsxCsrfToken)-3] + "="
			break
		}
	}

	if csrfToken == "" {
		errMsg := "get empty csrf_token from cookie"
		logs.Warn(errMsg)
		return "", errors.New(errMsg)
	}

	return csrfToken, nil
}

//GetCookies 获取cookie
func GetCookies() ([]*http.Cookie, error) {
	client := &http.Client{}

	response, err := client.Post(KCsrfTokenUrl, "application/json", nil)
	if err != nil {
		logs.Warn("post csrf api failed, err:", err)
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("post csrf api failed, err")
		logs.Warn("post csrf api get status code:", response.StatusCode)
		return nil, err
	}

	return response.Cookies(), nil
}
