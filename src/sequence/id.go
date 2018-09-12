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

type IDType uint8

const MaxIDValue = 127
const WrapAroundValue = 128
const HalfMaxIDValue = MaxIDValue / 2

type ID struct {
	id IDType
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// NewID : Creates an ID
func NewID(id IDType) (ID, error) {
	return ID{id: id}, nil
}

func (i ID) Raw() IDType {
	return i.id
}

func (i ID) Next() ID {
	return ID{id: (i.id + 1) % WrapAroundValue}
}

func (i ID) Distance(next ID) int {
	nextValue := int(next.id)
	if next.id < i.id {
		nextValue += WrapAroundValue
	}
	diff := nextValue - int(i.id)
	return diff
}

func (i ID) IsSuccessor(next ID) bool {
	diff := i.Distance(next)
	return diff <= HalfMaxIDValue
}
