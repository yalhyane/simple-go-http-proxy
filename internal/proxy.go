package internal

import (
	"github.com/elazarl/goproxy"
	"log"
	"net"
	"net/http"
)

type CustomProxyConfig struct {
	SimpleHttpProxyConfig
	Verbose bool
}

type CustomProxy struct {
	Config CustomProxyConfig
	proxy  *goproxy.ProxyHttpServer
}

func (p *CustomProxy) Start() {
	var (
		tcpListener net.Listener
		err         error
	)
	p.proxy = goproxy.NewProxyHttpServer()
	p.proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" && req.URL.Path == "/ping" {
			w.WriteHeader(200)
			_, _ = w.Write([]byte("pong"))
			return
		}
		http.Error(w, "This is a proxy server. Does not respond to non-proxy requests.", 500)
	})
	//var isPing goproxy.ReqConditionFunc = func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
	//	log.Println("Checking is PING:", req.Method, req.URL.Path)
	//	return req.Method == "GET" && req.URL.Path == "ping"
	//}
	//var pingHandler goproxy.FuncReqHandler = func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	//	resp := ctx.Resp
	//	resp.StatusCode = 200
	//	w := bytes.NewBufferString("pong")
	//	_ = resp.Write(w)
	//	return req, resp
	//}
	// p.proxy.OnRequest(isPing).Do(pingHandler)
	p.proxy.Verbose = p.Config.Verbose
	if tcpListener, err = net.Listen(
		"tcp",
		p.Config.Addr,
	); err != nil {
		log.Fatalf("Error listening on %s: %s\n", p.Config.Addr, err.Error())
	}
	log.Printf("Proxy listening at: %s\n", p.Config.Addr)
	if err := http.Serve(tcpListener, p.proxy); err != nil {
		log.Fatalf("Error starting proxy: %s\n", err.Error())
	}
}
