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
	"net/http"
	"testing"
)

func TestGetHeaderToken_MultipleValues(t *testing.T) {
	tests := []struct {
		name    string
		headers http.Header
		want    string
		wantErr error
	}{
		{
			name:    "Single value in Api-Key header",
			headers: http.Header{"Api-Key": []string{"value1"}},
			want:    "value1",
		},
		{
			name:    "No Api-Key header",
			headers: http.Header{},
			wantErr: ErrNoAuthHeaderIncluded,
		},
		{
			name:    "Empty Api-Key header",
			headers: http.Header{"Api-Key": []string{""}},
			wantErr: ErrNoAuthHeaderIncluded,
		},
		{
			name:    "Empty Api-Key header with other headers",
			headers: http.Header{"Api-Key": []string{""}, "Other-Header": []string{"value"}},
			wantErr: ErrNoAuthHeaderIncluded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHeaderToken(tt.headers)
			if err != tt.wantErr {
				t.Errorf("GetHeaderToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetHeaderToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
