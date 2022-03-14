#### Client结构
* CheckRedirect用来控制重定向
* Jar用来设置和存放Cookie
* Timeout用来设置超时时间


#### 片段示例

* 通过Get访问
```golang
// 默认使用http包中定义的http.DefaultClient

resp, err := http.Get("http://www.baidu.com")
if err != nil {
    fmt.Println(err.Error())
}

defer resp.Body.Close()

b,err := ioutils.ReadAll(resp.Body)
if err != nil {
    fmt.Println(err.Error())
}

fmt.Println(string(b))
```
* 使用Post和PostForm方法
```golang 

host := "http://www.baidu.com"
contentType :=map[string]string{
    "application/json":"{a:1}", // json
    "application/octet-stream":"", // 二进制文件
    "application/x-www-form-urlencoded":"a=1", //  表单
}

for k,v := range contentType {
    resp, err := http.Post(host,k,bytes.NewBufferString(v))
    if err != nil {
        continue
    }

    defer resp.Body.Close()

    b,err := ioutils.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err.Error())
    }

    fmt.Println(string(b))
}

// PostForm默认使用表单格式x-www-form-urlencoded解析数据，同时可添加Url参数
// 可替换Get使用
resp, err := http.PostForm(host,url.Values{"name":[]string{"winnerxue","jack"}})
if err != nil {
    return 
}

defer resp.Body.Close()

b,err := ioutils.ReadAll(resp.Body)
if err != nil {
    fmt.Println(err.Error())
}

fmt.Println(string(b))
```
* 服务器和客户端使用Cookie
```golang
//  客户端方式1： 
//  ==============================================================================
req, err := http.NewRequest("get", "http://192.168.2.109:8080", nil)
if err != nil {
    return
}

c := &http.Cookie{Name: "name", Value: "bbb"}

// 注意：客户端发送给服务器使用Cookie作为key，而不是Set-Cookie，Set-Cookie是服务器返回的key
req.Header.Add("Cookie", c.String()) 
// req.Header.Set("Cookie", "name=bbb")

resp, err := (&http.Client{}).Do(req)
if err != nil {
    return
}

defer resp.Body.Close()

for _, v := range resp.Cookies() {
    fmt.Println(v.Name, v.Value)
}

// 客户端方式2：
//  ==============================================================================
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


//  服务器
//  ==============================================================================
go func() {
    http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println(r.Header)

        if cookie, err := r.Cookie("name"); err != nil {
            http.SetCookie(w, &http.Cookie{Name: "name", Value: "aaa"})
        } else {
            cookie.Value = "123"
            http.SetCookie(w, cookie)
        }
    }))
}()

ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)

defer cancel()

for {
    select {
    case <-ctx.Done():
        break
    default:
    }
}
```
