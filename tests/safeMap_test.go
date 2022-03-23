package tests

import (
	"testing"

	"github.com/ALiwoto/StrongStringGo/strongStringGo"
)

type TopString string

type dummyStructType[T ~string] struct {
	member1 int
	member2 string
	tMember T
}

var defaultDummy dummyStructType[TopString] = dummyStructType[TopString]{
	member1: 20,
	member2: "20",
	tMember: "20",
}

func TestSafeMap01(t *testing.T) {
	m1 := strongStringGo.GetEmptySafeMap[int, string]()
	m2 := strongStringGo.GetEmptySafeMap[string, string]()
	m3 := strongStringGo.GetEmptySafeMap[string, dummyStructType[TopString]]()

	m1.Set(1, "1")
	m2.Set("1", "1")

	m1.Set(2, "2")
	m2.Set("2", "2")

	v10 := m1.GetValue(10)
	if v10 != "" {
		t.Error("Expected empty string for v10, got:", v10)
		return
	}

	m1.Set(10, "10")
	v10 = m1.GetValue(10)
	if v10 != "10" {
		t.Error("Expected 10 for v10, got:", v10)
		return
	}

	m1.Delete(10)
	v10 = m1.GetValue(10)
	if v10 != "" {
		t.Error("Expected empty string for v10, got:", v10)
		return
	}

	m1.SetDefault("NOT_FOUND")
	m3.SetDefault(defaultDummy)

	v10 = m1.GetValue(10)
	if v10 != "NOT_FOUND" {
		t.Error("Expected NOT_FOUND string for v10, got:", v10)
		return
	}

	m3.Add("first", &dummyStructType[TopString]{1, "1", "1"})
	m3.Add("second", &dummyStructType[TopString]{2, "2", "2"})
	m3.Add("third", &dummyStructType[TopString]{3, "3", "3"})

	if m3.Get("first").member1 != 1 {
		t.Error("Expected 1 for first member1, got:", m3.Get("first").member1)
		return
	}

	if m3.Get("first").member2 != "1" {
		t.Error("Expected 1 for first member2, got:", m3.Get("first").member2)
		return
	}

	if m3.Get("first").tMember != "1" {
		t.Error("Expected 1 for first tMember, got:", m3.Get("first").tMember)
		return
	}

	dummyNotFound := m3.GetValue("something something")
	if dummyNotFound != defaultDummy {
		t.Error("Expected defaultDummy for something something, got:", m3.GetValue("something something"))
		return
	}

	normalMap1 := m1.ToNormalMap()
	if normalMap1 == nil {
		t.Error("Expected not nil for normalMap1, got:", normalMap1)
	}

	normalMap2 := m2.ToNormalMap()
	if normalMap2 == nil {
		t.Error("Expected not nil for normalMap1, got:", normalMap1)
	}

	normalMap3 := m3.ToNormalMap()
	if normalMap3 == nil {
		t.Error("Expected not nil for normalMap1, got:", normalMap1)
	}

	if t.Failed() {
		return
	}
}
