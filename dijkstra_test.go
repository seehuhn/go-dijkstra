// seehuhn.de/go/dijkstra - find shortest paths in directed graphs
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

package dijkstra

import (
	"fmt"
	"testing"
)

type BinVertex uint32

func (v BinVertex) Edges() []BinEdge {
	var res []BinEdge
	res = append(res, BinEdge{from: v, to: v + 1})
	if v > 0 {
		res = append(res, BinEdge{from: v, to: v - 1})
	}
	if v > 0 && v%2 == 0 {
		res = append(res, BinEdge{from: v, to: v / 2})
	}
	res = append(res, BinEdge{from: v, to: 2 * v})
	return res
}

type BinEdge struct {
	from, to BinVertex
}

func (e BinEdge) From() BinVertex {
	return e.from
}

func (e BinEdge) To() BinVertex {
	return e.to
}

func (e BinEdge) Cost() int {
	return int(e.from) + 1
}

func TestBinary(t *testing.T) {
	ee, l := ShortestPath[int, BinEdge](BinVertex(100), BinVertex(1000))
	fmt.Println(ee, l)
	t.Error("fish")
}

type CircNode int

func (n CircNode) Edges() []CircEdge {
	var res CircNode
	if n >= 10 {
		res = 0
	} else {
		res = n + 1
	}
	return []CircEdge{{from: n, to: res}}
}

type CircEdge struct {
	from, to CircNode
}

func (e CircEdge) From() CircNode {
	return e.from
}

func (e CircEdge) To() CircNode {
	return e.to
}

func (e CircEdge) Cost() uint {
	return 1
}

func TestCircular(t *testing.T) {
	ee, l := ShortestPath[uint, CircEdge](CircNode(10), CircNode(11))
	fmt.Println(ee, l)
	t.Error("fish")
}
