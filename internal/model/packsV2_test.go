package model

import (
	"reflect"
	"testing"
)

func TestAddPacksV2(t *testing.T) {
	p := PacksV2{1: 2, 2: 3}
	p.AddPacks(1, 3)

	expected := PacksV2{1: 5, 2: 3}
	if !reflect.DeepEqual(p, expected) {
		t.Errorf("Expected %v, got %v", expected, p)
	}
}

func TestAddPackV2(t *testing.T) {
	p := PacksV2{1: 2, 2: 3}
	p.AddPack(1)

	expected := PacksV2{1: 3, 2: 3}
	if !reflect.DeepEqual(p, expected) {
		t.Errorf("Expected %v, got %v", expected, p)
	}
}

func TestRemovePackV2(t *testing.T) {
	p := PacksV2{
		1: 2,
		2: 3,
	}
	p.RemovePack(1)

	expected := PacksV2{1: 1, 2: 3}
	if !reflect.DeepEqual(p, expected) {
		t.Errorf("Expected %v, got %v", expected, p)
	}
}

func TestSetCountV2(t *testing.T) {
	p := PacksV2{1: 2, 2: 3}
	p.SetCount(1, 0)

	expected := PacksV2{2: 3}
	if !reflect.DeepEqual(p, expected) {
		t.Errorf("Expected %v, got %v", expected, p)
	}
}
