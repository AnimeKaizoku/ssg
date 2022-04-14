// ssg Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package tests

import (
	"log"
	"strconv"
	"testing"

	ws "github.com/AnimeKaizoku/ssg/ssg"
)

func TestTitleCase(t *testing.T) {
	const (
		str1 = "string1"
		str2 = "thisIsString2"
		str3 = "HelloThere"
	)

	tmp := ws.Title(str1)
	if tmp != "String1" {
		t.Errorf("Expected %s, got %s", "String1", tmp)
		return
	}

	tmp = ws.Title(str2)
	if tmp != "ThisIsString2" {
		t.Errorf("Expected %s, got %s", "ThisIsString2", tmp)
		return
	}

	tmp = ws.Title(str3)
	if tmp != "HelloThere" {
		t.Errorf("Expected %s, got %s", "HelloThere", tmp)
		return
	}
}

func TestStrong(t *testing.T) {
	LogStr("Hi")
	LogInt(5)
	s := ws.Qss("hello!; how are you? () are you okay?")
	if s == nil {
		t.FailNow()
	} else {
		array := s.SplitStr("; ", "() ")
		LogStr("real: " + s.GetValue())
		for i, str := range array {
			LogStr("NOW " + strconv.Itoa(i) + ": " + str.GetValue())
		}
	}
}

func TestToBool(t *testing.T) {
	s := ws.Qss("true")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = ws.Qss("false")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = ws.Qss("TRUE")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = ws.Qss("FALSE")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = ws.Qss("True")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = ws.Qss("False")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = ws.Qss("on")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = ws.Qss("off")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = ws.Qss("ON")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = ws.Qss("OFF")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = ws.Qss("yes")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = ws.Qss("no")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = ws.Qss("YES")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = ws.Qss("NO")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}
}

func TestIntegerHelpers(t *testing.T) {
	const (
		i1       = int64(1)
		i294     = int64(294)
		i356     = int64(356)
		i487     = int64(487)
		i5900    = int64(5900)
		i0x98760 = int64(0x98760)
	)

	s := ws.ToBase10(i1)
	i := ws.ToInt64(s)
	if ws.ToInt64(s) != i1 {
		t.Errorf("Expected %d, got %d", i1, i)
		return
	}

	s = ws.ToBase10(i294)
	i = ws.ToInt64(s)
	if ws.ToInt64(s) != i294 {
		t.Errorf("Expected %d, got %d", i294, i)
		return
	}

	s = ws.ToBase10(i356)
	i = ws.ToInt64(s)
	if ws.ToInt64(s) != i356 {
		t.Errorf("Expected %d, got %d", i356, i)
		return
	}

	s = ws.ToBase10(i487)
	i = ws.ToInt64(s)
	if ws.ToInt64(s) != i487 {
		t.Errorf("Expected %d, got %d", i487, i)
		return
	}

	s = ws.ToBase10(i5900)
	i = ws.ToInt64(s)
	if ws.ToInt64(s) != i5900 {
		t.Errorf("Expected %d, got %d", i5900, i)
		return
	}

	s = ws.ToBase10(i0x98760)
	i = ws.ToInt64(s)
	if ws.ToInt64(s) != i0x98760 {
		t.Errorf("Expected %d, got %d", i0x98760, i)
		return
	}
}

func LogStr(s string) {
	log.Println(s)
}

func LogInt(i int) {
	log.Println(i)
}
