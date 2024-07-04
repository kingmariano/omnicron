// Copyright (c) 2024 Charles Ozochukwu

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package utils

import (
	"strconv"
)

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}

// SetIntValue sets the value of the int pointer if the given string is a valid integer.
func SetIntValue(value string, param **int) {
	if value != "" { // Check if the value string is not empty
		intValue, err := strconv.Atoi(value) // Convert the string to an integer
		if err == nil {                      // Check if the conversion was successful
			*param = &intValue // Set the parameter to the address of the integer value
		}
	}
}

// SetFloatValue sets the value of the float64 pointer if the given string is a valid float.
func SetFloatValue(value string, param **float64) {
	if value != "" { // Check if the value string is not empty
		floatValue, err := strconv.ParseFloat(value, 64) // Convert the string to a float64
		if err == nil {                                  // Check if the conversion was successful
			*param = &floatValue // Set the parameter to the address of the float value
		}
	}
}

// SetBoolValue sets the value of the bool pointer if the given string is a valid boolean.
func SetBoolValue(value string, param **bool) {
	if value != "" { // Check if the value string is not empty
		boolValue, err := strconv.ParseBool(value) // Convert the string to a boolean
		if err == nil {                            // Check if the conversion was successful
			*param = &boolValue // Set the parameter to the address of the boolean value
		}
	}
}

// SetStringValue sets the value of the string pointer if the given string is not empty.
func SetStringValue(value string, param **string) {
	if value != "" { // Check if the value string is not empty
		*param = &value // Set the parameter to the address of the string value
	}
}
