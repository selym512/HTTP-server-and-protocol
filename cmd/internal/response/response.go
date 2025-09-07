package response

import (
	"fmt"
	"io"
	"main.go/cmd/internal/headers"
	"net"
)

type Writer struct {
	writer      io.Writer
	writerState WriterState
}
type WriterState int

const (
	Initiated WriterState = iota
	StatusLineWritten
	HeadersWritten
	BodyWritten
	TrailersWritten
)

func GetResponseWriter(w *net.Conn) *Writer {
	return &Writer{writer: *w, writerState: Initiated}
}
func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != HeadersWritten {
		return 0, fmt.Errorf("ERROR: Response body should be written after the headers and status line - %d", w.writerState)
	}

	chunkSize := len(p)
	nTotal := 0

	n, err := fmt.Fprintf(w.writer, "%X\r\n", chunkSize)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write(p)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return nTotal, err
	}
	nTotal += n
	return nTotal, err
}
func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != HeadersWritten {
		return 0, fmt.Errorf("ERROR: Response body should be written after the headers and status line - %d", w.writerState)
	}
	n, err := w.writer.Write([]byte("0\r\n"))
	if err != nil {
		return n, err
	}
	w.writerState = BodyWritten
	return n, err
}
func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != Initiated {
		return fmt.Errorf("ERROR: Status Line Should be written first")
	}
	_, err := w.writer.Write(getStatusLine(statusCode))
	if err == nil {
		w.writerState = StatusLineWritten
	}
	return err
}
func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != StatusLineWritten {
		return fmt.Errorf("ERROR: Headers should be written after the Status Line and before the response body")
	}
	for k, v := range headers {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	if err == nil {
		w.writerState = HeadersWritten
	}
	return err
}
func (w *Writer) WriteTrailers(t headers.Headers) error {
	if w.writerState != BodyWritten {
		return fmt.Errorf("Error: Trailer Header can only be written after the response body")
	}
	for k, v := range t {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	if err == nil {
		w.writerState = TrailersWritten
	}
	return err
}
func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != HeadersWritten {
		return 0, fmt.Errorf("ERROR: Response body should be written after the headers and status line")
	}
	int, err := w.writer.Write(p)
	w.writerState = BodyWritten
	return int, err
}
func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}
