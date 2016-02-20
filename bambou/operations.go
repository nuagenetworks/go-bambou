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
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// Interface for Operationables objects.
type Operationable interface {
	Fetch() *Error
	Save() *Error
	Delete() *Error
}

// Children fecthing information.
//
// This structure will be used to pass and get back information
// during the fetching of some children.
type FetchingInfo struct {
	Filter     string
	FilterType string
	GroupBy    []string
	OrderBy    string
	Page       int
	PageSize   int
	TotalCount int
}

// Returns a new *FetchingInfo
func NewFetchingInfo() *FetchingInfo {

	return &FetchingInfo{
		Page:     -1,
		PageSize: -1,
	}
}

// String representation of the FetchingInfo.
func (f *FetchingInfo) String() string {

	return fmt.Sprintf("{FetchingInfo: page: %d, pagesize: %d, totalcount: %d}", f.Page, f.PageSize, f.TotalCount)
}

// Private.
//
// Fills the HTTP headers of the given Request according to the given FetchingInfo.
func prepareHeaders(request *Request, info *FetchingInfo) {

	request.SetHeader("X-Nuage-PageSize", "50")

	if info == nil {
		return
	}

	if info.Filter != "" {
		request.SetHeader("X-Nuage-Filter", info.Filter)
	}

	if info.OrderBy != "" {
		request.SetHeader("X-Nuage-OrderBy", info.OrderBy)
	}

	if info.Page != -1 {
		request.SetHeader("X-Nuage-Page", strconv.Itoa(info.Page))
	}

	if info.PageSize > 0 {
		request.SetHeader("X-Nuage-PageSize", strconv.Itoa(info.PageSize))
	}

	if len(info.GroupBy) > 0 {
		request.SetHeader("X-Nuage-GroupBy", "true")
		request.SetHeader("X-Nuage-Attributes", strings.Join(info.GroupBy, ", "))
	}
}

// Private.
//
// Fills the given FetchingInfo according to the HTTP headers of the given Response.
func readHeaders(response *Response, info *FetchingInfo) {

	if info == nil {
		return
	}

	info.Filter = response.GetHeader("X-Nuage-Filter")
	info.FilterType = response.GetHeader("X-Nuage-FilterType")
	// info.GroupBy = response.GetHeader("X-Nuage-GroupBy")
	info.OrderBy = response.GetHeader("X-Nuage-OrderBy")
	info.Page, _ = strconv.Atoi(response.GetHeader("X-Nuage-Page"))
	info.PageSize, _ = strconv.Atoi(response.GetHeader("X-Nuage-PageSize"))
	info.TotalCount, _ = strconv.Atoi(response.GetHeader("X-Nuage-Count"))
}

// Fetchs the given Exposable from the server.
func FetchEntity(object Exposable) *Error {

	request := NewRequest(object.GetURL())
	connection := NewConnection()
	response, error := connection.Start(request)

	if error != nil {
		Logger().Errorf("Error during FetchEntity: %s", error.Error())
		return error
	}

	err := json.Unmarshal(response.Data[1:len(response.Data)-1], &object)

	if err != io.EOF && err != nil {
		panic("Unable to unmarshal json: " + err.Error())
	}

	return nil
}

// Saves the given Exposable into the server.
func SaveEntity(object Exposable) *Error {

	data, err := json.Marshal(object)

	if err != nil {
		panic("Unable to marshal json: " + err.Error())
	}

	request := NewRequest(object.GetURL())
	request.Method = RequestMethodPut
	request.Data = data

	connection := NewConnection()
	_, error := connection.Start(request)

	if error != nil {
		Logger().Errorf("Error during SaveEntity: %s", error.Error())
		return error
	}

	return nil
}

// Deletes the given Exposable from the server.
func DeleteEntity(object Exposable) *Error {

	request := NewRequest(object.GetURL())
	request.Method = RequestMethodDelete

	connection := NewConnection()
	_, error := connection.Start(request)

	if error != nil {
		Logger().Errorf("Error during DeleteEntity: %s", error.Error())
		return error
	}

	return nil
}

// Fetches the children with of given parent identified by the given Identify.
//
// The dest parameters must be a pointer to some Exposable object.
// The given FetchingInfo will be used to apply pagination, or filtering etc, and will
// be populated back according to the response (with the total cound of objects for instance).
// In case of error, an *Error is returned, otherwise nil.
func FetchChildren(parent Exposable, identity RESTIdentity, dest interface{}, info *FetchingInfo) *Error {

	request := NewRequest(parent.GetURLForChildrenIdentity(identity))
	prepareHeaders(request, info)

	connection := NewConnection()
	response, error := connection.Start(request)

	if error != nil {
		Logger().Errorf("Error during FetchChildren: %s", error.Error())
		return error
	}

	err := json.Unmarshal(response.Data, dest)

	if err != io.EOF && err != nil {
		panic("Unable to unmarshal json: " + err.Error())
	}

	readHeaders(response, info)

	Identify(dest, identity)

	return nil
}

// Creates a new child Exposable under the given parent Exposable in the server.
//
// In case of error, an *Error is returned, otherwise nil.
func CreateChild(parent Exposable, child Exposable) *Error {

	data, err := json.Marshal(child)

	if err != nil {
		panic("Unable to marshal json: " + err.Error())
	}

	request := NewRequest(parent.GetURLForChildrenIdentity(child.GetIdentity()))
	request.Method = RequestMethodPost
	request.Data = data

	connection := NewConnection()
	response, error := connection.Start(request)

	if error != nil {
		Logger().Errorf("Error during CreateChild: %s", error.Error())
		return error
	}

	err = json.Unmarshal(response.Data[1:len(response.Data)-1], child)

	if err != io.EOF && err != nil {
		panic("Unable to unmarshal json: " + err.Error())
	}

	return nil
}

// Assign the list of given child Exposables to the given Exposable parent in the server.
//
// You must provide the Identity of the children you want to assign. This is mandatory in
// case you want to unassign all objects.
// In case of error, an Error is returned, otherwise nil.
func AssignChildren(parent Exposable, children interface{}, identity RESTIdentity) *Error {

	ids := make([]string, 0)

	if children != nil {
		l := reflect.ValueOf(children).Len()

		for i := 0; i < l; i++ {

			o := reflect.ValueOf(children).Index(i).Elem()

			identityField := o.FieldByName("ID")
			ids = append(ids, identityField.String())
		}
	}

	data, err := json.Marshal(&ids)

	if err != nil {
		panic("Unable to marshal json: " + err.Error())
	}

	request := NewRequest(parent.GetURLForChildrenIdentity(identity))
	request.Method = RequestMethodPut
	request.Data = data

	connection := NewConnection()
	_, error := connection.Start(request)

	if error != nil {
		Logger().Errorf("Error during AssignChildren: %s", error.Error())
		return error
	}

	return nil
}
