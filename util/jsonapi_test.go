package util_test

import (
	"github.com/CactusDev/Xerophi/util"
)

// Stuff required to test MarshalResponse
type MarshalResponseTestObj struct {
	JSONAPITagTestField string `jsonapi:"test"`
	RandomTagTestField  string `random:"test"`
}

func (t MarshalResponseTestObj) GetAPITag(lookup string) string {
	return util.FieldTag(t, lookup, "jsonapi")
}

// func TestMarshalResponse(T *testing.T) {
// 	testTable := []struct {
// 		field    string
// 		tag      string
// 		expected string
// 		testObj  MarshalResponseTestObj
// 	}{
// 		{
// 			field:    "JSONAPITagTestField",
// 			tag:      "jsonapi",
// 			expected: "expected",
// 			testObj: MarshalResponseTestObj{
// 				JSONAPITagTestField: "expected",
// 				RandomTagTestField:  "",
// 			},
// 		},
// 		{
// 			field:    "RandomTagTestField",
// 			tag:      "random",
// 			expected: "expectedRandom",
// 			testObj: MarshalResponseTestObj{
// 				JSONAPITagTestField: "expected",
// 				RandomTagTestField:  "expectedRandom",
// 			},
// 		},
// 	}

// 	for _, test := range testTable {
// 		result := util.MarshalResponse(test.testObj)
// 		fmt.Println(test)
// 		fmt.Println(result)
// 	}
// }
