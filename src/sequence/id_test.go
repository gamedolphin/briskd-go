/*

MIT License

Copyright (c) 2017 Peter Bjorklund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package sequence

import (
	"testing"
)

func TestDistance(t *testing.T) {
	first, _ := NewID(10)
	second, _ := NewID(253)

	distance := second.Distance(first)
	if distance != 13 {
		t.Errorf("Not correct distance:%d, expected 13", distance)
	}
}

func TestPreviousID(t *testing.T) {
	first, _ := NewID(10)
	second, _ := NewID(9)

	distance := first.Distance(second)
	if distance != 255 {
		t.Errorf("Not correct distance:%d, expected 255", distance)
	}
}

func TestPreviousIDWrap(t *testing.T) {
	first, _ := NewID(0)
	second, _ := NewID(255)

	distance := first.Distance(second)
	if distance != 255 {
		t.Errorf("Not correct distance:%d, expected 255", distance)
	}
}

func TestNextWrap(t *testing.T) {
	first, _ := NewID(255)
	second, _ := NewID(125)

	distance := first.Distance(second)
	if distance != 126 {
		t.Errorf("Not correct distance:%d, expected 126", distance)
	}

	isSuccessor := first.IsSuccessor(second)
	if !isSuccessor {
		t.Errorf("Should be successor")
	}
}
