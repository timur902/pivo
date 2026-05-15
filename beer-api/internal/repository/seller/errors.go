package seller

import "errors"

var ErrLoginAlreadyExists = errors.New("login already exists")
var ErrSellerNotFound = errors.New("seller not found")
