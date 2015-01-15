package gofixedlength

import (
	"fmt"
	"testing"
	"time"
)

type Tester interface {
	Length() int
	ExpectedResult() string
	ExpectedErr() error
}

// Basic struct, valid
type testStruct1 struct {
	NumberA        int    `fixed:"0-5"`
	NumberB        int    `fixed:"5-10"`
	StringC        string `fixed:"10-15"`
	StringD        string `fixed:"15-35"`
	length         int
	expectedResult string
	expectedErr    error
}

func (t testStruct1) Length() int            { return t.length }
func (t testStruct1) ExpectedResult() string { return t.expectedResult }
func (t testStruct1) ExpectedErr() error     { return t.expectedErr }

// Basic struct, fields randomized, valid
type testStruct2 struct {
	StringD        string `fixed:"15-35"`
	NumberA        int    `fixed:"0-5"`
	StringC        string `fixed:"10-15"`
	NumberB        int    `fixed:"5-10"`
	length         int
	expectedResult string
	expectedErr    error
}

func (t testStruct2) Length() int            { return t.length }
func (t testStruct2) ExpectedResult() string { return t.expectedResult }
func (t testStruct2) ExpectedErr() error     { return t.expectedErr }

// Struct containing formatted date and float, valid
type testStruct3 struct {
	Uno            string    `fixed:"0-10"`
	Due            int       `fixed:"10-20"`
	Tre            time.Time `fixed:"20-30,2006-01-02"`
	Quattro        string    `fixed:"30-40"`
	Cinque         float32   `fixed:"40-50,2"`
	length         int
	expectedResult string
	expectedErr    error
}

func (t testStruct3) Length() int            { return t.length }
func (t testStruct3) ExpectedResult() string { return t.expectedResult }
func (t testStruct3) ExpectedErr() error     { return t.expectedErr }

// Basic struct with holes, valid
type testStruct4 struct {
	NumberA        int    `fixed:"0-5"`
	NumberB        int    `fixed:"5-10"`
	StringC        string `fixed:"10-15"`
	StringD        string `fixed:"20-40"` // There is a gap of 5 columns between StringC and StringD (15-20)
	length         int
	expectedResult string
	expectedErr    error
}

func (t testStruct4) Length() int            { return t.length }
func (t testStruct4) ExpectedResult() string { return t.expectedResult }
func (t testStruct4) ExpectedErr() error     { return t.expectedErr }

// Basic struct with overlap. Should only fail if the contested columns are not consistently defined.
type testStruct5 struct {
	NumberA        int    `fixed:"0-5"`
	NumberB        int    `fixed:"5-10"`
	StringC        string `fixed:"10-15"`
	StringD        string `fixed:"12-22"` // There is an overlap in pos. 12-15 between StringC and StringD
	length         int
	expectedResult string
	expectedErr    error
}

func (t testStruct5) Length() int            { return t.length }
func (t testStruct5) ExpectedResult() string { return t.expectedResult }
func (t testStruct5) ExpectedErr() error     { return t.expectedErr }

type embeddedStruct struct {
	Another    string `fixed:"50-60"`
	AndFinally string `fixed:"60-62"`
}

// Embedded struct, valid
type testStruct6 struct {
	Embedded1     testStruct1 // Basic struct, valid
	AnotherField1 int         `fixed:"35-40"`
	AnotherField2 int         `fixed:"40-42"`
	AnotherField3 string      `fixed:"42-47"`
	AnotherField4 string      `fixed:"47-50"`
	embeddedStruct
	length         int
	expectedResult string
	expectedErr    error
}

func (t testStruct6) Length() int            { return t.length }
func (t testStruct6) ExpectedResult() string { return t.expectedResult }
func (t testStruct6) ExpectedErr() error     { return t.expectedErr }

// Struct with invalid index. Should fail.
type testFailingStruct1 struct {
	Uno            string `fixed:"0-10"`
	Due            int    `fixed:"10-9"`
	length         int
	expectedResult string
	expectedErr    error
}

func (t testFailingStruct1) Length() int            { return t.length }
func (t testFailingStruct1) ExpectedResult() string { return t.expectedResult }
func (t testFailingStruct1) ExpectedErr() error     { return t.expectedErr }

// Struct with invalid index. Should fail.
type testFailingStruct2 struct {
	Uno            string `fixed:"0-10"`
	Due            int    `fixed:"10-"`
	length         int
	expectedResult string
	expectedErr    error
}

func (t testFailingStruct2) Length() int            { return t.length }
func (t testFailingStruct2) ExpectedResult() string { return t.expectedResult }
func (t testFailingStruct2) ExpectedErr() error     { return t.expectedErr }

