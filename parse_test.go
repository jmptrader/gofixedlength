package gofixedlength

import (
	"testing"
	"time"
)

const (
	basicParseTestString = "1234567890ABCDEFGHIJ012.87"
	parseTestWithComma   = "1234567890ABCDEFGHIJ012,87"
	dateParseTestString  = "20150114EX"
)

type basicParseTest struct {
	NumberA int     `fixed:"0-5"`
	NumberB int     `fixed:"2-5"` // test overlap
	StringC string  `fixed:"10-15"`
	StringD string  `fixed:"29-35"` // should fail
	FloatA  float64 `fixed:"20-26"`
}

type dateParseTest struct {
	DateField   time.Time `fixed:"0-8,20060102"`
	StringAfter string    `fixed:"8-10"`
}

func TestBasicParsing(t *testing.T) {
	t.Log("Basic parsing test")
	var out basicParseTest
	Unmarshal(basicParseTestString, &out)
	if out.NumberA != 12345 {
		t.Errorf("NumberA parsed as %d", out.NumberA)
	}
	if out.NumberB != 345 {
		t.Errorf("NumberB parsed as %d", out.NumberB)
	}
	if out.StringC != "ABCDE" {
		t.Errorf("StringC parsed as '%s'", out.StringC)
	}
	if out.StringD != "" {
		t.Errorf("StringD should have failed to parse")
	}
	if out.FloatA != 12.87 {
		t.Errorf("FloatA parsed as '%v'", out.FloatA)
	}
}

func TestLayeredParsing(t *testing.T) {
	t.Log("Layered parsing test")
	var out dateParseTest
	Unmarshal(dateParseTestString, &out)
	if out.StringAfter != "EX" {
		t.Errorf("Failed to parse after embedded struct/ptr\n")
	}
	if expectedTime, err := time.Parse("2006-01-02", "2015-01-14"); err != nil || out.DateField != expectedTime {
		t.Errorf("Failed to parse date (%d)\n", out.DateField)
	}
}

func TestBasicParsingWithComma(t *testing.T) {
	previousValue := DECIMAL_COMMA
	DECIMAL_COMMA = true
	defer func() { DECIMAL_COMMA = previousValue }()
	var out basicParseTest
	Unmarshal(parseTestWithComma, &out)
	if out.FloatA != 12.87 {
		t.Errorf("FloatA parsed as '%v'", out.FloatA)
	}
}
