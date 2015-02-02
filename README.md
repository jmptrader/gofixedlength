# GoFixedLength 

[![Build Status](https://secure.travis-ci.org/qrawl/gofixedlength.png)](http://travis-ci.org/qrawl/gofixedlength)
[![GoDoc](https://godoc.org/github.com/qrawl/gofixedlength?status.png)](https://godoc.org/github.com/qrawl/gofixedlength)

Go library to deal with extracting fixed field form values using struct tags.  

##Quickstart

**Unmarshal** unmarshals string data into an annotated interface. This should resemble:

	type SomeType struct {
		DateField   time.Time `fixed:"0-8,20060102"` // Standard Go time formatting
		StringAfter string    `fixed:"8-15"`
		SomeFloat   float64   `fixed:"15-25,4"`      // Four decimals
	}

	var out SomeType
	err := Unmarshal("20150202well   00012.1864", &out)


**Marshal** marshals struct data into a fixed-lenght formatted string.

	type SomeType struct {
		ValA string        `fixed:"0-10"`
		ValB int           `fixed:"10-20"`
		ValC time.Time     `fixed:"20-30,2006-01-02"` // Standard Go time formatting
		ValD float64       `fixed:"30-40,3"`          // Three decimals
	}

	myStruct := SomeType{
		"this",
		12345,
		time.Now(),
		123.1234,
	}

	out, err := gofixedlength.Marshal(myStruct)
	// out == "this      00000123452015-01-14000123.123"

Offsets are zero based.  
Field filling is based on data type: for text types it will be spaces, while numbers will be right-aligned and filled with zeroes.  
Floating-point values are printed with the specified number of decimals (two by default).  
`time.Time` fields are printed with the specified layout.


String offsets are zero based.

##European-styled numbers
To parse documents that use comma "," as decimal separator, just set to `true` the global variable:

	func init() {
		gofixedlength.DECIMAL_COMMA = true
	}

GoFixedLength is based on @jbuchbinder's [Gofixedfield](https://github.com/jbuchbinder/gofixedfield).
