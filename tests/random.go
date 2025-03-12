package tests

import (
	"encoding/json"
	"math/rand"
)

const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

func randomInt() int {
	return rand.Intn(1000)
}

func RandomString() string {
	str := make([]byte, 36)
	for i := range str {
		str[i] = chars[rand.Intn(len(chars))]
	}
	return string(str)
}

func randomFloat() float64 {
	return rand.Float64() * float64(randomInt())
}

func RandomItem() any {
	dice := rand.Intn(10)

	return generateItem(dice)
}

func generateItem(dice int) any {
	switch dice {
	case 0:
		return randomInt()
	case 1:
		return rand.Float64() * 1000
	case 2:
		return randomInt()%2 == 0
	case 3:
		return RandomString()
	case 4:
		jb, err := json.Marshal(TestingStruct1{AString: RandomString(), AnInt: randomInt()})
		if err != nil {
			return "an error"
		}
		return "[json]" + string(jb)
	case 5:
		return []int{randomInt(), randomInt(), randomInt(), randomInt(), randomInt()}
	case 6:
		return []float64{randomFloat(), randomFloat(), randomFloat(), randomFloat(), randomFloat(), randomFloat()}
	case 7:
		return TestingStruct1{AString: RandomString(), AnInt: randomInt()}
	case 8:
		return TestingStruct2{
			MoreStrings: []string{
				RandomString(), RandomString(), RandomString(),
			},
			MoreInts: []int{
				randomInt(), randomInt(), randomInt(),
			}, MaybeSomeOtherThings: []*TestingStruct2{
				{MoreStrings: []string{RandomString()}},
			}}
	default:
		return chars
	}
}
