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
		log.Println("real: " + s.GetValue())
		for i, str := range array {
			log.Println("NOW " + strconv.Itoa(i) + ": " + str.GetValue())
		}
	}
}

func LogStr(s string) {
	log.Println(s)
}

func LogInt(i int) {
	log.Println(i)
}
