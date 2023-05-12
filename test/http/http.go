package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	rs1 := &RealServer{Addr: "127.0.0.1:2003"}
	rs1.Run()
	rs2 := &RealServer{Addr: "127.0.0.1:2004"}
	rs2.Run()

	//监听关闭信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

type RealServer struct {
	Addr string
}

var upgrader = websocket.Upgrader{}

func (r *RealServer) WebSocketHandler(w http.ResponseWriter, req *http.Request) {
	println("websocket")
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s", message)

		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func (r *RealServer) Run() {
	log.Println("Starting httpserver at " + r.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/test_http_string/hello", r.HelloHandler)
	mux.HandleFunc("/test_http_string/ws", r.WebSocketHandler) // 添加 WebSocket handler
	mux.HandleFunc("/base/error", r.ErrorHandler)
	mux.HandleFunc("/test_http_string/test_http_string/aaa", r.TimeoutHandler)
	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}

// func (r *RealServer) HelloHandler(w http.ResponseWriter, req *http.Request) {
// 	//127.0.0.1:8008/abc?sdsdsa=11
// 	//r.Addr=127.0.0.1:8008
// 	//req.URL.Path=/abc
// 	//fmt.Println(req.Host)
// 	upath := fmt.Sprintf("http://%s%s\n", r.Addr, req.URL.Path)
// 	realIP := fmt.Sprintf("RemoteAddr=%s,X-Forwarded-For=%v,X-Real-Ip=%v\n", req.RemoteAddr, req.Header.Get("X-Forwarded-For"), req.Header.Get("X-Real-Ip"))
// 	header := fmt.Sprintf("headers =%v\n", req.Header)
// 	io.WriteString(w, upath)
// 	io.WriteString(w, realIP)
// 	io.WriteString(w, header)

// }

func (r *RealServer) ErrorHandler(w http.ResponseWriter, req *http.Request) {
	upath := "error handler"
	w.WriteHeader(500)
	io.WriteString(w, upath)
}

func (r *RealServer) TimeoutHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(6 * time.Second)
	upath := "timeout handler"
	w.WriteHeader(200)
	io.WriteString(w, upath)
}

func (r *RealServer) HelloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("hello")
	// 设置服务器 IP 地址
	serverIP := r.Addr

	// 获取请求的URL路径、X-Forwarded-For头部信息、X-Real-Ip头部信息，以及所有的HTTP请求头部信息
	upath := fmt.Sprintf("http://%s%s", r.Addr, req.URL.Path)
	realIP := fmt.Sprintf("RemoteAddr=%s,X-Forwarded-For=%v,X-Real-Ip=%v", req.RemoteAddr, req.Header.Get("X-Forwarded-For"), req.Header.Get("X-Real-Ip"))
	header := fmt.Sprintf("headers =%v", formatHeaders(req.Header))

	// 构造 HTML 页面
	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
			<head>
				<title>Server IP: %s</title>
			</head>
			<body>
				<h2>Server IP: %s</h2>
				<h3>URL Path:</h3>
				<pre>%s</pre>
				<h3>X-Forwarded-For and X-Real-Ip:</h3>
				<pre>%s</pre>
				<h3>HTTP Request Headers:</h3>
				<pre>%s</pre>
			</body>
		</html>
	`, serverIP, serverIP, upath, realIP, header)

	// 设置 HTTP 响应头和正文
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, html)
}

// 格式化 HTTP 请求头
func formatHeaders(headers http.Header) string {
	var headerString string
	for k, v := range headers {
		headerString += fmt.Sprintf("%s: %s\n", k, strings.Join(v, ", "))
	}
	return headerString
}
