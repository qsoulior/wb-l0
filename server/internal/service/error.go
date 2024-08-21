package service

import "errors"

var (
	ErrNotExist = errors.New("entity does not exist")
	ErrExists   = errors.New("entity already exists")
)
