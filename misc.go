package oid

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	sprintf func(string, ...any) string = fmt.Sprintf

	atoi func(string) (int, error) = strconv.Atoi
	itoa func(int) string          = strconv.Itoa

	contains   func(string, string) bool          = strings.Contains
	eq         func(string, string) bool          = strings.EqualFold
	fields     func(string) []string              = strings.Fields
	hasPrefix  func(string, string) bool          = strings.HasPrefix
	hasSuffix  func(string, string) bool          = strings.HasSuffix
	indexRune  func(string, rune) int             = strings.IndexRune
	join       func([]string, string) string      = strings.Join
	split      func(string, string) []string      = strings.Split
	splitAfter func(string, string) []string      = strings.SplitAfter
	splitN     func(string, string, int) []string = strings.SplitN
	trimL      func(string, string) string        = strings.TrimLeft
	trimR      func(string, string) string        = strings.TrimRight
)

func errorf(msg any, x ...any) error {
	switch tv := msg.(type) {
	case string:
		return errors.New(sprintf(tv, x...))
	case error:
		return errors.New(sprintf(tv.Error(), x...))
	}

	return nil
}

func strInSlice(str string, slice []string) bool {
	if len(str) == 0 || len(slice) == 0 {
		return false
	}

	for _, val := range slice {
		if eq(val, str) {
			return true
		}
	}

	return false
}

/*
is 'val' a digit?
*/
func isDigit(val string) bool {
	for _, c := range val {
		if '0' <= c && c <= '9' {
			continue
		}
		return false
	}
	return true
}

/*
compare slice members of two (2) []int instances.
*/
func intSliceEqual(s1, s2 []int) (equal bool) {
	if len(s1)|len(s2) == 0 || len(s1) != len(s2) {
		return
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return
		}
	}

	equal = true
	return
}

/*
compare slice members of two (2) []string instances.
*/
func strSliceEqual(s1, s2 []string) (equal bool) {
	if len(s1)|len(s2) == 0 || len(s1) != len(s2) {
		return
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return
		}
	}

	equal = true
	return
}
