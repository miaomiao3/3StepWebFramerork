package web

import (
	"bufio"
	"bytes"
	"github.com/miaomiao3/log"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		method := r.Method
		clientIP := getRealIp(r)

		var copyBody1 io.ReadCloser

		// 复制body用于打log
		buf, _ := ioutil.ReadAll(r.Body)
		copyBody1 = ioutil.NopCloser(bytes.NewBuffer(buf))
		copyBody2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.
		r.Body = copyBody2

		// to catch response
		bufWriter := bufio.NewWriter(w)
		buff := bytes.Buffer{}
		newWriter := &bufferedWriter{w, bufWriter, buff}

		// Call the next handler, which can be another middleware in the chain
		next.ServeHTTP(newWriter, r)

		// skip static path
		latency := time.Now().Sub(start)

		path := r.RequestURI

		reqBody := readBody(copyBody1)

		respBody := newWriter.RespBuffer.Bytes()

		log.Info("%13v| %s | %-7s %s  req: %s, resp: %s",
			latency,
			clientIP,
			method,
			path,
			reqBody,
			string(respBody),
		)

		// You have to flush the buffer at the end
		bufWriter.Flush()

	})
}

type bufferedWriter struct {
	http.ResponseWriter
	out        *bufio.Writer
	RespBuffer bytes.Buffer
}

func (g *bufferedWriter) Write(data []byte) (int, error) {
	g.RespBuffer.Write(data)
	return g.out.Write(data)
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := strings.Trim(buf.String(), " ")
	return s
}

func getRealIp(r *http.Request) string {
	clientIP := r.Header.Get("X-Forwarded-For")
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if len(clientIP) > 0 {
		return clientIP
	}

	if addr := r.Header.Get("X-Appengine-Remote-Addr"); addr != "" {
		return addr
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
func IsStringContainSliceElement(str string, data []string) bool {
	for _, v := range data {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
}
