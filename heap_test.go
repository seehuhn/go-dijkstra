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

import "testing"

func (h *heap[vertex, edge, length]) verify(t *testing.T, i int) {
	t.Helper()

	n := len(h.candidates)
	if n == 0 {
		return
	}

	if h.index[h.candidates[i].to] != i {
		t.Errorf("index corrupted: %d != %d", h.index[h.candidates[i].to], i)
	}

	j1 := 2*i + 1
	j2 := 2*i + 2

	if j1 < n {
		if h.Less(j1, i) {
			t.Error("heap invariant invalidated")
			return
		}
		h.verify(t, j1)
	}
	if j2 < n {
		if h.Less(j2, i) {
			t.Error("heap invariant invalidated")
			return
		}
		h.verify(t, j2)
	}
}

func TestHeap(t *testing.T) {
	h := &heap[int, int, int]{
		index: make(map[int]int),
	}

	type C = candidate[int, int, int]

	// increasing order
	for i := 100; i < 400; i += 10 {
		h.Push(&C{to: i, via: 0, total: i, prev: nil})
	}
	h.verify(t, 0)

	// decreasing order
	for i := 900; i > 600; i -= 10 {
		h.Push(&C{to: i, via: 0, total: i, prev: nil})
	}
	h.verify(t, 0)

	// interleaved with previous entries
	for i := 15; i < 1000; i += 30 {
		h.Push(&C{to: i, via: 0, total: i, prev: nil})
	}
	h.verify(t, 0)

	last := 0
mainLoop:
	for len(h.candidates) > 0 {
		for pos := 3; pos < 7; pos++ {
			if pos < len(h.candidates) && h.candidates[pos].total%10 == 5 {
				h.candidates[pos].total = 505 - h.candidates[pos].total
				if h.candidates[pos].total < last {
					last = h.candidates[pos].total
				}
				h.Fix(pos)
				h.verify(t, 0)
				continue mainLoop
			}
		}

		c := h.Pop()
		if last > c.total {
			t.Errorf("wrong order: %d > %d", last, c.total)
		}
		last = c.total
		h.verify(t, 0)
	}
}
