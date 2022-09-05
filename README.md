# UML
![](https://github.com/autsu/redtower/blob/master/doc/uml.svg)

# Example:

(在 `example` 文件夹中查看更多示例，testdata 中演示了一个粘包/拆包处理效果的 demo)
## Global

```go
// 自定义消息类型，需要保证不与库的内置消息类型（ErrorMsg = -1 和 HearbeatMsg = -2 ）的值相同
const HTTPMsg = 10	
const EchoMsg = 11
```



## Server

```go
func main() {
	s := server.NewTCPServer("localhost", "7788", "test1",
		func(req *server.Request) error {
			typ := req.MsgType()
			switch typ {
			case example.HTTPMsg:
				handleHTTP(req)
			case server.ErrorMsg:
				log.Printf("error: %v\n", string(req.Data()))
			default:
				log.Println("Unknown msg typ")
				log.Println(req.Data())
			}
			return nil
		})
	if err := s.Run(context.Background()); err != nil {
		panic(err)
	}
}

func handleHTTP(req *server.Request) {
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
	msg := server.NewMessage([]byte(form), example.HTTPMsg)
	req.Conn().Send(msg)
}
```

## client
### HTTP
```go
func main() {
	ctx := context.Background()
	conn, err := client.
		NewClientWithTCP("localhost", "7788").
		Init(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

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