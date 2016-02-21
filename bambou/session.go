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
	"encoding/base64"
)

var currentSession *Session

// Represents a user session.
type Session struct {
	Root         Rootable
	Certificate  string
	Username     string
	Password     string
	APIKey       string
	Organization string
	URL          string
}

// Returns a new *Session
//
// You need to provide a Rootable object that will be used to contain
// the results of the authentication process, like the api key for instance.
func NewSession(username, password, organization, url string, root Rootable) *Session {

	return &Session{
		Username:     username,
		Password:     password,
		Organization: organization,
		URL:          url,
		Root:         root,
	}
}

// Returns the current active and authenticated Session.
func CurrentSession() *Session {

	return currentSession
}

// Returns the computed Authorization HTTP header.
func (s *Session) MakeAuthorizationHeaders() string {

	if s.Username == "" {
		panic("Error while creating headers: User must be set")
	}

	if s.Password == "" && s.APIKey == "" {
		panic("Error while creating headers: Password or APIKey must be set")
	}

	var key string

	if s.APIKey != "" {
		key = s.APIKey
	} else {
		key = s.Password
	}

	return "XREST " + base64.StdEncoding.EncodeToString([]byte(s.Username+":"+key))
}

// Starts the session.
//
// At that point the authentication will be done.
func (s *Session) Start() *Error {

	if currentSession != nil {
		return nil
	}

	currentSession = s

	err := s.Root.Fetch()

	if err != nil {
		Logger().Errorf("Error during Authentication: %s", err.Error())
		return err
	}

	s.APIKey = s.Root.GetAPIKey()

	return nil
}

// Resets the session.
func (s *Session) Reset() {

	if currentSession == nil {
		return
	}

	s.APIKey = ""
	s.Root.SetAPIKey("")

	currentSession = nil
}
