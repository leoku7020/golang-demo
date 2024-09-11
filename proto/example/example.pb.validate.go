// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: example/example.proto

package example

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on LoginReq with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LoginReq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoginReq with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in LoginReqMultiError, or nil
// if none found.
func (m *LoginReq) ValidateAll() error {
	return m.validate(true)
}

func (m *LoginReq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetUsername()) < 1 {
		err := LoginReqValidationError{
			field:  "Username",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetPassword()) < 1 {
		err := LoginReqValidationError{
			field:  "Password",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return LoginReqMultiError(errors)
	}

	return nil
}

// LoginReqMultiError is an error wrapping multiple validation errors returned
// by LoginReq.ValidateAll() if the designated constraints aren't met.
type LoginReqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoginReqMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoginReqMultiError) AllErrors() []error { return m }

// LoginReqValidationError is the validation error returned by
// LoginReq.Validate if the designated constraints aren't met.
type LoginReqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoginReqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoginReqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoginReqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoginReqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoginReqValidationError) ErrorName() string { return "LoginReqValidationError" }

// Error satisfies the builtin error interface
func (e LoginReqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoginReq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoginReqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoginReqValidationError{}

// Validate checks the field values on LoginResp with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LoginResp) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoginResp with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in LoginRespMultiError, or nil
// if none found.
func (m *LoginResp) ValidateAll() error {
	return m.validate(true)
}

func (m *LoginResp) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Status

	// no validation rules for Description

	if len(errors) > 0 {
		return LoginRespMultiError(errors)
	}

	return nil
}

// LoginRespMultiError is an error wrapping multiple validation errors returned
// by LoginResp.ValidateAll() if the designated constraints aren't met.
type LoginRespMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoginRespMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoginRespMultiError) AllErrors() []error { return m }

// LoginRespValidationError is the validation error returned by
// LoginResp.Validate if the designated constraints aren't met.
type LoginRespValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoginRespValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoginRespValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoginRespValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoginRespValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoginRespValidationError) ErrorName() string { return "LoginRespValidationError" }

// Error satisfies the builtin error interface
func (e LoginRespValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoginResp.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoginRespValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoginRespValidationError{}

// Validate checks the field values on ListItemsReq with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ListItemsReq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListItemsReq with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ListItemsReqMultiError, or
// nil if none found.
func (m *ListItemsReq) ValidateAll() error {
	return m.validate(true)
}

func (m *ListItemsReq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetUsername()) < 1 {
		err := ListItemsReqValidationError{
			field:  "Username",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetItem()) < 1 {
		err := ListItemsReqValidationError{
			field:  "Item",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return ListItemsReqMultiError(errors)
	}

	return nil
}

// ListItemsReqMultiError is an error wrapping multiple validation errors
// returned by ListItemsReq.ValidateAll() if the designated constraints aren't met.
type ListItemsReqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListItemsReqMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListItemsReqMultiError) AllErrors() []error { return m }

// ListItemsReqValidationError is the validation error returned by
// ListItemsReq.Validate if the designated constraints aren't met.
type ListItemsReqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListItemsReqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListItemsReqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListItemsReqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListItemsReqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListItemsReqValidationError) ErrorName() string { return "ListItemsReqValidationError" }

// Error satisfies the builtin error interface
func (e ListItemsReqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListItemsReq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListItemsReqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListItemsReqValidationError{}

// Validate checks the field values on ListItemsResp with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ListItemsResp) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListItemsResp with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ListItemsRespMultiError, or
// nil if none found.
func (m *ListItemsResp) ValidateAll() error {
	return m.validate(true)
}

func (m *ListItemsResp) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Status

	for idx, item := range m.GetItem() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListItemsRespValidationError{
						field:  fmt.Sprintf("Item[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListItemsRespValidationError{
						field:  fmt.Sprintf("Item[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListItemsRespValidationError{
					field:  fmt.Sprintf("Item[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ListItemsRespMultiError(errors)
	}

	return nil
}

// ListItemsRespMultiError is an error wrapping multiple validation errors
// returned by ListItemsResp.ValidateAll() if the designated constraints
// aren't met.
type ListItemsRespMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListItemsRespMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListItemsRespMultiError) AllErrors() []error { return m }

// ListItemsRespValidationError is the validation error returned by
// ListItemsResp.Validate if the designated constraints aren't met.
type ListItemsRespValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListItemsRespValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListItemsRespValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListItemsRespValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListItemsRespValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListItemsRespValidationError) ErrorName() string { return "ListItemsRespValidationError" }

// Error satisfies the builtin error interface
func (e ListItemsRespValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListItemsResp.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListItemsRespValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListItemsRespValidationError{}

// Validate checks the field values on ItemData with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ItemData) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ItemData with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ItemDataMultiError, or nil
// if none found.
func (m *ItemData) ValidateAll() error {
	return m.validate(true)
}

func (m *ItemData) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ItemId

	// no validation rules for ItemName

	// no validation rules for Category

	if len(errors) > 0 {
		return ItemDataMultiError(errors)
	}

	return nil
}

// ItemDataMultiError is an error wrapping multiple validation errors returned
// by ItemData.ValidateAll() if the designated constraints aren't met.
type ItemDataMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ItemDataMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ItemDataMultiError) AllErrors() []error { return m }

// ItemDataValidationError is the validation error returned by
// ItemData.Validate if the designated constraints aren't met.
type ItemDataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ItemDataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ItemDataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ItemDataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ItemDataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ItemDataValidationError) ErrorName() string { return "ItemDataValidationError" }

// Error satisfies the builtin error interface
func (e ItemDataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sItemData.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ItemDataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ItemDataValidationError{}