package tests

import (
	"testing"

	"github.com/ALiwoto/StrongStringGo/strongStringGo"
)

func TestList01(t *testing.T) {
	l1 := strongStringGo.GetEmptyList[string]()
	l1.Add("1", "2", "3", "4")
	arr := l1.AsArray()
	if len(arr) != 4 {
		t.Error("Expected 4, got:", len(arr))
		return
	}

	l1.Remove("3")
	arr = l1.AsArray()
	if len(arr) != 3 {
		t.Error("Expected 4, got:", len(arr))
		return
	}

	l1.Add("1")

	countV1 := l1.Count("1")
	if countV1 != 2 {
		t.Error("Expected 2 for countV1, got:", countV1)
		return
	}

}
