package cmd

import (
	"encoding/json"
	"math/rand"
)

const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

type TestingStruct1 struct {
	AString       string
	AnInt         int
	Children      []TestingStruct1
	BrandNewField []float32
}

type TestingStruct2 struct {
	MoreStrings          []string
	MoreInts             []int
	MaybeSomeOtherThings []*TestingStruct2
	PossiblyNil          *TestingStruct1
}

var (
	quantity = 1_000_000
)

func ri() int {
	return rand.Intn(1000)
}

func rs() string {
	str := make([]byte, 36)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}

func rf() float64 {
	return rand.Float64() * float64(ri())
}

func randomItem() any {
	dice := rand.Intn(10)
	switch dice {
	case 0:
		return ri()
	case 1:
		return rand.Float64() * 1000
	case 2:
		return ri()%2 == 0
	case 3:
		return rs()
	case 4:
		jb, err := json.Marshal(TestingStruct1{AString: rs(), AnInt: ri()})
		if err != nil {
			return "an error"
		}
		return "[json]" + string(jb)
	case 5:
		return []int{ri(), ri(), ri(), ri(), ri()}
	case 6:
		return []float64{rf(), rf(), rf(), rf(), rf(), rf()}
	case 7:
		return TestingStruct1{AString: rs(), AnInt: ri()}
	case 8:
		return TestingStruct2{MoreStrings: []string{rs(), rs(), rs()}, MoreInts: []int{ri(), ri(), ri()}, MaybeSomeOtherThings: []*TestingStruct2{&TestingStruct2{MoreStrings: []string{rs()}}}}
	default:
		return chars
	}
}
