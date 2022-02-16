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

type BinGraph struct{}

type BinEdge struct {
	from, to uint32
}

func (g *BinGraph) Neighbours(v uint32) []BinEdge {
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

func (g *BinGraph) Length(e BinEdge) float64 {
	return 1 + 1/float64(e.from)
}

func (g *BinGraph) From(e BinEdge) uint32 {
	return e.from
}

func (g *BinGraph) To(e BinEdge) uint32 {
	return e.to
}

func TestBinary(t *testing.T) {
	g := &BinGraph{}
	path, err := ShortestPath[uint32, BinEdge, float64](g, 100, 1000)
	if err != nil {
		t.Fatal(err)
	}
	if len(path) < 2 || path[0].from != 100 || path[len(path)-1].to != 1000 {
		t.Error("wrong path")
	}
}

type Circle int

func (g Circle) Neighbours(n int) []CircEdge {
	var res int
	if n >= int(g) {
		res = 0
	} else {
		res = n + 1
	}
	return []CircEdge{{from: n, to: res}}
}

func (g Circle) Length(e CircEdge) int {
	return 1
}

func (g Circle) From(e CircEdge) int {
	return e.from
}

func (g Circle) To(e CircEdge) int {
	return e.to
}

type CircEdge struct {
	from, to int
}

func TestCircular(t *testing.T) {
	g := Circle(10)
	from := 10

	for _, to := range []int{9, 10, 11} {
		path, err := ShortestPath[int, CircEdge, int](g, from, to)
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
