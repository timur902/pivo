package sellerusecase

import "errors"

var ErrSellerNotFound = errors.New("seller not found")
var ErrLoginAlreadyExists = errors.New("login already exists")
