// Copyright 2016 Marcel Gotsch. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goserv

// RequestContext stores key-value pairs supporting all data types.
type RequestContext struct {
	store map[string]interface{}
}

// Set sets the value for the specified the key. It replaces any existing values.
func (c *RequestContext) Set(key string, value interface{}) {
	c.store[key] = value
}

// Get retrieves the value for key. If the key doesn't exist in the RequestContext,
// Get returns nil.
func (c *RequestContext) Get(key string) interface{} {
	return c.store[key]
}

// Delete deletes the value associated with key. If the key doesn't exist nothing happens.
func (c *RequestContext) Delete(key string) {
	delete(c.store, key)
}

// Exists returns true if the specified key exists in the RequestContext, otherwise false is returned.
func (c *RequestContext) Exists(key string) bool {
	_, exists := c.store[key]
	return exists
}

func newRequestContext() *RequestContext {
	return &RequestContext{make(map[string]interface{})}
}

// Stores a RequestContext for each Request.
var requestContextMap = make(map[*Request]*RequestContext)

// Context returns the corresponding RequestContext for the given Request.
func Context(r *Request) *RequestContext {
	return requestContextMap[r]
}

// Stores a new RequestContext for the specified Request in the requestContextMap.
// This may overwrite an existing RequestContext!
func createRequestContext(r *Request) {
	requestContextMap[r] = newRequestContext()
}
