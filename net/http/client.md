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
    "application/form-data":"a=1", //  文件表单
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

// PostForm默认使用表单格式x-www-form-urlencoded解析数据，注意：这里url.Value是放到请求体中的，而不是URL中
// 可替换Get使用，避免Get方式的URL大长度大小限制
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
* 常见请求自定义
```golang

// 客户端发送Post请求(类型为：x-www-form-urlencoded)
// 1. URL中有参数
// 2. Body中也有相同key的参数
//  ==============================================================================

urlString := "http://192.168.2.109:8080?name=aaa&age=111"

data := url.Values{}
data.Add("name", "bbb")
data.Set("age", "222")

// 构造请求体
req, err := http.NewRequest("POST", urlString, bytes.NewBufferString(data.Encode()))
if err != nil {
    panic(err)
}

req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

resp, err := (&http.Client{}).Do(req)
if err != nil {
    panic(err)
}

defer resp.Body.Close()

if b, _ := ioutil.ReadAll(resp.Body); err != nil {
    fmt.Println(string(b))
}

// 客户端发送Post请求(类型通过请求头中Header的Content-Type判定)
// 1. URL中有参数
// 2. Body中也有相同key的参数
//  ==============================================================================

urlString := "http://192.168.2.109:8080?name=aaa&age=111"

data := url.Values{}
data.Add("name", "bbb")
data.Set("age", "222")

// 构造请求体
req, err := http.NewRequest("POST", urlString, bytes.NewBufferString(data.Encode()))
if err != nil {
    panic(err)
}

req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

resp, err := (&http.Client{}).Do(req)
if err != nil {
    panic(err)
}

defer resp.Body.Close()

if b, _ := ioutil.ReadAll(resp.Body); err != nil {
    fmt.Println(string(b))
}

// 服务器处理Post请求
// 1. ParseForm会解析URL中的参数到request中的Form字段
// 2. 如果请求头中的Content-Type为x-www-form-urlencoded，则ParseForm会解析Body中的参数到request中的Form字段和PostForm字段，否则PostForm为空值
// 3. 如果请求头中Content-Type为json，则需要使用json.Unmashal或json.NewDecoder处理request.Body
// 4. 如果请求头中Content-Type为form-data，则使用ParseMultipartForm(100)解析文件和表单参数

//  ==============================================================================

// 常规表单
go func() {
    http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if err := r.ParseForm(); err != nil {
            panic(err)
        }

        fmt.Println(r.Form)
        fmt.Println(r.PostForm)
    }))
}()


// 文件表单处理
go func(){
    http.ListenAndServe(":8081",http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
        if err := r.ParseMultipartForm(10 << 20); err != nil {
            return
        }

        f, h, err := r.FormFile("name")
        if err != nil {
            return
        }

        b, err := ioutil.ReadAll(f)
        if err != nil {
            return
        }

        ioutil.WriteFile(h.Filename, b, os.ModePerm)
    }))
}

// json表单解析
go func() {
    http.ListenAndServe(":8082", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var p struct {
            Name string
        }

        // decoder方式
        jdecoder := json.NewDecoder(r.Body)
        if err := jdecoder.Decode(&p); err != nil {
            return
        }

        // unmarshal方式
        s, err := ioutil.ReadAll(r.Body)
        if err != nil {
            return
        }
        if err := json.Unmarshal(s, &p); err != nil {
            return
        }

        fmt.Println(p.Name)
    }))
}()

ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL)

defer cancel()

for {
    select {
    case <-ctx.Done():
        fmt.Println()
        break
    default:
    }
}


```
* Post请求时类型为form-data的请求体数据格式（RFC1867），[引用](https://github.com/codeLee321/http-rfc1867)

```golang
#请求头设置
Content-type: multipart/form-data, boundary=AaB03x                #设置请求头部类型

#请求体设置
--AaB03x                                                          #表单中某个数据的开始，下一部分换行（\r\n）
content-disposition: form-data; name="submitter"                  #表单中参数的Key，下一部分换行（\r\n\r\n）

Joe Blow                                                          #表单中参数的Value，下一部分换行（\r\n）
--AaB03x                                                          #表单中下个数据的开始，下一部分换行（\r\n）
content-disposition: form-data; name="pics"; filename="file1.txt" #表单中文件的Key和文件名称（\r\n）
Content-Type: text/plain                                          #表单中文件类型（\r\n\r\n）

... contents of file1.txt ...                                     #表单中文件的实体部分（\r\n）
--AaB03x--                                                        #表单数据传输结束标示（\r\n）

```
