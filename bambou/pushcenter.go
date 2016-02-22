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
)

// Signature of a Push Center Handler.
type EventHandler func(*Event)

// Represents a list of Push Center Handlers
type EventHandlers map[string]EventHandler

// Represents a Push Center
type PushCenter struct {
	IsRunning bool
	Channel   chan *Notification

	handlers      EventHandlers
	defaultHander EventHandler
	lastEventID   string
	stop          chan bool
}

// Creates a new Push Center
func NewPushCenter() *PushCenter {

	return &PushCenter{
		Channel:  make(chan *Notification),
		stop:     make(chan bool),
		handlers: EventHandlers{},
	}
}

// Registers the given EventHandler for the given Entity Identity.
//
// You can pass the bambou.AllIdentity as identity to register the handler
// for all events. If you pass a handler for an Identity that is already registered
// the previous handler will be silently overwriten.
func (p *PushCenter) RegisterHandlerForIdentity(handler EventHandler, identity Identity) {

	if identity.RESTName == AllIdentity.RESTName {
		p.defaultHander = handler
		return
	}

	p.handlers[identity.RESTName] = handler
}

// Registers the given EventHandler for the given Entity Identity.
func (p *PushCenter) UnregisterHandlerForIdentity(identity Identity) {

	if identity.RESTName == AllIdentity.RESTName {
		p.defaultHander = nil
		return
	}

	if _, exists := p.handlers[identity.RESTName]; exists {
		delete(p.handlers, identity.RESTName)
	}
}

func (p *PushCenter) HasHandlerForIdentity(identity Identity) bool {

	_, exists := p.handlers[identity.RESTName]
	return exists
}

// Starts the Push Center.
func (p *PushCenter) Start() {

	p.IsRunning = true
	p.lastEventID = ""

	go func() {
		for {
			go p.listen()

			select {
			case notification := <-p.Channel:

				for _, event := range notification.Events {

					event.Data, _ = json.Marshal(event.DataMap[0])

					if p.defaultHander != nil {
						p.defaultHander(event)
					}
					if handler, exists := p.handlers[event.EntityType]; exists {
						handler(event)
					}
				}
			case <-p.stop:
				return
			}
		}
	}()
}

// Stops a running PushCenter.
func (p *PushCenter) Stop() {

	p.IsRunning = false
	p.lastEventID = ""
	p.stop <- true
}

// Private.
//
// Will handle the creation of new *Notification
func (p *PushCenter) listen() {

	currentURL := CurrentSession().URL + "/events"

	if p.lastEventID != "" {
		currentURL += "?uuid=" + p.lastEventID
	}

	request := NewRequest(currentURL)
	connection := NewConnection()
	response, error := connection.Start(request)

	// if the push center not running anymore, return
	if !p.IsRunning {
		return
	}

	if error != nil {
		Logger().Errorf("Error during push: %s", error.Error())
		return
	}

	notification := NewNotification()
	err := json.Unmarshal(response.Data, &notification)

	if err != io.EOF && err != nil {
		Logger().Errorf("Error during push: %s", err.Error())
		return
	}

	p.lastEventID = notification.UUID

	if len(notification.Events) > 0 {
		p.Channel <- notification
	}
}
