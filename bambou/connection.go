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
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pborman/uuid"
)

// Represents a Connection.
//
// It is wrapper over the standard net/http
// client. You should never have to use it manually.
type Connection struct {
	HasTimeouted       bool
	Timeout            time.Duration
	TransactionID      string
	UsesAuthentication bool
	UserInfo           interface{}
}

// Returns a pointer to a new *Connection.
func NewConnection() *Connection {

	return &Connection{
		Timeout:       time.Duration(60) * time.Millisecond,
		TransactionID: uuid.New(),
	}
}

// Starts the connection with the given Request.
//
// If the request is a success, then a response will be returned and the error will be nil.
// In case of error, the error will be returned, and the response will be nil.

func (c *Connection) Start(request *Request) (*Response, *Error) {
	session := CurrentSession()
	request.SetHeader("X-Nuage-Organization", session.Organization)
	request.SetHeader("Authorization", session.MakeAuthorizationHeaders())
	request.SetHeader("Content-Type", "application/json")

	logger := Logger()
	logger.Infof("Req : %s %s %s (id: %s)", request.Method, request.URL, request.Parameters, c.TransactionID)
	logger.Debugf("Req : Headers: %s", request.Headers)
	logger.Debugf("Req : Data: %s", request.Data)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: transport,
	}

	nativeResponse, err := client.Do(request.ToNative())

	if nativeResponse != nil {
		defer nativeResponse.Body.Close()
	}

	if err != nil {
		panic("Error while performing the request: " + err.Error())
	}

	response := NewResponse()
	response.Code = nativeResponse.StatusCode

	for h, v := range nativeResponse.Header {
		response.SetHeader(h, strings.Join(v, ", "))
	}

	defer func() {
		logger.Debugf("Resp: %d (id: %s)", response.Code, c.TransactionID)
		logger.Debugf("Resp: Headers: %s", response.Headers)
		logger.Debugf("Resp: Data: %s", response.Data)
	}()

	switch response.Code {
	case ResponseCodeMultipleChoices:
		request.URL += "?responseChoice=1"
		return c.Start(request)

	case ResponseCodeSuccess, ResponseCodeCreated:
		data, err := ioutil.ReadAll(nativeResponse.Body)

		if err != nil {
			panic("Error while decoding the body data: " + err.Error())
		}

		response.Data = data

		return response, nil

	case ResponseCodeEmpty:
		return response, nil

	case ResponseCodeConflict:
		data, err := ioutil.ReadAll(nativeResponse.Body)

		if err != nil {
			panic("Error while decoding the body data: " + err.Error())
		}

		error := NewError(response.Code, string(data))
		json.Unmarshal(data, &error)

		return nil, error

	case ResponseCodeAuthenticationExpired:
		CurrentSession().Reset()
		CurrentSession().Start()
		return c.Start(request)

	case ResponseCodeBadRequest:
		return nil, NewError(response.Code, "Bad request.")

	case ResponseCodeUnauthorized:
		return nil, NewError(response.Code, "Unauthorized.")

	case ResponseCodePermissionDenied:
		return nil, NewError(response.Code, "Permission denied.")

	case ResponseCodeMethodNotAllowed:
		return nil, NewError(response.Code, "Not allowed.")

	case ResponseCodeNotFound:
		return nil, NewError(response.Code, "Not found.")

	case ResponseCodeConnectionTimeout:
		return nil, NewError(response.Code, "Timeout.")

	case ResponseCodePreconditionFailed:
		return nil, NewError(response.Code, "Precondition failed.")

	case ResponseCodeInternalServerError:
		return nil, NewError(response.Code, "Internal server error.")

	case ResponseCodeServiceUnavailable:
		return nil, NewError(response.Code, "Service unavailable.")

	default:
		return nil, NewError(response.Code, "Unknown error.")
	}
}
