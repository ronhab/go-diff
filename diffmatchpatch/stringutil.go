// Copyright (c) 2012-2016 The go-diff authors. All rights reserved.
// https://github.com/sergi/go-diff
// See the included LICENSE file for license details.
//
// go-diff is a Go implementation of Google's Diff, Match, and Patch library
// Original library is Copyright (c) 2006 Google Inc.
// http://code.google.com/p/google-diff-match-patch/

package diffmatchpatch

import (
	"strings"
	"unicode/utf8"
)

// unescaper unescapes selected chars for compatibility with JavaScript's encodeURI.
// In speed critical applications this could be dropped since the receiving application will certainly decode these fine. Note that this function is case-sensitive.  Thus "%3F" would not be unescaped.  But this is ok because it is only called with the output of HttpUtility.UrlEncode which returns lowercase hex. Example: "%3f" -> "?", "%24" -> "$", etc.
var unescaper = strings.NewReplacer(
	"%21", "!", "%7E", "~", "%27", "'",
	"%28", "(", "%29", ")", "%3B", ";",
	"%2F", "/", "%3F", "?", "%3A", ":",
	"%40", "@", "%26", "&", "%3D", "=",
	"%2B", "+", "%24", "$", "%2C", ",", "%23", "#", "%2A", "*")

// indexOf returns the first index of pattern in str, starting at str[i].
func indexOf(str string, pattern string, i int) int {
	if i > len(str)-1 {
		return -1
	}
	if i <= 0 {
		return strings.Index(str, pattern)
	}
	ind := strings.Index(str[i:], pattern)
	if ind == -1 {
		return -1
	}
	return ind + i
}

// lastIndexOf returns the last index of pattern in str, starting at str[i].
func lastIndexOf(str string, pattern string, i int) int {
	if i < 0 {
		return -1
	}
	if i >= len(str) {
		return strings.LastIndex(str, pattern)
	}
	_, size := utf8.DecodeRuneInString(str[i:])
	return strings.LastIndex(str[:i+size], pattern)
}

// charsIndexOf returns the index of pattern in target, starting at target[i].
func charsIndexOf(target, pattern []DiffChar, i int) int {
	if i > len(target)-1 {
		return -1
	}
	if i <= 0 {
		return charsIndex(target, pattern)
	}
	ind := charsIndex(target[i:], pattern)
	if ind == -1 {
		return -1
	}
	return ind + i
}

func charsEqual(c1, c2 []DiffChar) bool {
	if len(c1) != len(c2) {
		return false
	}
	for i, c := range c1 {
		if c != c2[i] {
			return false
		}
	}
	return true
}

// charsIndex is the equivalent of strings.Index for DiffChar slices.
func charsIndex(c1, c2 []DiffChar) int {
	last := len(c1) - len(c2)
	for i := 0; i <= last; i++ {
		if charsEqual(c1[i:i+len(c2)], c2) {
			return i
		}
	}
	return -1
}

func runesToChars(runes []rune) []DiffChar {
	chars := make([]DiffChar, len(runes))
	for i, r := range(runes) {
		chars[i] = DiffChar(r)
	}
	return chars
}

func charsToRunes(chars []DiffChar) []rune {
	runes := make([]rune, len(chars))
	for i, c := range(chars) {
		runes[i] = rune(c)
	}
	return runes
}

func encodeChar(c DiffChar) string {
	return string([]byte{ byte(c >> 24 & 0xFF), byte(c >> 16 & 0xFF), byte(c >> 8 & 0xFF), byte(c & 0xFF) })
}

func decodeChar(s string) DiffChar {
	return (DiffChar(s[0]) << 24) | (DiffChar(s[1]) << 16) | (DiffChar(s[2]) << 8) | DiffChar(s[3])
}

func charsToString(chars []DiffChar, encode bool) string {
	if encode {
		res := ""
		for _, c := range chars {
			res += encodeChar(c)
		}
		return res
	}
	return string(charsToRunes(chars))
}

func stringToChars(s string, decode bool) []DiffChar {
	if decode {
		chars := make([]DiffChar, len(s) / 4)
		for i := 0; i < len(s) / 4; i++ {
			chars[i] = decodeChar(s[i*4:i*4+4])
		}
		return chars
	}
	return runesToChars([]rune(s))
}
