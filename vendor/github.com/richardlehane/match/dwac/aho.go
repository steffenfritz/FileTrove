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

// This file implements Aho Corasick searching for the bytematcher package
package dwac

// out function
type out struct {
	max      int64 // maximum offset at which can occur
	seqIndex int   // index within all the Seqs in the Wac
	subIndex int   // index of the Choice within the Seq
	length   int   // length of byte slice
}

func contains(op []out, o out) bool {
	if op == nil {
		return false
	}
	for _, o1 := range op {
		if o == o1 {
			return true
		}
	}
	return false
}

func addOutput(op []out, o out, outMax int64, outMaxL int) ([]out, int64, int) {
	if op == nil {
		return []out{o}, o.max, o.length
	}
	if outMax > -1 && (o.max == -1 || o.max > outMax) {
		outMax = o.max
	}
	if o.length > outMaxL {
		outMaxL = o.length
	}
	return append(op, o), outMax, outMaxL
}

// regular node
type node struct {
	val     byte
	keys    []byte
	transit [256]*node // the goto function
	fail    *node      // the fail function
	output  []out      // the output function
	outMax  int64
	outMaxL int
}

func (start *node) addGotos(seqs []Seq) (int64, bool) {
	var maxOff int64
	var hasWild bool
	// iterate through byte sequences adding goto links to the link matrix
	for id, seq := range seqs {
		for i, choice := range seq.Choices {
			for _, byts := range choice {
				curr := start
				for _, byt := range byts {
					if curr.transit[byt] == nil {
						curr.transit[byt] = &node{
							val:  byt,
							keys: make([]byte, 0, 1),
						}
						curr.keys = append(curr.keys, byt)
					}
					curr = curr.transit[byt]
				}
				curr.output, curr.outMax, curr.outMaxL = addOutput(
					curr.output,
					out{seq.MaxOffsets[i], id, i, len(byts)},
					curr.outMax,
					curr.outMaxL)
				if seq.MaxOffsets[i] > maxOff {
					maxOff = seq.MaxOffsets[i]
				}
				if seq.MaxOffsets[i] == -1 {
					hasWild = true
				}
			}
		}
	}
	return maxOff, hasWild
}

func (start *node) addGotosIndexes(idxs []SeqIndex, seqs []Seq) {
	for _, idx := range idxs {
		for i, choice := range seqs[idx[0]].Choices[idx[1]:] {
			for _, byts := range choice {
				curr := start
				for _, byt := range byts {
					if curr.transit[byt] == nil {
						curr.transit[byt] = &node{
							val:  byt,
							keys: make([]byte, 0, 1),
						}
						curr.keys = append(curr.keys, byt)
					}
					curr = curr.transit[byt]
				}
				curr.output, curr.outMax, curr.outMaxL = addOutput(
					curr.output,
					out{-1, idx[0], i + idx[1], len(byts)},
					curr.outMax,
					curr.outMaxL)
			}
		}
	}
}

func (start *node) addFails() {
	// root and its children fail to root
	start.fail = start
	for _, byt := range start.keys {
		start.transit[byt].fail = start
	}
	// traverse tree in breadth first search adding fails
	queue := make([]*node, 0, 50)
	queue = append(queue, start)
	for len(queue) > 0 {
		pop := queue[0]
		for _, byt := range pop.keys {
			node := pop.transit[byt]
			queue = append(queue, node)
			// starting from the node's parent, follow the fails back towards root,
			// and stop at the first fail that has a goto to the node's value
			fail := pop.fail
			ok := fail.transit[node.val]
			for fail != start && ok == nil {
				fail = fail.fail
				ok = fail.transit[node.val]
			}
			fnode := fail.transit[node.val]
			if fnode != nil && fnode != node {
				node.fail = fnode
			} else {
				node.fail = start
			}
			// another traverse back to root following the fails. This time add any unique out functions to the node
			fail = node.fail
			for fail != start {
				for _, o := range fail.output {
					if !contains(node.output, o) {
						node.output, node.outMax, node.outMaxL = addOutput(node.output, o, node.outMax, node.outMaxL)
					}
				}
				fail = fail.fail
			}
		}
		queue = queue[1:]
	}
}

// preconditions ensure that subsequent (>0) Choices in a Seq are only sent when previous Choices have already matched
// previous matches are stored as offsets to prevent overlapping matches resulting in false positives
type precons [][]int64

func newPrecons(t []int) *precons {
	p := make(precons, len(t))
	for i, v := range t {
		p[i] = make([]int64, v)
	}
	return &p
}

func clear(p precons) *precons {
	for i := range p {
		for j := range p[i] {
			p[i][j] = 0
		}
	}
	return &p
}

func makeT(s []Seq) []int {
	t := make([]int, len(s))
	for i := range s {
		t[i] = len(s[i].Choices)
	}
	return t
}

func preconsFn(s []Seq) func() interface{} {
	t := makeT(s)
	return func() interface{} {
		return newPrecons(t)
	}
}
