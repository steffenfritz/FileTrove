// Copyright 2019 Richard Lehane. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// DWAC is a multiple string matching algorithm with choices and max offsets.
// It pauses matching when all strings with fixed offsets have been checked.
// To resume matching (a limited set) of wildcard sequences, send on the resume channel.
package dwac

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

// Result contains the index and offset of matches.
type Result struct {
	Index  [2]int // a double index: index of the Seq and index of the Choice
	Offset int64
	Length int
}

// Choice represents the different byte slices that can occur at each position of the Seq
type Choice [][]byte

// Seq is an ordered set of slices of Choices, with maximum offsets for each choice
type Seq struct {
	MaxOffsets []int64 // maximum offsets for each choice. Can be -1 for wildcard.
	Choices    []Choice
}

// SeqIndex is an index into the slice of Seqs and their containing choices used to create the Dwac
type SeqIndex [2]int

func (s Seq) String() string {
	str := "{Offsets:"
	for n, v := range s.MaxOffsets {
		if n > 0 {
			str += ","
		}
		str += fmt.Sprintf(" %d", v)
	}
	str += "; Choices:"
	for n, v := range s.Choices {
		if n > 0 {
			str += ","
		}
		str += " ["
		strs := make([]string, len(v))
		for i := range v {
			strs[i] = string(v[i])
		}
		str += strings.Join(strs, " | ")
		str += "]"
	}
	return str + "}"
}

type Dwac struct {
	maxOff  int64
	hasWild bool
	root    *node
	p       *sync.Pool
	seqs    []Seq
}

func New(seqs []Seq) *Dwac {
	d := &Dwac{}
	d.root = &node{}
	d.maxOff, d.hasWild = d.root.addGotos(seqs)
	d.root.addFails()
	if d.hasWild {
		d.seqs = seqs
	}
	d.p = &sync.Pool{New: preconsFn(seqs)}
	return d
}

// Dwac returns a channel of results, which are double indexes (of the Seq and of the Choice),
// and a resume channel, which is a slice of wild Seq indexes
func (d *Dwac) Index(rdr io.ByteReader) (<-chan Result, chan<- []SeqIndex) {
	output, resume := make(chan Result), make(chan []SeqIndex)
	go d.match(rdr, output, resume)
	return output, resume
}

func (dwac *Dwac) match(input io.ByteReader, results chan Result, resume chan []SeqIndex) {
	var offset int64
	p := *(dwac.p.Get().(*precons))
	curr := dwac.root
	var c byte
	var err error
	for c, err = input.ReadByte(); err == nil; c, err = input.ReadByte() {
		offset++
		if trans := curr.transit[c]; trans != nil {
			curr = trans
		} else {
			for curr != dwac.root {
				curr = curr.fail
				if trans := curr.transit[c]; trans != nil {
					curr = trans
					break
				}
			}
		}
		if curr.output != nil && (curr.outMax == -1 || curr.outMax >= offset-int64(curr.outMaxL)) {
			for _, o := range curr.output {
				if o.max == -1 || o.max >= offset-int64(o.length) {
					if o.subIndex == 0 || (p[o.seqIndex][o.subIndex-1] != 0 && offset-int64(o.length) >= p[o.seqIndex][o.subIndex-1]) {
						if p[o.seqIndex][o.subIndex] == 0 {
							p[o.seqIndex][o.subIndex] = offset
						}
						results <- Result{Index: [2]int{o.seqIndex, o.subIndex}, Offset: offset - int64(o.length), Length: o.length}
					}
				}
			}
		}
		if offset > int64(dwac.maxOff) && curr == dwac.root {
			break
		}
	}
	// if EOF not reached or other file read error, try the resume channel
	if err == nil && dwac.hasWild {
		results <- Result{Index: [2]int{-1, -1}, Offset: offset}
		seqIndexes := <-resume
		if len(seqIndexes) > 0 {
			root := &node{}
			root.addGotosIndexes(seqIndexes, dwac.seqs)
			root.addFails()
			curr = root
			for c, err = input.ReadByte(); err == nil; c, err = input.ReadByte() {
				offset++
				if trans := curr.transit[c]; trans != nil {
					curr = trans
				} else {
					for curr != root {
						curr = curr.fail
						if trans := curr.transit[c]; trans != nil {
							curr = trans
							break
						}
					}
				}
				if curr.output != nil {
					for _, o := range curr.output {
						if o.subIndex == 0 || (p[o.seqIndex][o.subIndex-1] != 0 && offset-int64(o.length) >= p[o.seqIndex][o.subIndex-1]) {
							if p[o.seqIndex][o.subIndex] == 0 {
								p[o.seqIndex][o.subIndex] = offset
							}
							results <- Result{Index: [2]int{o.seqIndex, o.subIndex}, Offset: offset - int64(o.length), Length: o.length}
						}
					}
				}
			}
		}
	}
	// return precons
	dwac.p.Put(clear(p))
	close(results)
}
