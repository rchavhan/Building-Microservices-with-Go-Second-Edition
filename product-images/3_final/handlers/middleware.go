package handlers

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-hclog"
)

type Middleware struct {
	maxSize int64
	logger  hclog.Logger
}

// NewMiddleware creates a new middleware handlers
func NewMiddleware(maxContentLength int64, logger hclog.Logger) *Middleware {
	return &Middleware{maxSize: maxContentLength, logger: logger}
}

// CheckContentLengthMiddleware ensures that the content length is not greater than
// our allowed max.
// This can not be 100% depended on as Content-Length might be reported incorrectly
// however it is a fast first pass check.
func (mw *Middleware) CheckContentLengthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// if the content length is greater than the max reject the request
		if r.ContentLength > mw.maxSize {
			http.Error(rw, "Unable to save file, content too large", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(rw, r)
	})
}

// GZipResponseMiddleware detects if the client can handle
// zipped content and if so returns the response in GZipped format
func (mw *Middleware) GZipResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// if client cant handle gzip send plain
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			//f.log.Debug("Unable to handle gzipped", "file", fp)
			fmt.Println("Not gzip")
			next.ServeHTTP(rw, r)
			return
		}

		// client can handle gziped content send gzipped to speed up transfer
		// set the content encoding header for gzip
		rw.Header().Add("Content-Encoding", "gzip")

		// file server sets the content stream
		// nice
		//rw.Header().Add("Content-Type", "application/octet-stream")

		wr := NewWrappedResponseWriter(rw)
		defer wr.Flush()

		// write the file
		next.ServeHTTP(wr, r)
	})
}

// WrappedResponseWriter wrapps the default http.ResponseWriter in a GZip stream
type WrappedResponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

// NewWrappedResponseWriter returns a new wrapped response writer
func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedResponseWriter {
	// wrap the default writer in a gzip writer
	gw := gzip.NewWriter(rw)

	return &WrappedResponseWriter{rw, gw}
}

// Header implements the http.ResponseWriter Header method
func (wr *WrappedResponseWriter) Header() http.Header {
	return wr.rw.Header()
}

// Write implements the http.ResponseWriter Write method
func (wr *WrappedResponseWriter) Write(d []byte) (int, error) {
	return wr.gw.Write(d)
}

// WriteHeader implements the http.ResponseWriter WriteHeader method
func (wr *WrappedResponseWriter) WriteHeader(statusCode int) {
	wr.rw.WriteHeader(statusCode)
}

// Flush implements the http.ResponseWriter Flush method
func (wr *WrappedResponseWriter) Flush() {
	// flush and close the writer
	wr.gw.Flush()
	wr.gw.Close()
}
