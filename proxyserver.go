package proxyhost

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
)

//StartServer start a proxy server
func StartServer(port string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Panic(err)
	}

	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}

		go handleClientRequest(client)
	}
}

func RandomPort() int {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	defer lis.Close()
	return lis.Addr().(*net.TCPAddr).Port
}

func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	var b [1024]byte
	n, err := client.Read(b[:])
	if err != nil {
		// log.Println(err)
		return
	}
	var method, host, address string
	index := bytes.IndexByte(b[:], '\n')
	fmt.Sscanf(string(b[:index]), "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		// log.Println(err)
		return
	}

	var ip string
	h, p, err := net.SplitHostPort(hostPortURL.Host) // 如果不带port，会有err
	if err != nil {
		ip, _ = FindIP(hostPortURL.Host)
	} else {
		ip, _ = FindIP(h)
	}

	if hostPortURL.Opaque == "443" { //https访问
		address = ip + ":443"
	} else { //http访问
		if strings.Index(hostPortURL.Host, ":") == -1 { //host不带端口， 默认80
			address = ip + ":80"
		} else {
			address = ip + ":" + p
		}
	}
	// address = "10.1.2.224:6688"
	if len(ip) > 0 {
		fmt.Println(address, ":", host)
	}

	//获得了请求的host和port，就开始拨号吧
	server, err := net.Dial("tcp", address)
	if err != nil {
		// log.Println(err)
		return
	}
	if method == "CONNECT" {
		fmt.Fprint(client, "HTTP/1.1 200 Connection established/r/n")
	} else {
		server.Write(b[:n])
	}
	//进行转发
	go io.Copy(server, client)
	io.Copy(client, server)
}
