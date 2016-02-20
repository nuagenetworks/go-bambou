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

const (
	ResponseCodeZero                  = 0
	ResponseCodeSuccess               = 200
	ResponseCodeCreated               = 201
	ResponseCodeEmpty                 = 204
	ResponseCodeMultipleChoices       = 300
	ResponseCodeBadRequest            = 400
	ResponseCodeUnauthorized          = 401
	ResponseCodePermissionDenied      = 403
	ResponseCodeNotFound              = 404
	ResponseCodeMethodNotAllowed      = 405
	ResponseCodeConnectionTimeout     = 408
	ResponseCodeConflict              = 409
	ResponseCodePreconditionFailed    = 412
	ResponseCodeAuthenticationExpired = 419
	ResponseCodeInternalServerError   = 500
	ResponseCodeServiceUnavailable    = 503
)

// Represents a Response.
type Response struct {
	Code    int
	Data    []byte
	Headers map[string]string
}

// Returns a new *Response.
func NewResponse() *Response {

	return &Response{
		Headers: make(map[string]string),
	}
}

// Sets the value of Header field.
func (r *Response) SetHeader(name, value string) {

	r.Headers[name] = value
}

// Returns the value of Header field.
func (r *Response) GetHeader(name string) string {

	return r.Headers[name]
}
