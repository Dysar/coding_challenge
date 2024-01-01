package model

import (
	"reflect"
	"testing"
)

func TestAddPacks(t *testing.T) {
	p := Packs{1: 2, 2: 3}
	p.AddPacks(1, 3)

	expected := Packs{1: 5, 2: 3}
	if !reflect.DeepEqual(p, expected) {
		t.Errorf("Expected %v, got %v", expected, p)
	}
}

func TestAddPack(t *testing.T) {
	p := Packs{1: 2, 2: 3}
	p.AddPack(1)

	expected := Packs{1: 3, 2: 3}
	if !reflect.DeepEqual(p, expected) {
		t.Errorf("Expected %v, got %v", expected, p)
	}
}

func TestRemovePack(t *testing.T) {
	p := Packs{
		1: 2,
		2: 3,
	}
	p.RemovePack(1)

	expected := Packs{1: 1, 2: 3}
	if !reflect.DeepEqual(p, expected) {
		t.Errorf("Expected %v, got %v", expected, p)
	}
}

func TestSetCount(t *testing.T) {
	p := Packs{1: 2, 2: 3}
	p.SetCount(1, 0)

	expected := Packs{2: 3}
	if !reflect.DeepEqual(p, expected) {
		t.Errorf("Expected %v, got %v", expected, p)
	}
}
