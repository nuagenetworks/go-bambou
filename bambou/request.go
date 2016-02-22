// Copyright (c) 2015, Alcatel-Lucent Inc.
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of bambou nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package bambou

import (
	"bytes"
	"net/http"
)

const (
	requestMethodDelete = "DELETE"
	requestMethodGet    = "GET"
	requestMethodHead   = "HEAD"
	requestMethodPost   = "POST"
	requestMethodPut    = "PUT"
)

// Represents a request.
type request struct {
	Data       []byte
	Headers    map[string]string
	Identifier string
	Parameters map[string]string
	URL        string
	Method     string
}

// Returns a new  *request.
func newRequest(url string) *request {

	return &request{
		URL:        url,
		Headers:    make(map[string]string),
		Parameters: make(map[string]string),
		Method:     requestMethodGet,
	}
}

// Sets the value of Header field.
func (r *request) setHeader(name, value string) {

	r.Headers[name] = value
}

// Returns the value of Header field.
func (r *request) getHeader(name string) string {

	return r.Headers[name]
}

// Sets the value of a query parameter.
func (r *request) setParameter(name, value string) {

	r.Parameters[name] = value
}

// Gets the value of a query parameter.
func (r *request) getParameter(name string) string {

	return r.Parameters[name]
}

// Returns a native http.Request.
func (r *request) toNative() *http.Request {

	req, _ := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.Data))

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	return req
}
