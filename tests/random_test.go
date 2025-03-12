package tests

import (
	"reflect"
	"testing"
)

func TestRandomInt(t *testing.T) {
	val := randomInt()

	if reflect.TypeOf(val).Kind() != reflect.Int {
		t.Errorf("randomInt should be int, but received %v", val)
	}
}

func TestRandomString(t *testing.T) {
	val := RandomString()
	if reflect.TypeOf(val).Kind() != reflect.String {
		t.Errorf("randomString should be string, but received %v", val)
	}
}

func TestRandomFloat(t *testing.T) {
	val := randomFloat()
	if reflect.TypeOf(val).Kind() != reflect.Float64 {
		t.Errorf("randomFloat should be float64, but received %v", val)
	}
}

func TestRandomItem(t *testing.T) {
	val := RandomItem()
	if val == nil {
		t.Errorf("randomItem should not be nil")
	}
}

func TestGenerateItem(t *testing.T) {
	val := generateItem(0)
	if reflect.TypeOf(val).Kind() != reflect.Int {
		t.Errorf("randomItem should be int, but received %v", val)
	}

	val = generateItem(1)
	if reflect.TypeOf(val).Kind() != reflect.Float64 {
		t.Errorf("randomItem should be float64, but received %v", val)
	}

	val = generateItem(2)
	if reflect.TypeOf(val).Kind() != reflect.Bool {
		t.Errorf("randomItem should be bool, but received %v", val)
	}

	val = generateItem(3)
	if reflect.TypeOf(val).Kind() != reflect.String {
		t.Errorf("randomItem should be string, but received %v", val)
	}

	val = generateItem(4)
	if reflect.TypeOf(val).Kind() != reflect.String {
		t.Errorf("randomItem should be string, but received %v", val)
	}

	values := generateItem(5)
	for _, num := range values.([]int) {
		if reflect.TypeOf(num).Kind() != reflect.Int {
			t.Errorf("randomItem should be int, but received %v", val)
		}
	}

	values = generateItem(6)
	for _, num := range values.([]float64) {
		if reflect.TypeOf(num).Kind() != reflect.Float64 {
			t.Errorf("randomItem should be float64, but received %v", val)
		}
	}

	val = generateItem(7)
	if reflect.TypeOf(val).Kind() != reflect.Struct {
		t.Errorf("randomItem should be struct, but received %v", val)
	} else {
		if reflect.TypeOf(val).Name() != "TestingStruct1" {
			t.Errorf("randomItem should be TestingStruct1, but received %v", val)
		}
	}

	val = generateItem(8)
	if reflect.TypeOf(val).Kind() != reflect.Struct {
		t.Errorf("randomItem should be struct, but received %v", val)
	} else {
		if reflect.TypeOf(val).Name() != "TestingStruct2" {
			t.Errorf("randomItem should be TestingStruct2, but received %v", val)
		}
	}

	val = generateItem(9)
	if reflect.TypeOf(val).Kind() != reflect.String && val != chars {
		t.Errorf("randomItem should be string of chars, but received %v", val)
	}

}
