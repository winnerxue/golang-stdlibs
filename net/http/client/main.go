package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func main() {
	urlString := "http://192.168.2.109:8080"

	// cookie的管理类
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	client := http.Client{
		Jar: jar,
	}

	u, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}

	// 添加cookie到指定指定的url中，在cookie中字段中选择Cookie生效的策略
	client.Jar.SetCookies(u, []*http.Cookie{{Name: "token", Value: "xxxx", MaxAge: 300}})

	postdata := url.Values{}
	postdata.Add("name", "winnerxue")
	postdata.Add("name", "jake")
	fmt.Println(postdata.Encode())

	// 构造请求体
	req, err := http.NewRequest("POST", urlString, strings.NewReader(postdata.Encode()))
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	for _, v := range resp.Cookies() {
		fmt.Println(v.Name, v.Value)
	}

}
