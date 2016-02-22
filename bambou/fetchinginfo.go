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
	"strconv"
	"strings"
)

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

	return fmt.Sprintf("<FetchingInfo page: %d, pagesize: %d, totalcount: %d>", f.Page, f.PageSize, f.TotalCount)
}

// Private.
//
// Fills the HTTP headers of the given request according to the given FetchingInfo.
func prepareHeaders(request *request, info *FetchingInfo) {

	request.setHeader("X-Nuage-PageSize", "50")

	if info == nil {
		return
	}

	if info.Filter != "" {
		request.setHeader("X-Nuage-Filter", info.Filter)
	}

	if info.OrderBy != "" {
		request.setHeader("X-Nuage-OrderBy", info.OrderBy)
	}

	if info.Page != -1 {
		request.setHeader("X-Nuage-Page", strconv.Itoa(info.Page))
	}

	if info.PageSize > 0 {
		request.setHeader("X-Nuage-PageSize", strconv.Itoa(info.PageSize))
	}

	if len(info.GroupBy) > 0 {
		request.setHeader("X-Nuage-GroupBy", "true")
		request.setHeader("X-Nuage-Attributes", strings.Join(info.GroupBy, ", "))
	}
}

// Private.
//
// Fills the given FetchingInfo according to the HTTP headers of the given response.
func readHeaders(response *response, info *FetchingInfo) {

	if info == nil {
		return
	}

	info.Filter = response.getHeader("X-Nuage-Filter")
	info.FilterType = response.getHeader("X-Nuage-FilterType")
	// info.GroupBy = response.getHeader("X-Nuage-GroupBy")
	info.OrderBy = response.getHeader("X-Nuage-OrderBy")
	info.Page, _ = strconv.Atoi(response.getHeader("X-Nuage-Page"))
	info.PageSize, _ = strconv.Atoi(response.getHeader("X-Nuage-PageSize"))
	info.TotalCount, _ = strconv.Atoi(response.getHeader("X-Nuage-Count"))
}
