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
	"fmt"
	"reflect"
)

// AllIdentity represents all possible Identities.
var AllIdentity = Identity{
	RESTName:     "__all__",
	ResourceName: "__all__",
}

// IdentifiablesList represents a list of Identifiable.
type IdentifiablesList []Identifiable

// Identifiable is the interface for objects that can be identified using an Identity.
type Identifiable interface {
	GetIdentity() Identity
}

// Identity is a structure that contains the necessary information about an Identifiable.
//
// The RESTName is usually the singular form of the ResourceName.
// For instance, "enterprise" and "enterprises".
type Identity struct {
	RESTName     string
	ResourceName string
}

// String returns the string representation of the identity.
func (i Identity) String() string {

	return fmt.Sprintf("<%s|%s>", i.RESTName, i.ResourceName)
}

// identify applies the given Identity to a list of Exposables.
func identify(list interface{}, identity Identity) {

	l := reflect.ValueOf(list).Elem().Len()

	for i := 0; i < l; i++ {

		o := reflect.ValueOf(list).Elem().Index(i).Elem()

		identityField := o.FieldByName("Identity")
		RESTNameField := identityField.FieldByName("RESTName")
		resourceNameField := identityField.FieldByName("ResourceName")

		RESTNameField.SetString(identity.RESTName)
		resourceNameField.SetString(identity.ResourceName)
	}
}
