package main

import (
	"log"
	"net/http"
	"os"
	"webhook-proxy/wechat"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main()  {
	if len(os.Args) != 2 {
		log.Panic("Please input wechat callback url")
		return
	}

	wechatCallbackUrl := os.Args[1]
	listenPort := ":5001"
	log.Println("Proxy start listen port ", listenPort, "callback url ", wechatCallbackUrl)

	handler := wechat.Handler{CallbackUrl:wechatCallbackUrl}
	mux := http.NewServeMux()
	mux.HandleFunc("/wechat", handler.HandleRequest)
	err := http.ListenAndServe(listenPort, mux)
	if err != nil {
		log.Println("Proxy start failed, error ", err)
	}
}
