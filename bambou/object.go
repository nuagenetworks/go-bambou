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
	"fmt"
)

// Exposable is the interface of objects that can be retrieved and sent to the server.
// An Exposable implements the Identifiable and Operationable interfaces.
type Exposable interface {
	Fetch() *Error
	Save() *Error
	Delete() *Error
	GetIdentity() Identity
	SetIdentity(Identity)
	GetID() string
}

// Rootable is the interface that must be implemented by the root object of the API.
// A Rootable also implements the Exposable. Rootable
type Rootable interface {
	Exposable

	GetAPIKey() string
	SetAPIKey(string)
}

// ExposedObject represents an object than contains information common to all objects.
// exposed by the server.
// This struct must be embedded into all objects that are available
// throught the ReST api.
type ExposedObject struct {
	ID           string   `json:"ID,omitempty"`
	ParentID     string   `json:"parentID,omitempty"`
	ParentType   string   `json:"parentType,omitempty"`
	Owner        string   `json:"owner,omitempty"`
	ParentObject string   `json:"-"`
	Identity     Identity `json:"-"`
}

// GetIdentity returns the Identity.
func (o *ExposedObject) GetIdentity() Identity {

	return o.Identity
}

// SetIdentity sets the Identity of the object.
func (o *ExposedObject) SetIdentity(identity Identity) {

	o.Identity = identity
}

// GetID returns the ID of the object.
func (o *ExposedObject) GetID() string {

	return o.ID
}

// String returns the string representation of the object.
func (o *ExposedObject) String() string {

	return fmt.Sprintf("<ExposedObject %s:%s>", o.Identity.RESTName, o.ID)
}
