// StrongStringGo Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package tests

import (
	"log"
	"strconv"
	"testing"

	"github.com/ALiwoto/StrongStringGo/strongStringGo"
)

func TestStrong(t *testing.T) {
	LogStr("Hi")
	LogInt(5)
	s := strongStringGo.Qss("hello!; how are you? () are you okay?")
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
	s := strongStringGo.Qss("true")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = strongStringGo.Qss("false")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = strongStringGo.Qss("TRUE")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = strongStringGo.Qss("FALSE")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = strongStringGo.Qss("True")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = strongStringGo.Qss("False")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = strongStringGo.Qss("on")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = strongStringGo.Qss("off")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = strongStringGo.Qss("ON")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = strongStringGo.Qss("OFF")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = strongStringGo.Qss("yes")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = strongStringGo.Qss("no")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}

	s = strongStringGo.Qss("YES")
	if !s.ToBool() {
		t.Error("Expected true, got false")
		return
	}

	s = strongStringGo.Qss("NO")
	if s.ToBool() {
		t.Error("Expected false, got true")
		return
	}
}

func LogStr(s string) {
	log.Println(s)
}

func LogInt(i int) {
	log.Println(i)
}
