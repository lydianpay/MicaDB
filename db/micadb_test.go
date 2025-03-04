package db

import (
	"fmt"
	"testing"
	"time"
)

func TestNewMicaDB(t *testing.T) {
	filename := fmt.Sprintf("safe-to-delete.%v.unit_test", time.Now())
	mdb, err := NewMicaDB(filename)
	if err != nil || mdb == nil {
		t.Fail()
	}

}
