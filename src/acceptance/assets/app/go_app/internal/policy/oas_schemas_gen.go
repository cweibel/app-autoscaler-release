// Code generated by ogen, DO NOT EDIT.

package api

import (
	"fmt"

	"github.com/go-faster/jx"
)

func (s *ErrorResponseStatusCode) Error() string {
	return fmt.Sprintf("code %d: %+v", s.StatusCode, s.Response)
}

type BasicAuthentication struct {
	Username string
	Password string
}

// GetUsername returns the value of Username.
func (s *BasicAuthentication) GetUsername() string {
	return s.Username
}

// GetPassword returns the value of Password.
func (s *BasicAuthentication) GetPassword() string {
	return s.Password
}

// SetUsername sets the value of Username.
func (s *BasicAuthentication) SetUsername(val string) {
	s.Username = val
}

// SetPassword sets the value of Password.
func (s *BasicAuthentication) SetPassword(val string) {
	s.Password = val
}

// Ref: #/components/schemas/ErrorResponse
type ErrorResponse struct{}

// ErrorResponseStatusCode wraps ErrorResponse with StatusCode.
type ErrorResponseStatusCode struct {
	StatusCode int
	Response   ErrorResponse
}

// GetStatusCode returns the value of StatusCode.
func (s *ErrorResponseStatusCode) GetStatusCode() int {
	return s.StatusCode
}

// GetResponse returns the value of Response.
func (s *ErrorResponseStatusCode) GetResponse() ErrorResponse {
	return s.Response
}

// SetStatusCode sets the value of StatusCode.
func (s *ErrorResponseStatusCode) SetStatusCode(val int) {
	s.StatusCode = val
}

// SetResponse sets the value of Response.
func (s *ErrorResponseStatusCode) SetResponse(val ErrorResponse) {
	s.Response = val
}

// NewOptInt64 returns new OptInt64 with value set to v.
func NewOptInt64(v int64) OptInt64 {
	return OptInt64{
		Value: v,
		Set:   true,
	}
}

// OptInt64 is optional int64.
type OptInt64 struct {
	Value int64
	Set   bool
}

// IsSet returns true if OptInt64 was set.
func (o OptInt64) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptInt64) Reset() {
	var v int64
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptInt64) SetTo(v int64) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptInt64) Get() (v int64, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptInt64) Or(d int64) int64 {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptPolicy returns new OptPolicy with value set to v.
func NewOptPolicy(v Policy) OptPolicy {
	return OptPolicy{
		Value: v,
		Set:   true,
	}
}

// OptPolicy is optional Policy.
type OptPolicy struct {
	Value Policy
	Set   bool
}

// IsSet returns true if OptPolicy was set.
func (o OptPolicy) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptPolicy) Reset() {
	var v Policy
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptPolicy) SetTo(v Policy) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptPolicy) Get() (v Policy, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptPolicy) Or(d Policy) Policy {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// Object creating Policy.
// Ref: #/components/schemas/Policy
type Policy struct {
	// Minimal number of instance count.
	InstanceMinCount OptInt64 `json:"instance_min_count"`
	// Maximal number of instance count.
	InstanceMaxCount OptInt64      `json:"instance_max_count"`
	ScalingRules     []ScalingRule `json:"scaling_rules"`
}

// GetInstanceMinCount returns the value of InstanceMinCount.
func (s *Policy) GetInstanceMinCount() OptInt64 {
	return s.InstanceMinCount
}

// GetInstanceMaxCount returns the value of InstanceMaxCount.
func (s *Policy) GetInstanceMaxCount() OptInt64 {
	return s.InstanceMaxCount
}

// GetScalingRules returns the value of ScalingRules.
func (s *Policy) GetScalingRules() []ScalingRule {
	return s.ScalingRules
}

// SetInstanceMinCount sets the value of InstanceMinCount.
func (s *Policy) SetInstanceMinCount(val OptInt64) {
	s.InstanceMinCount = val
}

// SetInstanceMaxCount sets the value of InstanceMaxCount.
func (s *Policy) SetInstanceMaxCount(val OptInt64) {
	s.InstanceMaxCount = val
}

// SetScalingRules sets the value of ScalingRules.
func (s *Policy) SetScalingRules(val []ScalingRule) {
	s.ScalingRules = val
}

type ScalingRule jx.Raw
