package gofixedlength

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrBeginOutOfRange     = errors.New("Begin index is out of range")
	ErrEndOutOfRange       = errors.New("End index is out of range")
	ErrTextTooLongForRange = errors.New("Text is longer than the range")
	ErrIncoherentOverlap   = errors.New("The function tried to rewrite a different value on the same column")
)

type Line []rune

// Marshal marshals struct data into a fixed-lenght formatted string.
//
// 	type SomeType struct {
// 		ValA string        `fixed:"0-10"`
//		ValB int           `fixed:"10-20"`
//		ValC time.Time     `fixed:"20-30,2006-01-02"`
//		ValD float         `fixed:"30-40,3"`
// 	}
//
//	myStruct := SomeType{
//		"this",
//		12345,
//		time.Now(),
//		123.1234,
//	}
//
//	out, err := Marshal(myStruct)
//	// out == "this      00000123452015-01-14000123.123"
//
// String offsets are zero based.
// Field filling is based on data type: for text types it will be spaces,
// while numbers will be right-aligned and filled with zeroes.
// Floating point-values are printed with the specified number of decimals (two by default).
// time.Time fields are printed in the specified layout.
func Marshal(v interface{}) (string, error) {
	var line Line // Build a rune array the length the output line is supposed to be
	line = make([]rune, LineLength(v))
	//debugStruct(v)
	var val reflect.Value
	val = reflect.ValueOf(v)
	/*
		var val reflect.Value
		if reflect.TypeOf(v).Name() != "" {
			val = reflect.ValueOf(v)
		} else {
			val = reflect.ValueOf(v).Elem()
		}
	*/
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		cArguments := strings.SplitN(tag.Get("fixed"), ",", 2)
		var cFormat string
		if len(cArguments) > 1 {
			cFormat = cArguments[1]
		}
		cRange := cArguments[0]
		cBookend := strings.Split(cRange, "-")
		if len(cBookend) != 2 {
			// If we don't have two values, skip
			continue
		}

		b, _ := strconv.Atoi(cBookend[0])
		e, _ := strconv.Atoi(cBookend[1])

		/*
			// Sanity check range before dying miserably
			if b < 0 || e > len(v) {
				continue
			}
		*/
		fieldLength := e - b

		//log.Println("CHE C'E QUA DENTRO?", reflect.ValueOf(v).Field(i))

		switch typeField.Type.Kind() {
		case reflect.Bool:
			/*
				format := fmt.Sprintf("%%%dv", fieldLength, fieldLength)
					if *val.Field(i) {
						out.WriteString(fmt.Sprintf(format, '1'))
					} else {
						out.WriteString(fmt.Sprintf(format, '0'))
					}
			*/
			break
		case reflect.Float32, reflect.Float64:
			// cFormat is the number of decimals
			decimals, err := strconv.Atoi(cFormat)
			if err != nil {
				log.Println("Found non-valid format for float:", cFormat)
			}
			integerPartLength := fieldLength - 1 - decimals
			integerPart := int(reflect.ValueOf(v).Field(i).Float())
			if integerPart >= pow(integerPartLength, 10) {
				log.Printf("This float number (%v) seems to be too big for output length (%v).\n", integerPart, integerPartLength)
			}
			format := fmt.Sprintf("%%0%d.%df", fieldLength, decimals)
			outstring := fmt.Sprintf(format, reflect.ValueOf(v).Field(i).Float()) // Doesn't check if the source float is too long
			if DECIMAL_COMMA {
				outstring = strings.Replace(outstring, ".", ",", 1)
			}
			err = line.WriteString(outstring, b, e)
			if err != nil {
				return line.String(), err
			}
			break
		case reflect.String:
			format := fmt.Sprintf("%%-%ds", fieldLength)
			outstring := fmt.Sprintf(format, reflect.ValueOf(v).Field(i).String())
			outstring = outstring[0:fieldLength]
			err := line.WriteString(outstring, b, e)
			if err != nil {
				return line.String(), err
			}
			break
		case reflect.Int8, reflect.Int32, reflect.Int, reflect.Int64, reflect.Uint:
			format := fmt.Sprintf("%%0%dv", fieldLength)
			outstring := fmt.Sprintf(format, reflect.ValueOf(v).Field(i).Int())
			err := line.WriteString(outstring, b, e)
			if err != nil {
				return line.String(), err
			}
			break
		case reflect.Ptr, reflect.Struct:
			if typeField.Type == reflect.TypeOf(time.Time{}) {
				// cFormat is the time.Format() format
				if len(cFormat) != fieldLength {
					log.Println("cFormat for this time.Time object doesn't match the field length") // Maybe this kind of parsing error check should be done elsewhere
				}
				outstring := reflect.ValueOf(v).Field(i).Interface().(time.Time).Format(cFormat)
				line.WriteString(outstring, b, e)
			} else {

				// Handle embedded objects by recursively parsing
				// the object with the range we passed.
				if val.Field(i).IsNil() {
					// Initialize pointer to avoid panic
					val.Field(i).Set(reflect.New(val.Field(i).Type().Elem()))
				}
				marshalledStruct, err := Marshal(val.Field(i).Interface())
				if err != nil {
					return line.String(), err
				}
				err = line.WriteString(marshalledStruct, 0, line.Length())
				if err != nil {
					return line.String(), err
				}
			}
			break
		default:
			break
		}
	}
	return line.String(), nil
}

func pow(a, b int) int {
	var c int
	c = 1
	for i := 0; i < b; i++ {
		c = c * a
	}
	return c
}

// Returns the total length of the line we're going to marshal the data to, iterating
// all the struct's fields and returning the higher number in the `field` tag.
func LineLength(v interface{}) int {
	var higherNumber int

	var val reflect.Value
	val = reflect.ValueOf(v)

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		cArguments := strings.SplitN(tag.Get("fixed"), ",", 2)
		cRange := cArguments[0]
		cBookend := strings.Split(cRange, "-")

		if len(cBookend) == 2 {
			if e, _ := strconv.Atoi(cBookend[1]); e > higherNumber {
				higherNumber = e
			}
		}

		// Iterate thgough the embedded struct if it's not a time.Time object
		if typeField.Type.Kind() == reflect.Struct || typeField.Type.Kind() == reflect.Ptr { // Should I handle the reflect.Ptr case?
			if typeField.Type != reflect.TypeOf(time.Time{}) {
				higherSubNumber := LineLength(val.Field(i).Interface())
				if higherSubNumber > higherNumber {
					higherNumber = higherSubNumber
				}
			}
		}
	}
	return higherNumber
}

func (l Line) WriteString(text string, begin, end int) error {
	textRunesCount := utf8.RuneCountInString(text)
	if begin < 0 || begin > l.Length()-1 {
		return ErrBeginOutOfRange
	}
	if end < 1 || end > l.Length() {
		return ErrEndOutOfRange
	}
	if textRunesCount > end-begin {
		return ErrTextTooLongForRange
	}
	for j, i, w := begin, 0, 0; i < len(text); i += w {
		runeValue, width := utf8.DecodeRuneInString(text[i:])
		if l[j] != '\x00' && l[j] != runeValue {
			return ErrIncoherentOverlap
		}
		l[j] = runeValue

		w = width // Next iteration will start with the cursor on the next rune of the input
		j++       // Next iteration will affect the next rune of the output
	}
	return nil
}

func (l Line) Length() int {
	return len(l)
}

func (l Line) String() string {
	return string(l)
}
