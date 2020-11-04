package deta

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

type nestedCustomTestStruct struct {
	TestInt    int               `json:"test_int"`
	TestBool   bool              `json:"test_bool"`
	TestString string            `json:"test_string"`
	TestMap    map[string]string `json:"test_map"`
	TestList   []string          `json:"test_list"`
}

type customTestStruct struct {
	TestKey    string                  `json:"key"`
	TestValue  string                  `json:"test_value"`
	TestNested *nestedCustomTestStruct `json:"test_nested_struct"`
}

func SetUp() (*Base, error) {
	projectKey := os.Getenv("TEST_PROJECT_KEY")
	baseName := os.Getenv("TEST_BASE_NAME")
	rootEndpoint := os.Getenv("TEST_ENDPOINT")
	return newBase(projectKey, baseName, rootEndpoint)
}

func TestModifyItem(t *testing.T) {

	base, err := SetUp()
	if err != nil {
		t.Fatalf("Failed to set up test base: %v", err)
	}

	testStructCases := []struct {
		item         customTestStruct
		modifiedItem baseItem
	}{
		{
			item: customTestStruct{
				TestKey:   "key",
				TestValue: "value",
				TestNested: &nestedCustomTestStruct{
					TestInt:    1,
					TestBool:   true,
					TestList:   []string{"a", "b"},
					TestString: "test",
					TestMap:    map[string]string{"a": "b"},
				},
			},
			modifiedItem: baseItem{
				"key":        "key",
				"test_value": "value",
				"test_nested_struct": map[string]interface{}{
					"test_int":    float64(1),
					"test_bool":   true,
					"test_list":   []interface{}{"a", "b"},
					"test_map":    map[string]interface{}{"a": "b"},
					"test_string": "test",
				},
			},
		},
	}

	testMapCases := []struct {
		item         map[string]interface{}
		modifiedItem baseItem
	}{
		{
			item: map[string]interface{}{
				"key":  "abcd",
				"name": "test",
			},
			modifiedItem: baseItem{
				"key":  "abcd",
				"name": "test",
			},
		},
	}

	testBadCases := []struct {
		item interface{}
		err  error
	}{
		{"a string", ErrBadItem},
		{1, ErrBadItem},
		{true, ErrBadItem},
		{[]string{"a", "b"}, ErrBadItem},
	}

	for _, tc := range testStructCases {
		o, _ := base.modifyItem(tc.item)
		if !reflect.DeepEqual(o, tc.modifiedItem) {
			t.Errorf("Failed to modify struct.\nExpected: %v\nGot: %v", tc.modifiedItem, o)
		}
	}
	for _, tc := range testMapCases {
		o, _ := base.modifyItem(tc.item)
		if !reflect.DeepEqual(o, tc.modifiedItem) {
			t.Errorf("Failed to modify map.\nExpected: %v\nGot: %v", tc.modifiedItem, o)
		}
	}
	for _, tc := range testBadCases {
		o, err := base.modifyItem(tc.item)
		if !errors.Is(err, tc.err) {
			t.Errorf("Unexpected error value for item %v.\nExpected: %v\nGot Item: %v\nGot Error:%v", tc.item, tc.err, o, err)
		}
	}
}
