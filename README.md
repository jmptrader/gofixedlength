# GoFixedLength 

[![Build Status](https://secure.travis-ci.org/qrawl/gofixedlength.png)](http://travis-ci.org/qrawl/gofixedlength)
[![GoDoc](https://godoc.org/github.com/qrawl/gofixedlength?status.png)](https://godoc.org/github.com/qrawl/gofixedlength)

Go library to deal with extracting fixed field form values using struct tags.  

##Quickstart

**Unmarshal** unmarshals string data into an annotated interface. This should
resemble:

	type SomeType struct {
		ValA string        `fixed:"0-5"`
		ValB int           `fixed:"9-15"`
		ValC *EmbeddedType `fixed:"15-22"`
 	}
	type EmbeddedType struct {
		ValX string `fixed:"0-3"`
		ValY string `fixed:"3-6"`
	}

	var out SomeType
	err := gofixedlength.Unmarshal("some string here", &out)

**Marshal** marshals struct data into a fixed-lenght formatted string.

 	type SomeType struct {
 		ValA string        `fixed:"0-10"`
		ValB int           `fixed:"10-20"`
		ValC time.Time     `fixed:"20-30,2006-01-02"`
		ValD float         `fixed:"30-40,3"`
 	}

	myStruct := SomeType{
		"this",
		12345,
		time.Now(),
		123.1234,
	}

	out, err := gofixedlength.Marshal(myStruct)
	// out == "this      00000123452015-01-14000123.123"

String offsets are zero based.  
Field filling is based on data type: for text types it will be spaces,
while numbers will be right-aligned and filled with zeroes.  
Floating-point values are printed with the specified number of decimals (two by default).  
`time.Time` fields are printed with the specified layout.


String offsets are zero based.

##European-styled numbers
To parse documents that use comma "," as decimal separator, just set to `true` the global variable:

	func init() {
		gofixedlength.DECIMAL_COMMA = true
	}

GoFixedLength is based on @jbuchbinder's [Gofixedfield](https://github.com/jbuchbinder/gofixedfield).
