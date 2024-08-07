package product

import "errors"

var (
	ErrRemoteProductNotFound      = errors.New("no remote product found")
	ErrRemoteLocalProductNotFound = errors.New("no remote local product found")
	ErrProductNoChange            = errors.New("no product attributes change")
)
