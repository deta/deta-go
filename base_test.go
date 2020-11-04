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

func Setup(t *testing.T) *Base {
	projectKey := os.Getenv("TEST_PROJECT_KEY")
	baseName := os.Getenv("TEST_BASE_NAME")
	rootEndpoint := os.Getenv("TEST_ENDPOINT")
	base, err := newBase(projectKey, baseName, rootEndpoint)
	if err != nil {
		t.Fatalf("Failed to set up test base: %v", err)
	}
	return base
}

func TestModifyItem(t *testing.T) {
	base := Setup(t)

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

func TestPut(t *testing.T) {
	base := Setup(t)
	testCases := []struct {
		item customTestStruct
		err  error
	}{
		{
			item: customTestStruct{
				TestKey:   "key",
				TestValue: "value",
				TestNested: &nestedCustomTestStruct{
					TestInt: 1,
				},
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		_, err := base.Put(tc.item)
		if !errors.Is(err, tc.err) {
			t.Errorf("Unexpected error value for item %v\n. Expected %v got %v", tc.item, tc.err, err)
		}
		var storedItem customTestStruct
		err = base.Get(tc.item.TestKey, &storedItem)
		if err != nil {
			t.Errorf("Failed to get item with key %s", tc.item.TestKey)
		}
		if !reflect.DeepEqual(storedItem, tc.item) {
			t.Errorf("Items not equal.\nExpected:\n%v\nStored:\n%v", tc.item, storedItem)
		}
	}
}

func TestGet(t *testing.T) {
	base := Setup(t)

	// put items
	testItems := []*customTestStruct{
		&customTestStruct{
			TestKey:   "a",
			TestValue: "value",
			TestNested: &nestedCustomTestStruct{
				TestInt: 1,
			},
		},
		&customTestStruct{
			TestKey:   "b",
			TestValue: "value",
		},
	}

	var keys []string
	for _, item := range testItems {
		key, err := base.Put(item)
		if err != nil {
			t.Fatalf("Failed to put item %v with error %v", item, err)
		}
		keys = append(keys, key)
	}

	type testCase struct {
		key  string
		item customTestStruct
		dest customTestStruct
		err  error
	}

	var testCases []*testCase
	for n, item := range testItems {
		testCases = append(testCases, &testCase{
			key:  keys[n],
			item: *item,
			dest: customTestStruct{},
			err:  nil,
		})
	}

	for _, tc := range testCases {
		err := base.Get(tc.key, &tc.dest)
		if !errors.Is(err, tc.err) {
			t.Errorf("Unexpected error value. Expected: %v Got %v", tc.err, err)
		}
		if !reflect.DeepEqual(tc.item, tc.dest) {
			t.Errorf("Items not equal.\nExpected:\n%v\nGot:\n%v", tc.item, tc.dest)
		}
	}

}
