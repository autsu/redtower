# UML
![](https://github.com/autsu/redtower/blob/master/doc/uml.svg)

# Example:
(View the `example` folder in detail)
## Server
```go
var (
	EchoMsg = server.GenMsgTyp("echo")
	HTTPMsg = server.GenMsgTyp("HTTP")
)

// EchoHandler 自己实现 Echo message 的处理函数
type EchoHandler struct {
	server.BasicHandler
}

func (e *EchoHandler) Handle(r *server.Request) {
	data := r.Data()

	msg := server.NewMessage(data, r.MsgType())
	_, err := r.Conn().Send(msg)
	if err != nil {
		log.Println(err)
		return
	}
}

// HttpEchoPostFormHandler 回显 post 表单数据（该 handler 仅仅是用来测试 http 请求）
type HttpEchoPostFormHandler struct {
	server.BasicHandler
}

func (h *HttpEchoPostFormHandler) Handle(req *server.Request) {
	bsr := bytes.NewReader(req.Data())
	br := bufio.NewReader(bsr)

	r, err := http.ReadRequest(br)
	if err != nil {
		log.Println(err)
		msg := server.NewErrorMessage([]byte(err.Error()))
		req.Conn().Send(msg)
		return
	}
	log.Printf("%+v", r)

	r.ParseForm()
	form := r.Form.Encode()
	msg := server.NewMessage([]byte(form), HTTPMsg)
	req.Conn().Send(msg)
}
```

## client
### HTTP
```go
func main() {
	ctx := context.Background()
	conn, err := client.
		NewClientWithTCP("localhost", "8080").
		Init(ctx)

	if err != nil {
		log.Fatalln(err)
	}

	data := `POST /?123=456 HTTP/1.1

User-Agent: PostmanRuntime/7.28.1
Accept: */*
Postman-Token: e1e457b7-d713-443d-8022-04a2d1d9697a
Host: 127.0.0.1:8080
Accept-Encoding: gzip, deflate, br
Connection: keep-alive
Content-Length: 0`

	msg := server.NewMessage([]byte(data), example.HTTPMsg)
	conn.Send(msg)

	msg, err = conn.Receive()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(msg.Data()))
}
```