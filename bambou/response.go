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
	responseCodeZero                  = 0
	responseCodeSuccess               = 200
	responseCodeCreated               = 201
	responseCodeEmpty                 = 204
	responseCodeMultipleChoices       = 300
	responseCodeBadrequest            = 400
	responseCodeUnauthorized          = 401
	responseCodePermissionDenied      = 403
	responseCodeNotFound              = 404
	responseCodeMethodNotAllowed      = 405
	responseCodeConnectionTimeout     = 408
	responseCodeConflict              = 409
	responseCodePreconditionFailed    = 412
	responseCodeAuthenticationExpired = 419
	responseCodeInternalServerError   = 500
	responseCodeServiceUnavailable    = 503
)

// Represents a response.
type response struct {
	Code    int
	Data    []byte
	Headers map[string]string
}

// Returns a new *response.
func newResponse() *response {

	return &response{
		Headers: make(map[string]string),
	}
}

// Sets the value of Header field.
func (r *response) setHeader(name, value string) {

	r.Headers[name] = value
}

// Returns the value of Header field.
func (r *response) getHeader(name string) string {

	return r.Headers[name]
}
