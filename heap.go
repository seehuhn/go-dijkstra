// seehuhn.de/go/dijkstra - shortest paths in directed graphs
// Copyright (C) 2022  Jochen Voss <voss@seehuhn.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Part of this code is copied from the Go standard library
// https://golang.org/src/container/heap/ and then modified.  Use of this
// source code is governed by a BSD-style license, which is reproduced here:
//
//     Copyright (c) 2009 The Go Authors. All rights reserved.
//
//     Redistribution and use in source and binary forms, with or without
//     modification, are permitted provided that the following conditions are
//     met:
//
//        * Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer.
//        * Redistributions in binary form must reproduce the above
//     copyright notice, this list of conditions and the following disclaimer
//     in the documentation and/or other materials provided with the
//     distribution.
//        * Neither the name of Google Inc. nor the names of its
//     contributors may be used to endorse or promote products derived from
//     this software without specific prior written permission.
//
//     THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
//     "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
//     LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
//     A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
//     OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
//     SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
//     LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
//     DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
//     THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
//     (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
//     OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package dijkstra

type candidate[vertex V, edge E, length L] struct {
	to    vertex
	via   edge
	total length
	prev  *candidate[vertex, edge, length]
}

type heap[vertex V, edge E, length L] struct {
	candidates []*candidate[vertex, edge, length]
	index      map[vertex]int
}

func (h *heap[vertex, edge, length]) Less(i, j int) bool {
	cand := h.candidates
	return cand[i].total < cand[j].total
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func (h *heap[vertex, edge, length]) Push(cand *candidate[vertex, edge, length]) {
	n := len(h.candidates)
	h.index[cand.to] = n
	h.candidates = append(h.candidates, cand)
	h.up(n)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to Remove(h, 0).
func (h *heap[vertex, edge, length]) Pop() *candidate[vertex, edge, length] {
	cc := h.candidates

	n := len(cc) - 1

	x := cc[0]
	delete(h.index, x.to)

	if n > 0 {
		cc[0] = cc[n]
		h.index[cc[0].to] = 0
		h.down(0, n)
	}
	h.candidates = cc[:n]

	return x
}

// Re-establish the heap ordering after the element idx
// has changed its value.  The complexity is O(log n) where n
// is the number of candidates.
func (h *heap[vertex, edge, length]) Fix(idx int) {
	if !h.down(idx, len(h.candidates)) {
		h.up(idx)
	}
}

func (h *heap[vertex, edge, length]) up(j int) {
	cand := h.candidates

	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.Less(j, i) {
			break
		}

		// swap
		cand[i], cand[j] = cand[j], cand[i]
		h.index[cand[j].to] = j

		j = i
	}
	h.index[cand[j].to] = j
}

func (h *heap[vertex, edge, length]) down(i0, n int) bool {
	cand := h.candidates

	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
			j = j2 // right child
		}
		if !h.Less(j, i) {
			break
		}

		// swap
		cand[i], cand[j] = cand[j], cand[i]
		h.index[cand[i].to] = i

		i = j
	}

	if i > i0 {
		h.index[cand[i].to] = i
		return true
	}
	return false
}
