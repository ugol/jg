// Copyright © 2024 JR team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package ctx

import (
	"sync"
	"time"
)

var JrContext *Context

// Context is the object passed on the templates which contains all the needed details.
type Context struct {
	StartTime                 time.Time
	GeneratedObjects          int64
	ExpectedObjects           int64
	GeneratedBytes            int64
	Locale                    string
	CtxCounters               map[string]int
	CtxCountersLock           sync.RWMutex
	Ctx                       map[string]string
	CtxLock                   sync.RWMutex
	CtxList                   map[string][]string
	CtxListLock               sync.RWMutex
	CtxCSV                    map[int]map[string]string
	CtxCSVLock                sync.RWMutex
	CtxGeoJson                [][]float64
	CtxGeoJsonLock            sync.RWMutex
	CtxLastPointLat           []float64
	CtxLastPointLon           []float64
	CtxForward                bool
	CtxIndex                  int
	LastIndex                 int
	CountryIndex              int
	CityIndex                 int
	CurrentIterationLoopIndex int
}

func init() {
	var ctxgeojson [][]float64
	JrContext = &Context{
		StartTime:        time.Now(),
		GeneratedBytes:   0,
		GeneratedObjects: 0,
		Locale:           "us",
		CtxCounters:      make(map[string]int),
		CtxCountersLock:  sync.RWMutex{},
		Ctx:              make(map[string]string),
		CtxLock:          sync.RWMutex{},
		CtxList:          make(map[string][]string),
		CtxListLock:      sync.RWMutex{},
		CtxCSV:           make(map[int]map[string]string),
		CtxCSVLock:       sync.RWMutex{},
		CtxGeoJson:       ctxgeojson,
		CtxGeoJsonLock:   sync.RWMutex{},
		CtxLastPointLat:  []float64{},
		CtxLastPointLon:  []float64{},
		CtxForward:       true,
		CtxIndex:         0,
		LastIndex:        -1,
		CountryIndex:     232,
		CityIndex:        -1,
	}
}
