package simulator

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
