package deta

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

const (
	readChunkSize = 1024 * 1024 * 10
)

func SetupDrive() *Drive {
	projectKey := os.Getenv("DETA_SDK_TEST_PROJECT_KEY")
	driveName := os.Getenv("DETA_SDK_TEST_DRIVE_NAME")
	rootEndpoint := os.Getenv("DETA_SDK_TEST_ENDPOINT")
	return newDrive(projectKey, driveName, rootEndpoint)
}

func TearDownDrive(d *Drive, t *testing.T) {
	lr, err := d.List(1000, "", "")
	if !errors.Is(err, nil) {
		t.Log("Failed to list names in teardown, further tests might fail")
	}
	for _, name := range lr.Names {
		_, err = d.Delete(name)
		if !errors.Is(err, nil) {
			t.Logf("Failed to delete test file with name '%s'.\nFurther tests might fail", name)
		}
	}

}

func TestPutDrive(t *testing.T) {
	drive := SetupDrive()
	defer TearDownDrive(drive, t)

	testCases := []struct {
		Name        string
		Body        io.Reader
		ContentType string
		Content     string
	}{
		{
			Name:        "test_file_1.txt",
			Body:        strings.NewReader("this is a string."),
			ContentType: "text/plain",
			Content:     "this is a string.",
		},
		{
			Name:        "name with spaces.txt",
			Body:        strings.NewReader("lorem ipsum"),
			ContentType: "text/plain",
			Content:     "lorem ipsum",
		},
		{
			Name:        "test_file_1.txt",
			Body:        strings.NewReader("same file name should be overwritten"),
			ContentType: "text/plain",
			Content:     "same file name should be overwritten",
		},
	}

	for _, tc := range testCases {
		name, err := drive.Put(&PutInput{
			Name:        tc.Name,
			Body:        tc.Body,
			ContentType: tc.ContentType,
		})

		if !errors.Is(err, nil) {
			t.Errorf("Unexpected error value. Expected %v Got %v", nil, err)
		}

		if !reflect.DeepEqual(tc.Name, name) {
			t.Errorf("File names are not equal. \nExpected \n%vGot:\n%v", tc.Name, name)
		}

		driveContent, err := drive.Get(tc.Name)
		if !errors.Is(err, nil) {
			t.Errorf("Unexpected error while trying to get file content.")
		}
		defer driveContent.Close()

		if b, err := ioutil.ReadAll(driveContent); err == nil {
			if !reflect.DeepEqual(tc.Content, string(b)) {
				t.Errorf("Fetched content not equal to expected. \nExpected:\n%v\nGot:\n%v", tc.Content, string(b))
			}
		}

	}

}

func TestGetDrive(t *testing.T) {
	drive := SetupDrive()
	defer TearDownDrive(drive, t)

	testCases := []struct {
		Name        string
		Body        io.Reader
		ContentType string
		Content     string
	}{
		{
			Name:        "string_stream.txt",
			Body:        strings.NewReader("string stream"),
			ContentType: "text/plain",
			Content:     "string stream",
		},
	}

	for _, tc := range testCases {
		_, err := drive.Put(&PutInput{
			Name:        tc.Name,
			Body:        tc.Body,
			ContentType: tc.ContentType,
		})

		if !errors.Is(err, nil) {
			t.Errorf("Unexpected error value. Expected %v Got %v", nil, err)
		}

		driveContent, err := drive.Get(tc.Name)
		if !errors.Is(err, nil) {
			t.Errorf("Unexpected error while trying to get file content.")
		}
		defer driveContent.Close()

		if b, err := ioutil.ReadAll(driveContent); err == nil {
			if !reflect.DeepEqual(tc.Content, string(b)) {
				t.Errorf("Fetched content not equal to expected. \nExpected:\n%v\nGot:\n%v", tc.Content, string(b))
			}
		}
	}
}

func TestDeleteDrive(t *testing.T) {
	drive := SetupDrive()
	defer TearDownDrive(drive, t)

	testCases := []struct {
		Name        string
		Body        io.Reader
		ContentType string
		Content     string
	}{
		{
			Name:        "to_del_1.txt",
			Body:        strings.NewReader("hello"),
			ContentType: "text/plain",
			Content:     "hello",
		},
		{
			Name:        "to del name with spaces.txt",
			Body:        strings.NewReader("hola"),
			ContentType: "text/plain",
			Content:     "hola",
		},
	}

	for _, tc := range testCases {
		name, err := drive.Put(&PutInput{
			Name:        tc.Name,
			Body:        tc.Body,
			ContentType: tc.ContentType,
		})
		if !errors.Is(err, nil) {
			t.Fatalf("Failed to put file %v with error %v", name, err)
		}

		_, err = drive.Delete(name)
		if !errors.Is(err, nil) {
			t.Errorf("Unexpected error value. Expected: %v Got %v", nil, err)
		}
	}
}

func contains(slice []string, name string) bool {
	for _, v := range slice {
		if v == name {
			return true
		}
	}
	return false
}

