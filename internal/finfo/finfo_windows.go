//go:build windows

package finfo

import (
	"errors"
)


// Stat reads path's attributes, of a link itself rather than its target.
func Stat(path string) (Info, error) {
	return Info{}, errors.New("not implemented")
}
