package base

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/internal/client"
)

type nestedCustomTestStruct struct {
	TestInt    int      `json:"test_int"`
	TestBool   bool     `json:"test_bool"`
	TestString string   `json:"test_string"`
	TestList   []string `json:"test_list"`
}

type customTestStruct struct {
	TestKey    string                  `json:"key"`
	TestValue  string                  `json:"test_value"`
	TestNested *nestedCustomTestStruct `json:"test_nested_struct"`
}

func Setup() *Base {
	projectKey := os.Getenv("DETA_SDK_TEST_PROJECT_KEY")
	baseName := os.Getenv("DETA_SDK_TEST_BASE_NAME")
	d, _ := deta.New(deta.WithProjectKey(projectKey))
	return New(d, baseName)
}

func TearDown(b *Base, t *testing.T) {
	var items []map[string]interface{}
	_, err := b.Fetch(&FetchInput{
		Q:    nil,
		Dest: &items,
	})
	if err != nil {
		t.Log("Failed to fetch items in teardown, further tests might fail")
	}
	for _, item := range items {
		key := item["key"].(string)
		err := b.Delete(item["key"].(string))
		if err != nil {
			t.Logf("Failed to delete test item with key '%s'.\nFurther tests might fail", key)
		}
	}
}

func TestModifyItem(t *testing.T) {
	base := Setup()
	defer TearDown(base, t)

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
				},
			},
			modifiedItem: baseItem{
				"key":        "key",
				"test_value": "value",
				"test_nested_struct": map[string]interface{}{
					"test_int":    float64(1),
					"test_bool":   true,
					"test_list":   []interface{}{"a", "b"},
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
	base := Setup()
	defer TearDown(base, t)

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
	base := Setup()
	defer TearDown(base, t)

	// put items
	testItems := []*customTestStruct{
		&customTestStruct{
			TestKey:   "a",
			TestValue: "value",
			TestNested: &nestedCustomTestStruct{
				TestInt:  1,
				TestList: []string{"a", "b"},
				TestBool: true,
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

func TestPutMany(t *testing.T) {
	base := Setup()
	defer TearDown(base, t)

	testItems := []*customTestStruct{
		&customTestStruct{
			TestKey:   "key",
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

	testCases := []struct {
		items []*customTestStruct
		err   error
	}{
		{testItems, nil},
	}

	for _, tc := range testCases {
		_, err := base.PutMany(tc.items)
		if !errors.Is(err, tc.err) {
			t.Errorf("Unexpected error value. Expected: %v Got %v", tc.err, err)
		}
		for _, item := range tc.items {
			var dest customTestStruct
			err = base.Get(item.TestKey, &dest)
			if err != nil {
				t.Fatalf("Failed to get item with key %s", item.TestKey)
			}
			if !reflect.DeepEqual(*item, dest) {
				t.Errorf("Item not equal.\nExpected:\n%v\nStored Items:\n%v", item, dest)
			}
		}
	}
}

func TestUpdate(t *testing.T) {
	base := Setup()
	defer TearDown(base, t)

	testCases := []struct {
		item       *customTestStruct
		updates    Updates
		resultItem *customTestStruct
		err        error
	}{
		{
			item: &customTestStruct{
				TestKey:   "a",
				TestValue: "value",
				TestNested: &nestedCustomTestStruct{
					TestInt:    1,
					TestBool:   true,
					TestList:   []string{"b"},
					TestString: "little",
				},
			},
			updates: Updates{
				"test_value":                     "changed value",
				"test_nested_struct.test_int":    base.Util.Increment(1),
				"test_nested_struct.test_string": base.Util.Trim(),
				"test_nested_struct.test_list":   base.Util.Prepend("c"),
			},
			resultItem: &customTestStruct{
				TestKey:   "a",
				TestValue: "changed value",
				TestNested: &nestedCustomTestStruct{
					TestInt:  2,
					TestBool: true,
					TestList: []string{"c", "b"},
				},
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		key, err := base.Put(tc.item)
		if err != nil {
			t.Fatalf("Failed to put test item with key %s", key)
		}

		err = base.Update(key, tc.updates)
		if !errors.Is(err, tc.err) {
			t.Errorf("Unexpected error value. Expected: %v Got %v", tc.err, err)
		}

		var dest customTestStruct
		err = base.Get(key, &dest)
		if err != nil {
			t.Fatalf("Failed to get test item with key %s", key)
		}

		if !reflect.DeepEqual(*tc.resultItem, dest) {
			t.Errorf("Item not equal.\nExpected:\n%v\nStored Items:\n%v", *tc.resultItem, dest)
		}
	}
}

func TestInsert(t *testing.T) {
	base := Setup()
	defer TearDown(base, t)

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
		{
			item: customTestStruct{
				TestKey:   "key",
				TestValue: "value",
				TestNested: &nestedCustomTestStruct{
					TestInt: 1,
				},
			},
			err: client.ErrConflict,
		},
	}

	for _, tc := range testCases {
		_, err := base.Insert(tc.item)
		if !errors.Is(err, tc.err) {
			t.Errorf("Unexpected error value for item %v\nExpected: %v Got: %v", tc.item, tc.err, err)
		}
		var storedItem customTestStruct
		if tc.err == nil {
			err = base.Get(tc.item.TestKey, &storedItem)
			if err != nil {
				t.Errorf("Failed to get item with key %s", tc.item.TestKey)
			}
			if !reflect.DeepEqual(storedItem, tc.item) {
				t.Errorf("Items not equal.\nExpected:\n%v\nStored:\n%v", tc.item, storedItem)
			}
		}
	}
}

func TestDelete(t *testing.T) {
	base := Setup()
	defer TearDown(base, t)

	// put items
	testItems := []*customTestStruct{
		&customTestStruct{
			TestKey:   "a",
			TestValue: "value",
			TestNested: &nestedCustomTestStruct{
				TestInt:  1,
				TestList: []string{"a", "b"},
				TestBool: true,
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
		key string
		err error
	}

	var testCases []*testCase
	for n := range testItems {
		testCases = append(testCases, &testCase{
			key: keys[n],
			err: nil,
		})
	}

	for n, tc := range testCases {
		key := keys[n]
		err := base.Delete(key)
		if !errors.Is(err, tc.err) {
			t.Errorf("Unexpected error value. Expected: %v Got %v", tc.err, err)
		}
		var dest customTestStruct
		err = base.Get(key, &dest)
		if !errors.Is(err, client.ErrNotFound) {
			t.Errorf("Item with key %s not deleted from database", key)
		}
	}
}

func TestFetch(t *testing.T) {
	base := Setup()
	defer TearDown(base, t)

	// put items
	testItems := []*customTestStruct{
		&customTestStruct{
			TestKey:   "a",
			TestValue: "value",
			TestNested: &nestedCustomTestStruct{
				TestInt:  1,
				TestList: []string{"a", "b"},
				TestBool: true,
			},
		},
		&customTestStruct{
			TestKey:   "b",
			TestValue: "value",
		},
	}
	_, err := base.PutMany(testItems)
	if err != nil {
		t.Fatalf("Failed to put items with error %v", err)
	}

	type testCase struct {
		query         Query
		expectedItems []customTestStruct
		dest          []customTestStruct
		err           error
	}

	testCases := []testCase{
		{
			query: Query{
				{"test_nested_struct.test_int?lte": 1},
			},
			expectedItems: []customTestStruct{
				customTestStruct{
					TestKey:   "a",
					TestValue: "value",
					TestNested: &nestedCustomTestStruct{
						TestInt:  1,
						TestList: []string{"a", "b"},
						TestBool: true,
					},
				},
			},
			dest: []customTestStruct{},
			err:  nil,
		},
		{
			query: Query{
				{"test_value": "value"},
			},
			expectedItems: []customTestStruct{
				customTestStruct{
					TestKey:   "a",
					TestValue: "value",
					TestNested: &nestedCustomTestStruct{
						TestInt:  1,
						TestList: []string{"a", "b"},
						TestBool: true,
					},
				},
				customTestStruct{
					TestKey:   "b",
					TestValue: "value",
				},
			},
			dest: []customTestStruct{},
			err:  nil,
		},
	}

	for _, tc := range testCases {
		_, err := base.Fetch(&FetchInput{
			Q:    tc.query,
			Dest: &tc.dest,
		})
		if !errors.Is(err, tc.err) {
			t.Errorf("Unexpected error value. Expected: %v Got %v", tc.err, err)
		}
		if !reflect.DeepEqual(tc.expectedItems, tc.dest) {
			t.Errorf("Items not equal.\nExpected:\n%v\nGot:\n%v", tc.expectedItems, tc.dest)
		}
	}
}

func TestFetchPaginated(t *testing.T) {
	base := Setup()
	defer TearDown(base, t)

	// put items
	testItems := []customTestStruct{
		customTestStruct{
			TestKey:   "a",
			TestValue: "a",
		},
		customTestStruct{
			TestKey:   "b",
			TestValue: "b",
		},
	}
	_, err := base.PutMany(testItems)
	if err != nil {
		t.Fatalf("Failed to put items with error %v", err)
	}

	var dest []customTestStruct
	lastKey, err := base.Fetch(&FetchInput{
		Q:     nil,
		Dest:  &dest,
		Limit: 1,
	})
	if err != nil {
		t.Fatalf("Failed to fetch items with err %v", err)
	}

	key := testItems[0].TestKey
	if lastKey != key {
		t.Errorf("Failed to get correct last key. Expected: %s Got: %s", key, lastKey)
	}

	if !reflect.DeepEqual(testItems[0], dest[0]) {
		t.Errorf("Fetched item not equal to expected.\nExpected:\n%v\nGot: %v", testItems[0], dest[0])
	}

	lastKey, err = base.Fetch(&FetchInput{
		Q:       nil,
		Dest:    &dest,
		LastKey: lastKey,
	})
	if err != nil {
		t.Errorf("Failed to fetch items with err %v", err)
	}

	if lastKey != "" {
		t.Errorf("Failed to get correct last key. Expected: %s Got: %s", key, lastKey)
	}
	if !reflect.DeepEqual(testItems[1], dest[0]) {
		t.Errorf("Fetched item not equal to expected.\nExpected:\n%v\nGot: %v", testItems[1], dest[0])
	}
}
