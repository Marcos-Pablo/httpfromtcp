package response

import (
	"fmt"
	"io"

	"github.com/Marcos-Pablo/httpfromtcp/internal/headers"
)

type WriterState int

const (
	WriterStateStatusLine WriterState = iota
	WriterStateHeaders
	WriterStateBody
	WriterStateDone
)

type Writer struct {
	State WriterState
	buff  io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		State: WriterStateStatusLine,
		buff:  w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.State != WriterStateStatusLine {
		return fmt.Errorf("unable to write status line")
	}
	var reasonPhrase string
	switch statusCode {
	case StatusOK:
		reasonPhrase = "HTTP/1.1 200 OK"
	case StatusBadRequest:
		reasonPhrase = "HTTP/1.1 400 Bad Request"
	case StatusInternalServerError:
		reasonPhrase = "HTTP/1.1 500 Internal Server Error"
	}

	_, err := w.buff.Write([]byte(reasonPhrase + crlf))

	if err != nil {
		return fmt.Errorf("Error writing status line: %s", err)
	}

	w.State = WriterStateHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.State != WriterStateHeaders {
		return fmt.Errorf("unable to write headers")
	}
	for k, v := range headers {
		_, err := fmt.Fprintf(w.buff, "%s: %s%s", k, v, crlf)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w.buff, crlf)
	if err != nil {
		return err
	}
	w.State = WriterStateBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.State != WriterStateBody {
		return 0, fmt.Errorf("unable to write body")
	}
	n, err := w.buff.Write(p)
	w.State = WriterStateDone
	return n, err
}
