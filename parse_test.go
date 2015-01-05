package gofixedlength

import (
	"testing"
)

const (
	basicParseTestString   = "1234567890ABCDEFGHIJ"
	layeredParseTestString = "20091010EX"
)

type basicParseTest struct {
	NumberA int    `fixed:"0-5"`
	NumberB int    `fixed:"2-5"` // test overlap
	StringC string `fixed:"10-15"`
	StringD string `fixed:"29-35"` // should fail
}

type layeredParseTest struct {
	DateField   *dateStruct `fixed:"0-8"`
	StringAfter string      `fixed:"8-10"`
}

type dateStruct struct {
	Y int `fixed:"0-4"`
	M int `fixed:"4-6"`
	D int `fixed:"6-8"`
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
}

func TestLayeredParsing(t *testing.T) {
	t.Log("Layered parsing test")
	var out layeredParseTest
	Unmarshal(layeredParseTestString, &out)
	if out.StringAfter != "EX" {
		t.Errorf("Failed to parse after embedded struct/ptr\n")
	}
	if out.DateField.Y != 2009 {
		t.Errorf("Failed to parse embedded Y (Y=%d)\n", out.DateField.Y)
	}
	if out.DateField.M != 10 {
		t.Errorf("Failed to parse embedded M (M=%d)\n", out.DateField.M)
	}
	if out.DateField.D != 10 {
		t.Errorf("Failed to parse embedded D (D=%d)\n", out.DateField.D)
	}
}
