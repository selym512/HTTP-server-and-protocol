package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"main.go/cmd/internal/headers"
	"main.go/cmd/internal/request"
	"main.go/cmd/internal/response"
	"main.go/cmd/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
	defer server.Close()

	log.Println("Server started on port: ", port)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
func handler(w *response.Writer, req *request.Request) {
	switch {
	case req.RequestLine.RequestTarget == "/yourproblem":
		handler400(w, req)
		return
	case req.RequestLine.RequestTarget == "/myproblem":
		handler500(w, req)
		return
	case strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin"):
		handlerBin(w, req)
		return
	case strings.HasPrefix(req.RequestLine.RequestTarget, "/video"):
		getVideoHandler(w, req)
		return
	default:
		handler200(w, req)
		return
	}
}
func getVideoHandler(w *response.Writer, req *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	file, err := os.ReadFile("/home/myless/Code/go/httpInGo/assets/vim.mp4")
	if err != nil {
		fmt.Print("Err reading video file:", err.Error())
	}
	fmt.Print(file)
	h := response.GetDefaultHeaders(len(file))
	h.Overwrite("Content-Type", "video/mp4")
	err = w.WriteHeaders(h)
	_, err = w.WriteBody(file)
	if err != nil {
		fmt.Print("Err writing video file:", err.Error())
	}

}
func handlerBin(w *response.Writer, req *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	h := response.GetDefaultHeaders(0)
	h.Remove("Content-Length")
	h.Overwrite("Content-Type", "text/plain")
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Trailer", "X-Content-Length, X-Content-SHA256")
	w.WriteHeaders(h)
	t := headers.NewHeaders()
	httpBinProxyURL := fmt.Sprintf("https://httpbin.org%v", strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin"))
	res, err := http.Get(httpBinProxyURL)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to get %v", err.Error()))
		return
	}
	defer res.Body.Close()
	const maxChunkSize = 1024
	buf := make([]byte, maxChunkSize)
	var bodyTotalBytes = 0
	checksumBody := make([]byte, 0, maxChunkSize)
	for {
		n, err := res.Body.Read(buf)
		fmt.Println("Read", n, "bytes")
		if n > 0 {
			bodyBytesWritten, err := w.WriteChunkedBody(buf[:n])
			bodyTotalBytes += n
			fmt.Printf("Adding %v and %x\r\n", bodyBytesWritten, buf[:n])
			checksumBody = append(checksumBody, buf[:n]...)
			if err != nil {
				fmt.Println("error writing chunked body:", err)
				break
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading response body:", err)
			break
		}
	}
	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Println("Error writing chunked body done:", err)
	}
	//fmt.Printf("Entire Body in bytes %x", checksumBody)
	t.Set("X-Content-Length", strconv.Itoa(bodyTotalBytes))
	t.Set("X-Content-SHA256", fmt.Sprintf("%x", sha256.Sum256(checksumBody)))
	err = w.WriteTrailers(t)
	if err != nil {
		fmt.Println(err)
	}
}
func handler200(w *response.Writer, req *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	body := []byte(fmt.Sprint("<html>" +
		"\n <head>" +
		"\n  <title>200 OK</title>" +
		"\n </head>" +
		"\n <body>" +
		"\n  <h1>Success!</h1>" +
		"\n  <p>Your request was an absolute banger.</p>" +
		"\n </body>" +
		"\n</html>",
	))
	h := response.GetDefaultHeaders(len(body))
	h.Overwrite("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}
func handler400(w *response.Writer, req *request.Request) {
	w.WriteStatusLine(response.StatusBadRequest)
	body := []byte(fmt.Sprint("<html>" +
		"\n <head>" +
		"\n  <title>400 Bad Request</title>" +
		"\n </head>" +
		"\n <body>" +
		"\n  <h1>Bad Request</h1>" +
		"\n  <p>Your request honestly kinda sucked.</p>" +
		"\n </body>" +
		"\n</html>",
	))
	h := response.GetDefaultHeaders(len(body))
	h.Overwrite("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}
func handler500(w *response.Writer, req *request.Request) {
	w.WriteStatusLine(response.StatusServerError)
	body := []byte(fmt.Sprint("<html>" +
		"\n <head>" +
		"\n  <title>500 Internal Server Error</title>" +
		"\n </head>" +
		"\n <body>" +
		"\n  <h1>Internal Server Error</h1>" +
		"\n  <p>Okay, you know what? This one is on me.</p>" +
		"\n </body>" +
		"\n</html>",
	))
	h := response.GetDefaultHeaders(len(body))
	h.Overwrite("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}
