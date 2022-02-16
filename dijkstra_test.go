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

func (e BinEdge) Length() float64 {
	return 1 + 1/float64(e.from)
}

func TestBinary(t *testing.T) {
	path, err := ShortestPath[BinEdge, BinVertex, float64](BinVertex(100), BinVertex(1000))
	if err != nil {
		t.Fatal(err)
	}
	if len(path) < 2 || path[0].from != 100 || path[len(path)-1].to != 1000 {
		t.Error("wrong path")
	}
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

func (e CircEdge) Length() uint {
	return 1
}

func TestCircular(t *testing.T) {
	from := CircNode(10)

	for _, to := range []CircNode{9, 10, 11} {
		path, err := ShortestPath[CircEdge, CircNode, uint](from, to)
		switch {
		case to > 10:
			if err != ErrNoPath || path != nil {
				t.Error("wrong path found")
			}
		case to == from:
			if err != nil || len(path) != 0 {
				t.Error("wrong empty path")
			}
		default:
			if err != nil {
				t.Error(err)
			}
			if len(path) != int(to+10-from)%10+1 {
				t.Error("wrong path length")
			}
			if path[0].from != from || path[len(path)-1].to != to {
				t.Error("wrong path")
			}
		}
	}
}
