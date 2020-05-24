package exception

import "errors"

// InvalidID Entity has an invalid ID
var InvalidID = errors.New("invalid id")

// RequiredField A field is required
var RequiredField = errors.New("missing required request field")
var RequiredFieldString = "missing required request field %v"

// InvalidFieldFormat A field has a bad format
var InvalidFieldFormat = errors.New("request field has an invalid format")
var InvalidFieldFormatString = "request field %v has an invalid format, expected %v"

// InvalidFieldRange A field is out of range
var InvalidFieldRange = errors.New("request field is out of range")
var InvalidFieldRangeString = "request field %s is out of range [%s, %s)"

// EmptyBody Body is empty
var EmptyBody = errors.New("no parameters found")

// EntityNotFound Entity was not found
var EntityNotFound = errors.New("resource not found")

// EntitiesNotFound Entities were not found
var EntitiesNotFound = errors.New("resources not found")

// EntityExists Entity was already created
var EntityExists = errors.New("resource already exists")