// Struct with invalid index. Should fail.
type testFailingStruct3 struct {
	Uno            string `fixed:"0-10"`
	Due            int    `fixed:"-20"`
	length         int
	expectedResult string
	expectedErr    error
}

func (t testFailingStruct3) Length() int            { return t.length }
func (t testFailingStruct3) ExpectedResult() string { return t.expectedResult }
func (t testFailingStruct3) ExpectedErr() error     { return t.expectedErr }

// Embedded struct with fixed length (shorter than embedded struct's length). Should fail.
type testFailingStruct4 struct {
	EmbeddedStruct testStruct1 `fixed:"0-30"` // Basic struct, valid
	AnotherField1  int         `fixed:"30-40"`
	AnotherField2  int         `fixed:"40-42"`
	AnotherField3  string      `fixed:"42-47"`
	AnotherField4  string      `fixed:"47-50"`
	length         int
	expectedResult string
	expectedErr    error
}

func (t testFailingStruct4) Length() int            { return t.length }
func (t testFailingStruct4) ExpectedResult() string { return t.expectedResult }
func (t testFailingStruct4) ExpectedErr() error     { return t.expectedErr }

var (
	timeObject = time.Now()
	TestSuite  = [...]Tester{
		testStruct1{
			NumberA:        123,                 // `fixed:"0-5"`
			NumberB:        12345,               // `fixed:"5-10"`
			StringC:        "ohmy",              // `fixed:"10-15"`
			StringD:        "What's happening?", // `fixed:"15-35"`
			length:         35,
			expectedResult: "0012312345ohmy What's happening?   ",
		},
		testStruct2{
			StringD:        "Some more text 12345", // `fixed:"15-35"`
			NumberA:        12345,                  // `fixed:"0-5"`
			StringC:        "AND",                  // `fixed:"10-15"`
			NumberB:        67890,                  // `fixed:"5-10"`
			length:         35,
			expectedResult: "1234567890AND  Some more text 12345",
		},
		testStruct3{
			Uno:            "just atext", // string    `fixed:"0-10"`
			Due:            1234567890,   // int       `fixed:"10-20"`
			Tre:            timeObject,   // time.Time `fixed:"20-30,2006-01-02"`
			Quattro:        "anotherStr", // string    `fixed:"30-40"`
			Cinque:         321.458,      // float32   `fixed:"40-50,2"`
			length:         50,           // int
			expectedResult: fmt.Sprintf("just atext1234567890%vanotherStr0000321.46", timeObject.Format("2006-01-02")),
		},
		testStruct4{
			NumberA:        12345,                  // int    `fixed:"0-5"`
			NumberB:        67890,                  // int    `fixed:"5-10"`
			StringC:        "short",                // string `fixed:"10-15"`
			StringD:        "and then some       ", // string `fixed:"20-40"` // There is a gap of 5 columns between StringC and StringD (15-20)
			length:         40,                     // int
			expectedResult: "1234567890short     and then some       ",
		},
		testStruct5{
			NumberA:        123,         // int    `fixed:"0-5"`
			NumberB:        0,           // int    `fixed:"5-10"`
			StringC:        "overl",     // string `fixed:"10-15"`
			StringD:        "erlapping", // string `fixed:"12-22"` // There is an overlap in pos. 12-15 between StringC and StringD
			length:         22,
			expectedResult: "0012300000overlapping ",
		},
		testStruct5{
			NumberA:     123,          // int    `fixed:"0-5"`
			NumberB:     0,            // int    `fixed:"5-10"`
			StringC:     "overl",      // string `fixed:"10-15"`
			StringD:     "ER THE TOP", // string `fixed:"12-22"` // There is an overlap in pos. 12-15 between StringC and StringD
			length:      22,
			expectedErr: ErrIncoherentOverlap,
		},
	}
)

func TestLinelength(t *testing.T) {
	t.Log("length calculator test")
	for i, target := range TestSuite {
		if length := LineLength(target); length != target.Length() {
			t.Errorf("Failed to find the line length for test no.%v (found %v, expected %v)", i+1, length, target.Length())
		}
	}
}

func TestMarshal(t *testing.T) {
	t.Log("Marshaller test")
	for i, target := range TestSuite {
		out, err := Marshal(target)
		if err != target.ExpectedErr() {
			if err != nil {
				t.Errorf("Error while marshaling the test struct no.%v: %v\n", i+1, err)
			} else {
				t.Errorf("Struct no.%v should have failed and passed instead, resulting: %v\n", i+1, out)
			}
		}
		if err == nil && out != target.ExpectedResult() {
			t.Errorf("Marshalled string no.%v doesn't match the expected output:\n'%v'\n", i+1, out)
		}
	}
}
