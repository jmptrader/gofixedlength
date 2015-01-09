package gofixedlength

import (
	"log"
	"testing"
	"time"
)

// Basic struct, valid
type testStruct1 struct {
	NumberA int    `fixed:"0-5"`
	NumberB int    `fixed:"5-10"`
	StringC string `fixed:"10-15"`
	StringD string `fixed:"15-35"`
}

// Basic struct, fields randomized, valid
type testStruct2 struct {
	StringD string `fixed:"15-35"`
	NumberA int    `fixed:"0-5"`
	StringC string `fixed:"10-15"`
	NumberB int    `fixed:"5-10"`
}

// Struct containing formatted date and float, valid
type testStruct3 struct {
	Uno     string    `fixed:"0-10"`
	Due     int       `fixed:"10-20"`
	Tre     time.Time `fixed:"20-30,2015-10-31"`
	Quattro string    `fixed:"30-40"`
	Cinque  float32   `fixed:"40-50,2"`
}

// Basic struct with holes, valid
type testStruct4 struct {
	NumberA int    `fixed:"0-5"`
	NumberB int    `fixed:"5-10"`
	StringC string `fixed:"10-15"`
	StringD string `fixed:"20-40"` // There is a gap of 5 columns between StringC and StringD (15-20)
}

// Basic struct with overlap. Should only fail if the contested columns are not consistently defined.
type testStruct5 struct {
	NumberA int    `fixed:"0-5"`
	NumberB int    `fixed:"5-10"`
	StringC string `fixed:"10-15"`
	StringD string `fixed:"12-35"` // There is an overlap in pos. 12-15 between StringC and StringD
}

type embeddedStruct struct {
	Another    string `fixed:"50-60"`
	AndFinally string `fixed:"60-62"`
}

// Embedded struct, valid
type testStruct6 struct {
	Embedded1        testStruct1 // Basic struct, valid
	AdditionalField1 int         `fixed:"35-40"`
	AdditionalField2 int         `fixed:"40-42"`
	AdditionalField3 string      `fixed:"42-47"`
	AdditionalField4 string      `fixed:"47-50"`
	embeddedStruct
}

// Struct with invalid index. Should fail.
type testFailingStruct1 struct {
	Uno string `fixed:"0-10"`
	Due int    `fixed:"10-9"`
}

// Struct with invalid index. Should fail.
type testFailingStruct2 struct {
	Uno string `fixed:"0-10"`
	Due int    `fixed:"10-"`
}

// Struct with invalid index. Should fail.
type testFailingStruct3 struct {
	Uno string `fixed:"0-10"`
	Due int    `fixed:"-20"`
}

// Embedded struct with fixed length (shorter than embedded struct's length). Should fail.
type testFailingStruct4 struct {
	EmbeddedStruct   testStruct1 `fixed:"0-30"` // Basic struct, valid
	AdditionalField1 int         `fixed:"30-40"`
	AdditionalField2 int         `fixed:"40-42"`
	AdditionalField3 string      `fixed:"42-47"`
	AdditionalField4 string      `fixed:"47-50"`
}

func TestLineLength(t *testing.T) {
	t.Log("Length calculator test")
	if length := LineLength(testStruct1{}); length != 35 {
		t.Errorf("Failed to find the line length (found %v, expected %v)", length, 35)
	}
	if length := LineLength(testStruct2{}); length != 35 {
		t.Errorf("Failed to find the line length (found %v, expected %v)", length, 35)
	}
	if length := LineLength(testStruct3{}); length != 50 {
		t.Errorf("Failed to find the line length (found %v, expected %v)", length, 50)
	}
	if length := LineLength(testStruct4{}); length != 40 {
		t.Errorf("Failed to find the line length (found %v, expected %v)", length, 40)
	}
	if length := LineLength(testStruct5{}); length != 35 {
		t.Errorf("Failed to find the line length (found %v, expected %v)", length, 35)
	}
	if length := LineLength(testStruct6{}); length != 62 {
		t.Errorf("Failed to find the line length (found %v, expected %v)", length, 62)
	}
}

func TestMarshal(t *testing.T) {
	var out string
	expectedOut := "0012312345ohmy What's happening?  "
	t1 := testStruct1{
		NumberA: 123,                 // `fixed:"0-5"`
		NumberB: 12345,               // `fixed:"5-10"`
		StringC: "ohmy",              // `fixed:"10-15"`
		StringD: "What's happening?", // `fixed:"15-35"`
	}
	err := Marshal(t1, out)
	if err != nil {
		t.Errorf("Error while marshaling the test struct 1: %v\n", err)
	}
	log.Println(out)
	if out != expectedOut {
		t.Errorf("Marshalled string doesn't match the expected output:\n%v\n", out)
	}
}
