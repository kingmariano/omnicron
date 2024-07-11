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

package auth

import (
	"errors"
	"net/http"
	"testing"
)

func TestGetHeaderToken(t *testing.T) {
	tests := []struct {
		name        string
		headers     http.Header
		expected    string
		expectedErr error
	}{
		{
			name:        "No Authorization header",
			headers:     http.Header{},
			expected:    "",
			expectedErr: ErrNoAuthHeaderIncluded,
		},
		{
			name: "Malformed Authorization header without Bearer",
			headers: http.Header{
				"Authorization": {"Token abcdef12345"},
			},
			expected:    "",
			expectedErr: errors.New("malformed authorization header"),
		},
		{
			name: "Malformed Authorization header with only Bearer",
			headers: http.Header{
				"Authorization": {"Bearer"},
			},
			expected:    "",
			expectedErr: errors.New("malformed authorization header"),
		},
		{
			name: "Valid Authorization header",
			headers: http.Header{
				"Authorization": {"Bearer abcdef12345"},
			},
			expected:    "abcdef12345",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetHeaderToken(tt.headers)
			if token != tt.expected {
				t.Errorf("expected token %v, got %v", tt.expected, token)
			}
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err == nil && tt.expectedErr != nil {
				t.Errorf("expected error %v, got nil", tt.expectedErr)
			}
		})
	}
}
