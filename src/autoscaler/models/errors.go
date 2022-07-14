package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CFErrorResponse struct {
	Description string `json:"description"`
	ErrorCode   string `json:"error_code"`
	Code        int    `json:"code"`
}

var CfResourceNotFound = &CfError{Errors: []CfErrorItem{{Detail: "App usage event not found", Title: "CF-ResourceNotFound", Code: 10010}}}
var CfInternalServerError = &CfError{Errors: []CfErrorItem{{Detail: "An unexpected, uncaught error occurred; the CC logs will contain more information", Title: "UnknownError", Code: 10001}}}
var _ error = &CfError{}
var ErrInvalidJson = fmt.Errorf("invalid error json")

func NewCfError(resourceId string, statusCode int, respBody []byte) error {
	var cfError = &CfError{}
	err := json.Unmarshal(respBody, &cfError)
	if err != nil {
		return fmt.Errorf("failed to unmarshal id:%s error status '%d' body:'%s' : %w", resourceId, statusCode, truncateString(string(respBody), 512), err)
	}
	cfError.ResourceId = resourceId
	cfError.StatusCode = statusCode

	if !cfError.IsValid() {
		return fmt.Errorf("invalid cfError: resource %s status:%d body:%s :%w", resourceId, statusCode, truncateString(string(respBody), 512), ErrInvalidJson)
	}
	return cfError
}

//CfError cf V3 Error type
type CfError struct {
	Errors     []CfErrorItem `json:"errors"`
	StatusCode int
	ResourceId string
}

type CfErrorItem struct {
	Code   int    `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (c *CfError) Error() string {
	errors := []string{}
	message := "None found"
	for _, errorItem := range c.Errors {
		errorsString := fmt.Sprintf("['%s' code: %d, Detail: '%s']", errorItem.Title, errorItem.Code, errorItem.Detail)
		errors = append(errors, errorsString)
	}
	if len(errors) > 0 {
		message = strings.Join(errors, ", ")
	}
	return fmt.Sprintf("cf api Error: %s", message)
}

func (c *CfError) IsNotFound() bool {
	if c.IsValid() {
		for _, item := range c.Errors {
			if item.Code == 10010 {
				return true
			}
		}
	}
	return false
}

func (c *CfError) IsValid() bool {
	return c != nil && len(c.Errors) > 0
}

func truncateString(stringToTrunk string, length int) string {
	if len(stringToTrunk) > length {
		return stringToTrunk[:length]
	}
	return stringToTrunk
}
