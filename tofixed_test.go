package gofixedlength

import (
	"testing"
	"time"
)

// Basic struct, valid
type testStruct1 struct {
	NumberA        int    `fixed:"0-5"`
	NumberB        int    `fixed:"5-10"`
	StringC        string `fixed:"10-15"`
	StringD        string `fixed:"15-35"`
	Length         int
	ExpectedResult string
}

// Basic struct, fields randomized, valid
type testStruct2 struct {
	StringD        string `fixed:"15-35"`
	NumberA        int    `fixed:"0-5"`
	StringC        string `fixed:"10-15"`
	NumberB        int    `fixed:"5-10"`
	Length         int
	ExpectedResult string
}

// Struct containing formatted date and float, valid
type testStruct3 struct {
	Uno            string    `fixed:"0-10"`
	Due            int       `fixed:"10-20"`
	Tre            time.Time `fixed:"20-30,2006-01-02"`
	Quattro        string    `fixed:"30-40"`
	Cinque         float32   `fixed:"40-50,2"`
	Length         int
	ExpectedResult string
}

// Basic struct with holes, valid
type testStruct4 struct {
	NumberA        int    `fixed:"0-5"`
	NumberB        int    `fixed:"5-10"`
	StringC        string `fixed:"10-15"`
	StringD        string `fixed:"20-40"` // There is a gap of 5 columns between StringC and StringD (15-20)
	Length         int
	ExpectedResult string
}

// Basic struct with overlap. Should only fail if the contested columns are not consistently defined.
type testStruct5 struct {
	NumberA        int    `fixed:"0-5"`
	NumberB        int    `fixed:"5-10"`
	StringC        string `fixed:"10-15"`
	StringD        string `fixed:"12-35"` // There is an overlap in pos. 12-15 between StringC and StringD
	Length         int
	ExpectedResult string
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
	Length         int
	ExpectedResult string
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
	t.Log("Marshaller test")
	t1 := testStruct1{
		NumberA:        123,                                   // `fixed:"0-5"`
		NumberB:        12345,                                 // `fixed:"5-10"`
		StringC:        "ohmy",                                // `fixed:"10-15"`
		StringD:        "What's happening?",                   // `fixed:"15-35"`
		ExpectedResult: "0012312345ohmy What's happening?   ", // no tag
	}
	out1, err := Marshal(t1)
	if err != nil {
		t.Errorf("Error while marshaling the test struct 1: %v\n", err)
	}
	if out1 != t1.ExpectedResult {
		t.Errorf("Marshalled string doesn't match the expected output:\n'%v'\n", out1)
	}

	t2 := testStruct2{
		StringD:        "Some more text 12345", // `fixed:"15-35"`
		NumberA:        12345,                  // `fixed:"0-5"`
		StringC:        "AND",                  // `fixed:"10-15"`
		NumberB:        67890,                  // `fixed:"5-10"`
		Length:         35,
		ExpectedResult: "1234567890AND  Some more text 12345",
	}
	out2, err := Marshal(t2)
	if err != nil {
		t.Errorf("Error while marshaling the test struct 2: %v\n", err)
	}
	if out2 != t2.ExpectedResult {
		t.Errorf("Marshalled string doesn't match the expected output:\n'%v'\n", out2)
	}

	date3, _ := time.Parse("2006-01-02", "1984-10-31")
	t3 := testStruct3{
		Uno:            "just atext",                                         // string    `fixed:"0-10"`
		Due:            1234567890,                                           // int       `fixed:"10-20"`
		Tre:            date3,                                                // time.Time `fixed:"20-30,2015-10-31"`
		Quattro:        "anotherStr",                                         // string    `fixed:"30-40"`
		Cinque:         321.458,                                              // float32   `fixed:"40-50,2"`
		Length:         50,                                                   // int
		ExpectedResult: "just atext12345678901984-10-31anotherStr0000321.46", // string
	}
	out3, err := Marshal(t3)
	if err != nil {
		t.Errorf("Error while marshaling the test struct 3: %v\n", err)
	}
	if out3 != t3.ExpectedResult {
		t.Errorf("Marshalled string doesn't match the expected output:\n'%v'\n", out3)
	}

	t4 := testStruct4{
		NumberA:        12345,                                      // int    `fixed:"0-5"`
		NumberB:        67890,                                      // int    `fixed:"5-10"`
		StringC:        "short",                                    // string `fixed:"10-15"`
		StringD:        "and then some       ",                     // string `fixed:"20-40"` // There is a gap of 5 columns between StringC and StringD (15-20)
		Length:         40,                                         // int
		ExpectedResult: "1234567890short     and then some       ", // string
	}
	out4, err := Marshal(t4)
	if err != nil {
		t.Errorf("Error while marshaling the test struct 4: %v\n", err)
	}
	if out4 != t4.ExpectedResult {
		t.Errorf("Marshalled string doesn't match the expected output:\n'%v'\n", out4)
	}
}
