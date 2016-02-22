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
	"io"
	"reflect"
)

// Interface for Operationables objects.
type Operationable interface {
	Fetch() *Error
	Save() *Error
	Delete() *Error
	GetPersonalURL() string
	GetGeneralURL() string
	GetURLForChildrenIdentity(Identity) string
}

// Fetchs the given Exposable from the server.
func FetchEntity(object Operationable) *Error {

	request := NewRequest(object.GetPersonalURL())
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

	data, _ := json.Marshal(object)

	request := NewRequest(object.GetPersonalURL())
	request.Method = RequestMethodPut
	request.Data = data

	connection := NewConnection()
	response, err1 := connection.Start(request)

	if err1 != nil {
		Logger().Errorf("Error during SaveEntity: %s", err1.Error())
		return err1
	}

	err2 := json.Unmarshal(response.Data[1:len(response.Data)-1], &object)

	if err2 != io.EOF && err2 != nil {
		panic("Unable to unmarshal json: " + err2.Error())
	}

	return nil
}

// Deletes the given Exposable from the server.
func DeleteEntity(object Exposable) *Error {

	request := NewRequest(object.GetPersonalURL())
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
func FetchChildren(parent Exposable, identity Identity, dest interface{}, info *FetchingInfo) *Error {

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

	data, _ := json.Marshal(child)

	request := NewRequest(parent.GetURLForChildrenIdentity(child.GetIdentity()))
	request.Method = RequestMethodPost
	request.Data = data

	connection := NewConnection()
	response, error := connection.Start(request)

	if error != nil {
		Logger().Errorf("Error during CreateChild: %s", error.Error())
		return error
	}

	err := json.Unmarshal(response.Data[1:len(response.Data)-1], child)

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
func AssignChildren(parent Exposable, children interface{}, identity Identity) *Error {

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
