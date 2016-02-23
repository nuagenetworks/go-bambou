// Copyright (c) 2015, Alcatel-Lucent Inc.
// All rights reserved.
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
// * Neither the name of bambou nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
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
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// connection epresents a connection to a remote server.
// It is wrapper over the standard net/http
// client. You should never have to use it manually.
type connection struct {
	HasTimeouted       bool
	Timeout            time.Duration
	UsesAuthentication bool
	UserInfo           interface{}
}

var sendNativeRequest = func(request *request) *response {

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: transport,
	}

	nativeResponse, err := client.Do(request.toNative())

	if nativeResponse != nil {
		defer nativeResponse.Body.Close()
	}

	if err != nil {
		panic("Error while performing the request: " + err.Error())
	}

	response := newResponse()
	response.Code = nativeResponse.StatusCode

	for h, v := range nativeResponse.Header {
		response.setHeader(h, strings.Join(v, ", "))
	}

	response.Data, _ = ioutil.ReadAll(nativeResponse.Body)

	return response
}

// Returns a pointer to a new *Connection.
func newConnection() *connection {

	return &connection{
		Timeout: time.Duration(60) * time.Second,
	}
}

// Starts the connection with the given request.
// If the request is a success, then a response will be returned and the error will be nil.
// In case of error, the error will be returned, and the response will be nil.
func (c *connection) start(request *request) (*response, *Error) {
	session := CurrentSession()
	request.setHeader("X-Nuage-Organization", session.Organization)
	request.setHeader("Authorization", session.makeAuthorizationHeaders())
	request.setHeader("Content-Type", "application/json")

	logger := Logger()
	logger.Infof("Req : %s %s %s", request.Method, request.URL, request.Parameters)
	logger.Debugf("Req : Headers: %s", request.Headers)
	logger.Debugf("Req : Data: %s", request.Data)

	response := sendNativeRequest(request)

	defer func() {
		logger.Debugf("Resp: %d %s %s %s", response.Code, request.Method, request.URL, request.Parameters)
		logger.Debugf("Resp: Headers: %s", response.Headers)
		logger.Debugf("Resp: Data: %s", response.Data)
	}()

	switch response.Code {

	case responseCodeSuccess, responseCodeCreated, responseCodeEmpty:
		return response, nil

	case responseCodeMultipleChoices:
		request.URL += "?responseChoice=1"
		return c.start(request)

	case responseCodeConflict:
		error := NewError(response.Code, string(response.Data))
		json.Unmarshal(response.Data, &error)
		return nil, error

	case responseCodeAuthenticationExpired:
		CurrentSession().Reset()
		CurrentSession().Start()
		return c.start(request)

	case responseCodeBadrequest:
		return nil, NewError(response.Code, "Bad request.")

	case responseCodeUnauthorized:
		return nil, NewError(response.Code, "Unauthorized.")

	case responseCodePermissionDenied:
		return nil, NewError(response.Code, "Permission denied.")

	case responseCodeMethodNotAllowed:
		return nil, NewError(response.Code, "Not allowed.")

	case responseCodeNotFound:
		return nil, NewError(response.Code, "Not found.")

	case responseCodeConnectionTimeout:
		return nil, NewError(response.Code, "Timeout.")

	case responseCodePreconditionFailed:
		return nil, NewError(response.Code, "Precondition failed.")

	case responseCodeInternalServerError:
		return nil, NewError(response.Code, "Internal server error.")

	case responseCodeServiceUnavailable:
		return nil, NewError(response.Code, "Service unavailable.")

	default:
		return nil, NewError(response.Code, "Unknown error.")
	}
}
