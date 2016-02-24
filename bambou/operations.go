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
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

// FetchEntity fetchs the given Exposable from the server.
// You should not use this function by yourself.
func FetchEntity(object Exposable) *Error {

	request := newRequest(object.GetPersonalURL())
	connection := newConnection()
	response, error := connection.start(request)

	if error != nil {
		Logger().Errorf("Error during FetchEntity: %s", error.Error())
		return error
	}

	data := response.Data[1 : len(response.Data)-1]
	err := json.Unmarshal(data, &object)

	if err != io.EOF && err != nil {
		panic(fmt.Sprintf("Unable to unmarshal json %s: %s", string(data), err.Error()))
	}

	return nil
}

// SaveEntity saves the given Exposable into the server.
// You should not use this function by yourself.
func SaveEntity(object Exposable) *Error {

	data, _ := json.Marshal(object)

	request := newRequest(object.GetPersonalURL())
	request.Method = requestMethodPut
	request.Data = data

	connection := newConnection()
	response, err1 := connection.start(request)

	if err1 != nil {
		Logger().Errorf("Error during SaveEntity: %s", err1.Error())
		return err1
	}

	if len(response.Data) > 0 {
		data := response.Data[1 : len(response.Data)-1]
		err2 := json.Unmarshal(data, &object)

		if err2 != io.EOF && err2 != nil {
			panic(fmt.Sprintf("Unable to unmarshal json %s: %s", string(data), err2.Error()))
		}
	}

	return nil
}

// DeleteEntity deletes the given Exposable from the server.
// You should not use this function by yourself.
func DeleteEntity(object Exposable) *Error {

	request := newRequest(object.GetPersonalURL())
	request.Method = requestMethodDelete

	connection := newConnection()
	_, error := connection.start(request)

	if error != nil {
		Logger().Errorf("Error during DeleteEntity: %s", error.Error())
		return error
	}

	return nil
}

// FetchChildren fetches the children with of given parent identified by the given identify.
// You should not use this function by yourself.
func FetchChildren(parent Exposable, identity Identity, dest interface{}, info *FetchingInfo) *Error {

	request := newRequest(parent.GetURLForChildrenIdentity(identity))
	prepareHeaders(request, info)

	connection := newConnection()
	response, error := connection.start(request)

	if error != nil {
		Logger().Errorf("Error during FetchChildren: %s", error.Error())
		return error
	}

	readHeaders(response, info)

	if response.Code == responseCodeEmpty {
		return nil
	}

	data := response.Data
	err := json.Unmarshal(data, dest)

	if err != io.EOF && err != nil {
		panic(fmt.Sprintf("Unable to unmarshal json %s: %s", string(data), err.Error()))
	}

	identify(dest, identity)

	return nil
}

// CreateChild creates a new child Exposable under the given parent Exposable in the server.
// You should not use this function by yourself.
func CreateChild(parent Exposable, child Exposable) *Error {

	data, _ := json.Marshal(child)

	request := newRequest(parent.GetURLForChildrenIdentity(child.GetIdentity()))
	request.Method = requestMethodPost
	request.Data = data

	connection := newConnection()
	response, error := connection.start(request)

	if error != nil {
		Logger().Errorf("Error during CreateChild: %s", error.Error())
		return error
	}

	data = response.Data[1 : len(response.Data)-1]
	err := json.Unmarshal(data, child)

	if err != io.EOF && err != nil {
		panic(fmt.Sprintf("Unable to unmarshal json %s: %s", string(data), err.Error()))
	}

	return nil
}

// AssignChildren assigns the list of given child Exposables to the given Exposable parent in the server.
// You should not use this function by yourself.
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

	request := newRequest(parent.GetURLForChildrenIdentity(identity))
	request.Method = requestMethodPut
	request.Data = data

	connection := newConnection()
	_, error := connection.start(request)

	if error != nil {
		Logger().Errorf("Error during AssignChildren: %s", error.Error())
		return error
	}

	return nil
}
