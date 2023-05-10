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
)

func main() {
	rs1 := &RealServer{Addr: "127.0.0.1:2003"}
	rs1.Run()
	rs2 := &RealServer{Addr: "127.0.0.1:2004"}
	rs2.Run()

	//监听关闭信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

type RealServer struct {
	Addr string
}

func (r *RealServer) Run() {
	log.Println("Starting httpserver at " + r.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.HelloHandler)
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

func (r *RealServer) HelloHandler(w http.ResponseWriter, req *http.Request) {
	// 设置服务器 IP 地址
	serverIP := r.Addr

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
				<h3>X-Forwarded-For:</h3>
				<pre>%s</pre>
				<h3>X-Real-Ip:</h3>
				<pre>%s</pre>
				<h3>HTTP Request Headers:</h3>
				<pre>%s</pre>
			</body>
		</html>
	`, serverIP, serverIP, req.URL.Path, req.Header.Get("X-Forwarded-For"), req.Header.Get("X-Real-Ip"), formatHeaders(req.Header))

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