func TestDeleteManyDrive(t *testing.T) {
	drive := SetupDrive()
	defer TearDownDrive(drive, t)

	testCases := []struct {
		Name        string
		Body        io.Reader
		ContentType string
		Content     string
	}{
		{
			Name:        "to_del_1.txt",
			Body:        strings.NewReader("hello"),
			ContentType: "text/plain",
			Content:     "hello",
		},
		{
			Name:        "to del name with spaces.txt",
			Body:        strings.NewReader("hola"),
			ContentType: "text/plain",
			Content:     "hola",
		},
	}

	var names []string
	for _, tc := range testCases {
		name, err := drive.Put(&PutInput{
			Name:        tc.Name,
			Body:        tc.Body,
			ContentType: tc.ContentType,
		})
		if !errors.Is(err, nil) {
			t.Errorf("Failed to put file %v with error %v", name, err)
		}
		names = append(names, name)
	}

	dr, err := drive.DeleteMany(names)
	if !errors.Is(err, nil) {
		t.Errorf("Failed to delete files %v with error %v", strings.Join(names, ", "), err)
	}

	if !reflect.DeepEqual(true, contains(names, testCases[0].Name)) {
		t.Errorf("Failed to delete file %v with error %v", testCases[0].Name, dr.Failed[testCases[0].Name])
	}

	if !reflect.DeepEqual(true, contains(names, testCases[1].Name)) {
		t.Errorf("Failed to delete file %v with error %v", testCases[1].Name, dr.Failed[testCases[1].Name])
	}

	lr, err := drive.List(1000, "", "")
	if !errors.Is(err, nil) {
		t.Errorf("Failed to fetch file names from drive with the err %v", err)
	}

	if !reflect.DeepEqual(0, len(lr.Names)) {
		t.Errorf("Failed to delete files %v with error %v", strings.Join(names, ", "), err)
	}
}

func TestListDrive(t *testing.T) {
	drive := SetupDrive()
	defer TearDownDrive(drive, t)

	testCases := []struct {
		Name        string
		Body        io.Reader
		ContentType string
		Content     string
	}{
		{
			Name:        "a",
			Body:        strings.NewReader("a"),
			ContentType: "text/plain",
			Content:     "a",
		},
		{
			Name:        "b",
			Body:        strings.NewReader("b"),
			ContentType: "text/plain",
			Content:     "b",
		},
		{
			Name:        "c/d",
			Body:        strings.NewReader("c/d"),
			ContentType: "text/plain",
			Content:     "c and d",
		},
	}

	for _, tc := range testCases {
		name, err := drive.Put(&PutInput{
			Name:        tc.Name,
			Body:        tc.Body,
			ContentType: tc.ContentType,
		})

		if !errors.Is(err, nil) {
			t.Errorf("Failed to put file %v with error %v", name, err)
		}
	}

	lr, err := drive.List(1000, "", "")
	if !errors.Is(err, nil) {
		t.Errorf("Failed to fetch filenames with error %v", err)
	}
	if !reflect.DeepEqual([]string{"a", "b", "c/d"}, lr.Names) {
		t.Errorf("Fetched list not equal expected.\nExpected:\n%v\nGot:\n%v", []string{"a", "b", "c/d"}, lr.Names)
	}

	lr, err = drive.List(1, "", "")
	if !errors.Is(err, nil) {
		t.Errorf("Failed to fetch filenames with error %v", err)
	}
	if !reflect.DeepEqual([]string{"a"}, lr.Names) {
		t.Errorf("Fetched list not equal expected.\nExpected:\n%v\nGot:\n%v", []string{"a"}, lr.Names)
	}

	lr, err = drive.List(2, "", "")
	if !errors.Is(err, nil) {
		t.Errorf("Failed to fetch filenames with error %v", err)
	}
	if !reflect.DeepEqual("b", *lr.Paging.Last) {
		t.Errorf("Fetched list not equal expected.\nExpected:\n%v\nGot:\n%v", "b", *lr.Paging.Last)
	}

	lr, err = drive.List(1000, "c/", "")
	if !errors.Is(err, nil) {
		t.Errorf("Failed to fetch filenames with error %v", err)
	}
	if !reflect.DeepEqual([]string{"c/d"}, lr.Names) {
		t.Errorf("Fetched list not equal expected.\nExpected:\n%v\nGot:\n%v", []string{"c/d"}, lr.Names)
	}

}

func TestPutLargeFile(t *testing.T) {
	drive := SetupDrive()
	defer TearDownDrive(drive, t)

	// create a large random byte chunk
	b := make([]byte, readChunkSize*2+1000)
	_, err := rand.Read(b)
	if !errors.Is(err, nil) {
		t.Errorf("Failed to generate random large file with error %v", err)
	}

	name, err := drive.Put(&PutInput{
		Name: "large_binary_file",
		Body: bytes.NewReader(b),
	})
	if !errors.Is(err, nil) {
		t.Errorf("Failed to put random large file %v with error %v", name, err)
	}

	driveContent, err := drive.Get(name)
	if !errors.Is(err, nil) {
		t.Errorf("Unexpected error while trying to get file content.")
	}
	defer driveContent.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(driveContent)

	if !reflect.DeepEqual(b, buf.Bytes()) {
		t.Errorf("Fetched content not equal to expected.")
	}

}
