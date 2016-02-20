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
type PushCenterHander func(*Notification)

// Represents a Push Center
type PushCenter struct {
	IsRunning bool
	Channel   chan *Notification
	URL       string
}

// Creates a new Push Center
func NewPushCenter() *PushCenter {

	return &PushCenter{}
}

// Starts the Push Center with the given PushCenterHander.
//
// When a Push Notification is received from the server, the PushCenter
// will call the handler and pass it the *Notification.
func (p *PushCenter) StartWithHandler(hander PushCenterHander) {

	go func() {

		p.Start()

	loop:
		for {
			select {
			case notification, ok := <-p.Channel:
				if !ok {
					break loop
				} else {
					hander(notification)
				}
			}
		}
	}()
}

// Starts the Push Center.
//
// The Push Center will be started, but you will need to use the
// internal PushCenter.Channel to deal with the event yourself.
func (p *PushCenter) Start() {

	p.IsRunning = true
	p.URL = CurrentSession().URL + "/events"
	p.Channel = make(chan *Notification)

	go func() { p.listen() }()
}

// Stops a running PushCenter.
func (p *PushCenter) Stop() {

	p.IsRunning = false
	p.URL = ""
	close(p.Channel)
}

// Private.
//
// Will handle the creation of new *Notification
func (p *PushCenter) listen() {

	lastEventID := ""

	for {
		currentURL := p.URL

		if lastEventID != "" {
			currentURL += "?uuid=" + lastEventID
		}

		request := NewRequest(currentURL)
		connection := NewConnection()
		response, error := connection.Start(request)

		if error != nil {
			Logger().Errorf("Error during push: %s", error.Error())
			continue
		}

		notification := NewNotification()
		err := json.Unmarshal(response.Data, &notification)

		if err != io.EOF && err != nil {
			Logger().Errorf("Error during push: %s", err.Error())
			continue
		}

		lastEventID = notification.UUID

		if len(notification.Events) > 0 {
			p.Channel <- notification
		}
	}
}
