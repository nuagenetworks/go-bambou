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
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var _currentSession Operationable

// CurrentSession returns the current active and authenticated Session.
func CurrentSession() Operationable {

	return _currentSession
}

// Operationable is the the interface that any kind of session must implement
type Operationable interface {
	Start() *Error
	Reset()
	FetchEntity(Exposable) *Error
	SaveEntity(Exposable) *Error
	DeleteEntity(Exposable) *Error
	FetchChildren(Exposable, Identity, interface{}, *FetchingInfo) *Error
	CreateChild(Exposable, Exposable) *Error
	AssignChildren(Exposable, interface{}, Identity) *Error
	NextEvent(NotificationsChannel, *string)
	Root() Rootable
}

// Session represents a user session. It provides the entire
// communication layer with the backend. It must implement the Operationable interface.
type Session struct {
	root         Rootable
	Certificate  string
	Username     string
	Password     string
	Organization string
	URL          string
	client       *http.Client
}

// NewSession returns a new *Session
// You need to provide a Rootable object that will be used to contain
// the results of the authentication process, like the api key for instance.
func NewSession(username, password, organization, url string, root Rootable) *Session {

	return &Session{
		Username:     username,
		Password:     password,
		Organization: organization,
		URL:          url,
		root:         root,
		client:       &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}},
	}
}

func (s *Session) makeAuthorizationHeaders() string {

	if s.Username == "" {
		panic("Error while creating headers: User must be set")
	}

	key := s.root.APIKey()
	if s.Password == "" && key == "" {
		panic("Error while creating headers: Password or APIKey must be set")
	}

	if key == "" {
		key = s.Password
	}

	return "XREST " + base64.StdEncoding.EncodeToString([]byte(s.Username+":"+key))
}

func (s *Session) prepareHeaders(request *http.Request, info *FetchingInfo) {

	request.Header.Set("X-Nuage-PageSize", "50")
	request.Header.Set("X-Nuage-Organization", s.Organization)
	request.Header.Set("Authorization", s.makeAuthorizationHeaders())
	request.Header.Set("Content-Type", "application/json")

	if info == nil {
		return
	}

	if info.Filter != "" {
		request.Header.Set("X-Nuage-Filter", info.Filter)
	}

	if info.OrderBy != "" {
		request.Header.Set("X-Nuage-OrderBy", info.OrderBy)
	}

	if info.Page != -1 {
		request.Header.Set("X-Nuage-Page", strconv.Itoa(info.Page))
	}

	if info.PageSize > 0 {
		request.Header.Set("X-Nuage-PageSize", strconv.Itoa(info.PageSize))
	}

	if len(info.GroupBy) > 0 {
		request.Header.Set("X-Nuage-GroupBy", "true")
		request.Header.Set("X-Nuage-Attributes", strings.Join(info.GroupBy, ", "))
	}
}

func (s *Session) readHeaders(response *http.Response, info *FetchingInfo) {

	if info == nil {
		return
	}

	info.Filter = response.Header.Get("X-Nuage-Filter")
	info.FilterType = response.Header.Get("X-Nuage-FilterType")
	info.OrderBy = response.Header.Get("X-Nuage-OrderBy")
	info.Page, _ = strconv.Atoi(response.Header.Get("X-Nuage-Page"))
	info.PageSize, _ = strconv.Atoi(response.Header.Get("X-Nuage-PageSize"))
	info.TotalCount, _ = strconv.Atoi(response.Header.Get("X-Nuage-Count"))

	// info.GroupBy = response.Header.Get("X-Nuage-GroupBy")
}

func (s *Session) send(request *http.Request, info *FetchingInfo) (*http.Response, *Error) {

	s.prepareHeaders(request, info)

	response, err := s.client.Do(request)

	if err != nil {
		return response, NewError(0, err.Error())
	}

	switch response.StatusCode {

	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		s.readHeaders(response, info)
		return response, nil

	case http.StatusMultipleChoices:
		newURL := request.URL.String() + "?responseChoice=1"
		request, _ = http.NewRequest(request.Method, newURL, request.Body)
		return s.send(request, info)

	case http.StatusConflict:
		data, _ := ioutil.ReadAll(response.Body)
		error := NewError(response.StatusCode, string(data))
		json.Unmarshal(data, error)
		return nil, error

	default:
		return nil, NewError(response.StatusCode, response.Status)
	}
}

func (s *Session) getGeneralURL(o Exposable) string {

	if o.Identity().ResourceName == "" {
		panic("Cannot GetGeneralURL of that as no ResourceName in its Identity")
	}

	return s.URL + "/" + o.Identity().ResourceName
}

func (s *Session) getPersonalURL(o Exposable) string {

	if _, ok := o.(Rootable); ok {
		return s.URL + "/" + o.Identity().RESTName
	}

	if o.Identifier() == "" {
		panic("Cannot GetPersonalURL of an object with no ID set")
	}

	return s.getGeneralURL(o) + "/" + o.Identifier()
}

func (s *Session) getURLForChildrenIdentity(o Exposable, childrenIdentity Identity) string {

	if _, ok := o.(Rootable); ok {
		return s.URL + "/" + childrenIdentity.ResourceName
	}

	return s.getPersonalURL(o) + "/" + childrenIdentity.ResourceName
}

// Root returns the Root API object.
func (s *Session) Root() Rootable {

	return s.root
}

// Start starts the session.
// At that point the authentication will be done.
func (s *Session) Start() *Error {

	_currentSession = s

	err := s.FetchEntity(s.root)

	if err != nil {
		Logger().Errorf("Error during Authentication: %s", err.Error())
		return err
	}

	return nil
}

// Reset resets the session.
func (s *Session) Reset() {

	s.root.SetAPIKey("")

	_currentSession = nil
}

