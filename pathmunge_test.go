package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestGenList(t *testing.T) {
	cases := []struct {
		paths  []string
		dir    string
		result []string
	}{
		{
			[]string{"abc"},
			"def",
			[]string{"def", "abc"},
		},
		{
			[]string{"abc"},
			"abc",
			[]string{"abc"},
		},
		{
			[]string{"def", "abc"},
			"abc",
			[]string{"def", "abc"},
		},
		{
			[]string{"def", "abc", "ghi"},
			"abc",
			[]string{"def", "abc", "ghi"},
		},
		{
			[]string{"def", "abc", "abc", "def"},
			"ghi",
			[]string{"ghi", "def", "abc"},
		},
	}
	for _, c := range cases {
		p := strings.Join(c.paths, string(filepath.ListSeparator))
		r := genList(p, c.dir, optType{})
		if !equal(c.result, r) {
			t.Errorf("expectd: %v : %v", c.result, r)
		}
	}
}


func TestListDel(t *testing.T) {
	cases := []struct {
		paths  []string
		dir    string
		result []string
	}{
		{
			[]string{"abc"},
			"def",
			[]string{"abc"},
		},
		{
			[]string{"abc"},
			"abc",
			[]string{},
		},
		{
			[]string{"def", "abc"},
			"abc",
			[]string{"def"},
		},
		{
			[]string{"def", "abc"},
			"def",
			[]string{"abc"},
		},
		{
			[]string{"def", "abc", "ghi", "def"},
			"abc",
			[]string{"def", "ghi"},
		},
		{
			[]string{"def", "abc", "abc", "def"},
			"ghi",
			[]string{"def", "abc"},
		},
	}
	for _, c := range cases {
		p := strings.Join(c.paths, string(filepath.ListSeparator))
		r := genList(p, c.dir, optType{del:true})
		if !equal(c.result, r) {
			t.Errorf("expectd: %v : %v", c.result, r)
		}
	}
}

func TestAddToHead(t *testing.T) {
	cases := []struct {
		paths  []string
		dir    string
		result []string
	}{
		{
			[]string{"abc"},
			"def",
			[]string{"def", "abc"},
		},
		{
			[]string{"abc"},
			"abc",
			[]string{"abc"},
		},
		{
			[]string{"def", "abc"},
			"abc",
			[]string{"abc", "def"},
		},
		{
			[]string{"def", "abc"},
			"def",
			[]string{"def", "abc"},
		},
		{
			[]string{"def", "abc", "ghi", "def"},
			"abc",
			[]string{"abc", "def", "ghi"},
		},
		{
			[]string{"def", "abc", "abc", "def"},
			"ghi",
			[]string{"ghi", "def", "abc"},
		},
	}
	for _, c := range cases {
		p := strings.Join(c.paths, string(filepath.ListSeparator))
		r := genList(p, c.dir, optType{head:true})
		if !equal(c.result, r) {
			t.Errorf("expectd: %v : %v", c.result, r)
		}
	}
}

func equal[T comparable](x, y []T) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}
