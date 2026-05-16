package client

import "errors"

var ErrLoginAlreadyExists = errors.New("login already exists")
var ErrClientNotFound = errors.New("client not found")