// FetchEntity fetchs the given Exposable from the server.
// You should not use this function by yourself.
func (s *Session) FetchEntity(object Exposable) *Error {

	request, _ := http.NewRequest("GET", s.getPersonalURL(object), nil)
	response, err1 := s.send(request, nil)

	if response != nil {
		defer response.Body.Close()
	}

	if err1 != nil {
		Logger().Errorf("Error during FetchEntity: %s", err1.Error())
		return err1
	}

	data, _ := ioutil.ReadAll(response.Body)
	data = data[1 : len(data)-1]
	err2 := json.Unmarshal(data, &object)

	if err2 != io.EOF && err2 != nil {
		panic(fmt.Sprintf("Unable to unmarshal json %s: %s", string(data), err2.Error()))
	}

	return nil
}

// SaveEntity saves the given Exposable into the server.
// You should not use this function by yourself.
func (s *Session) SaveEntity(object Exposable) *Error {

	data, _ := json.Marshal(object)
	request, _ := http.NewRequest("PUT", s.getPersonalURL(object), bytes.NewBuffer(data))
	response, err1 := s.send(request, nil)

	if response != nil {
		defer response.Body.Close()
	}

	if err1 != nil {
		Logger().Errorf("Error during SaveEntity: %s", err1.Error())
		return err1
	}

	data, _ = ioutil.ReadAll(response.Body)

	if len(data) > 0 {
		data := data[1 : len(data)-1]
		err2 := json.Unmarshal(data, &object)

		if err2 != io.EOF && err2 != nil {
			panic(fmt.Sprintf("Unable to unmarshal json %s: %s", string(data), err2.Error()))
		}
	}

	return nil
}

// DeleteEntity deletes the given Exposable from the server.
// You should not use this function by yourself.
func (s *Session) DeleteEntity(object Exposable) *Error {

	request, _ := http.NewRequest("DELETE", s.getPersonalURL(object), nil)
	_, error := s.send(request, nil)

	if error != nil {
		Logger().Errorf("Error during DeleteEntity: %s", error.Error())
		return error
	}

	return nil
}

// FetchChildren fetches the children with of given parent identified by the given identify.
// You should not use this function by yourself.
func (s *Session) FetchChildren(parent Exposable, identity Identity, dest interface{}, info *FetchingInfo) *Error {

	request, _ := http.NewRequest("GET", s.getURLForChildrenIdentity(parent, identity), nil)
	response, err1 := s.send(request, nil)

	if response != nil {
		defer response.Body.Close()
	}

	if err1 != nil {
		Logger().Errorf("Error during FetchChildren: %s", err1.Error())
		return err1
	}

	if response.StatusCode == http.StatusNoContent || response.ContentLength == 0 {
		return nil
	}

	data, _ := ioutil.ReadAll(response.Body)
	err2 := json.Unmarshal(data, dest)

	if err2 != io.EOF && err2 != nil {
		panic(fmt.Sprintf("Unable to unmarshal json %s: %s", string(data), err2.Error()))
	}

	identify(dest, identity)

	return nil
}

// CreateChild creates a new child Exposable under the given parent Exposable in the server.
// You should not use this function by yourself.
func (s *Session) CreateChild(parent Exposable, child Exposable) *Error {

	data, _ := json.Marshal(child)
	request, _ := http.NewRequest("POST", s.getURLForChildrenIdentity(parent, child.Identity()), bytes.NewBuffer(data))
	response, err1 := s.send(request, nil)

	if response != nil {
		defer response.Body.Close()
	}

	if err1 != nil {
		Logger().Errorf("Error during CreateChild: %s", err1.Error())
		return err1
	}

	data, _ = ioutil.ReadAll(response.Body)
	data = data[1 : len(data)-1]
	err2 := json.Unmarshal(data, child)

	if err2 != io.EOF && err2 != nil {
		panic(fmt.Sprintf("Unable to unmarshal json %s: %s", string(data), err2.Error()))
	}

	return nil
}

// AssignChildren assigns the list of given child Exposables to the given Exposable parent in the server.
// You should not use this function by yourself.
func (s *Session) AssignChildren(parent Exposable, children interface{}, identity Identity) *Error {

	var ids []string

	if children != nil {
		l := reflect.ValueOf(children).Len()
		for i := 0; i < l; i++ {
			o := reflect.ValueOf(children).Index(i).Elem()
			identityField := o.FieldByName("ID")
			ids = append(ids, identityField.String())
		}
	}

	data, _ := json.Marshal(&ids)
	request, _ := http.NewRequest("PUT", s.getURLForChildrenIdentity(parent, identity), bytes.NewBuffer(data))
	_, error := s.send(request, nil)

	if error != nil {
		Logger().Errorf("Error during AssignChildren: %s", error.Error())
		return error
	}

	return nil
}

// NextEvent will return the next notification from the backend as it occurs and will
// send it to the correct channel.
func (s *Session) NextEvent(channel NotificationsChannel, lastEventID *string) {

	currentURL := s.URL + "/events"

	if *lastEventID != "" {
		currentURL += "?uuid=" + *lastEventID
	}

	request, _ := http.NewRequest("GET", currentURL, nil)
	response, err1 := s.send(request, nil)

	if err1 != nil {
		Logger().Errorf("Error during push: %s", err1.Error())
		return
	}

	data, _ := ioutil.ReadAll(response.Body)
	notification := NewNotification()
	err2 := json.Unmarshal(data, &notification)

	if err2 != io.EOF && err2 != nil {
		Logger().Errorf("Error during push: %s", err2.Error())
		return
	}

	*lastEventID = notification.UUID

	if len(notification.Events) > 0 {
		channel <- notification
	}
}
