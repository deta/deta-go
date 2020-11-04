package deta

import (
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
	TestKey    string                 `json:"key"`
	TestValue  string                 `json:"test_value"`
	TestNested nestedCustomTestStruct `json:"test_nested_struct"`
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
		modifiedItem customTestStruct
	}{
		{
			item: customTestStruct{
				TestKey:   "key",
				TestValue: "value",
				TestNested: nestedCustomTestStruct{
					TestInt:  1,
					TestBool: true,
					TestList: []string{"a", "b"},
					TestMap:  map[string]string{"a": "b"},
				},
			},
			modifiedItem: customTestStruct{
				TestKey:   "key",
				TestValue: "value",
				TestNested: nestedCustomTestStruct{
					TestInt:  1,
					TestBool: true,
					TestList: []string{"a", "b"},
					TestMap:  map[string]string{"a": "b"},
				},
			},
		},
	}

	testMapCases := []struct {
		item         map[string]interface{}
		modifiedItem map[string]interface{}
	}{
		{
			item: map[string]interface{}{
				"key":  "abcd",
				"name": "test",
			},
			modifiedItem: map[string]interface{}{
				"key":  "abcd",
				"name": "test",
			},
		},
	}

	testNativeCases := []struct {
		item         interface{}
		modifiedItem interface{}
	}{
		{"a string", map[string]interface{}{"value": "a string"}},
		{1, map[string]interface{}{"value": 1}},
		{true, map[string]interface{}{"value": true}},
		{[]string{"a", "b"}, map[string]interface{}{"value": []string{"a", "b"}}},
	}

	for _, tc := range testStructCases {
		o := base.modifyItem(tc.item)
		if !reflect.DeepEqual(o, tc.modifiedItem) {
			t.Errorf("Failed to modify struct.\nExpected: %v\nGot: %v", tc.modifiedItem, o)
		}
	}
	for _, tc := range testMapCases {
		o := base.modifyItem(tc.item)
		if !reflect.DeepEqual(o, tc.modifiedItem) {
			t.Errorf("Failed to modify map.\nExpected: %v\nGot: %v", tc.modifiedItem, o)
		}
	}
	for _, tc := range testNativeCases {
		o := base.modifyItem(tc.item)
		if !reflect.DeepEqual(o, tc.modifiedItem) {
			t.Errorf("Failed to modify native.\nExpected: %v\nGot: %v", tc.modifiedItem, o)
		}
	}
}
