# GoFixedField 

[![Build Status](https://secure.travis-ci.org/qrawl/gofixedfield.png)](http://travis-ci.org/qrawl/gofixedfield)
[![GoDoc](https://godoc.org/github.com/qrawl/gofixedfield?status.png)](https://godoc.org/github.com/qrawl/gofixedfield)

Go library to deal with extracting fixed field form values using struct tags.  
This is a fork of [@jbuchbinder](https://github.com/jbuchbinder)'s GOFIXEDFIELD. The only noticeable change is that I prefer zero-based character count.

##Quickstart

Unmarshal unmarshals string data into an annotated interface. This should
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
	err := Unmarshal("some string here", &out)

String offsets are zero based.

##European-styled numbers
To parse documents that use a comma "," instead of the decimal point, just set to `true` the corresponding global variable:

	func init() {
		gofixedfield.DECIMAL_COMMA = true
	}
