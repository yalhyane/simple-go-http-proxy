package internal

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	XForwardForHeaderName = "X-Forwarded-For"
	ConnectionHeaderName  = "Connection"
)

var schemeRegex = regexp.MustCompile("^(https?)$")

// Hop-by-hop headers. These are removed when sent to the backend.
// https://0xn3va.gitbook.io/cheat-sheets/web-application/abusing-http-hop-by-hop-request-headers
var hopByHopHeaders = []string{
	"Keep-Alive",
	"Transfer-Encoding",
	"TE",
	"Connection",
	"Trailer",
	"Upgrade",
	"Proxy-Authorization",
	"Proxy-Authenticate",
}

type SimpleHttpProxyConfig struct {
	Addr          string `json:"listen"`
	TargetTimeout time.Duration
}

func (cfg *SimpleHttpProxyConfig) FillDefaults() {
	if cfg.Addr == "" {
		cfg.Addr = "0.0.0.0:8889"
	}
	if cfg.TargetTimeout == 0 {
		cfg.TargetTimeout = 10 * time.Second
	}
}

func (cfg *SimpleHttpProxyConfig) Validate() {
	if err := cfg.ValidateE(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}

// ValidateE This function ensures an Addr address is valid.
func (cfg *SimpleHttpProxyConfig) ValidateE() error {
	if cfg.Addr == "" {
		return errors.New("invalid Addr address")
	}
	return nil
}

// This struct provides an easy way to configure a simple HTTP proxy.
type SimpleHttpProxy struct {
	Config SimpleHttpProxyConfig
}

// ServeHTTP This function is handling requests, validating and forwarding traffic with logging and cleaning up headers.
func (p *SimpleHttpProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println(req.RemoteAddr, "\t\t", req.Method, "\t\t", req.URL, "\t\t Host:", req.Host, "\n\tPath: ", req.URL.Path)
	log.Println("\t\t\t\t\t", req.Header)
	if !isValidScheme(req.URL.Scheme) {
		err := fmt.Sprintf("Invalid scheme: %s\n", req.URL.Scheme)
		log.Println(err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	if req.Method == "GET" && req.URL.Path == "ping" {
		w.WriteHeader(200)

		_, _ = w.Write([]byte("Ping succeeded"))
		return
	}

	c := &http.Client{
		Timeout: p.Config.TargetTimeout,
	}
	req.RequestURI = ""

	p.cleanUpHeaders(req.Header)
	p.handelXForwardHeader(req)

	res, err := c.Do(req)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Printf("Server error: %s\n", err.Error())
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("[Warning] Could not close body: %s\n", err.Error())
		}
	}(res.Body)

	log.Printf("Response from target: %s => status: %s\n", req.RemoteAddr, res.Status)

	p.cleanUpHeaders(res.Header)
	p.copyHeaders(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)
	bw, err := io.Copy(w, res.Body)
	if err != nil {
		log.Printf("Could not forward response body from %s: %s\n", req.RemoteAddr, err.Error())
		return
	} else {
		log.Printf("[Debug] %d bytes writing from %s\n", bw, req.RemoteAddr)
	}

}

func isValidScheme(scheme string) bool {
	log.Printf("Validating scheme: %s\n", scheme)
	return schemeRegex.MatchString(scheme)
}

func (p *SimpleHttpProxy) Start() {
	var (
		tcpListener net.Listener
		err         error
	)

	if tcpListener, err = net.Listen(
		"tcp",
		p.Config.Addr,
	); err != nil {
		log.Fatalf("Error listening on %s: %s\n", p.Config.Addr, err.Error())
	}
	log.Printf("Proxy listening at: %s\n", p.Config.Addr)
	if err := http.Serve(tcpListener, p); err != nil {
		log.Fatalf("Error starting proxy: %s\n", err.Error())
	}
}

// cleanupHeaders removes hop-by-hop headers listed in headers and  in the "Connection" header
func (p *SimpleHttpProxy) cleanUpHeaders(h http.Header) {
	for _, hh := range hopByHopHeaders {
		h.Del(hh)
	}
	if _, ok := h[ConnectionHeaderName]; ok {
		// header not found...
		return
	}
	for _, ch := range h[ConnectionHeaderName] {
		for _, hh := range strings.Split(ch, ",") {
			if nhh := strings.TrimSpace(hh); nhh != "" {
				h.Del(nhh)
			}
		}
	}

	// Parse the Connection header to remove any hop-by-hop headers from the request
	//ch := strings.Split(h.Get(ConnectionHeaderName), ",")
	//for _, hh := range ch {
	//	normalizedHh := strings.TrimSpace(hh)
	//
	//	if strings.EqualFold(strings.TrimSpace(normalizedHh), "close") {
	//		h = http.Header{}
	//		break
	//	}
	//	h.Del(normalizedHh)
	//}

}

func (p *SimpleHttpProxy) handelXForwardHeader(req *http.Request) {
	clientIP, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		log.Printf("Error while parsing HostPort: %s", err.Error())
		return
	}
	host := ""
	if prevHost, ok := req.Header[XForwardForHeaderName]; ok {
		hosts := append(prevHost[:], clientIP)
		host = strings.Join(hosts, ", ")
	}

	req.Header.Set(XForwardForHeaderName, host)

}

func (p *SimpleHttpProxy) copyHeaders(destH, srcH http.Header) {
	for k, h := range srcH {
		for _, v := range h {
			destH.Add(k, v)
		}
	}
}
