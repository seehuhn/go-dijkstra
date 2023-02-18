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

package dijkstra

import (
	"testing"
)

type Circle int // a circle of vertices 0, 1, ..., n-1

type CircEdge struct {
	from, to int
}

func (g Circle) AppendEdges(ee []CircEdge, v int) []CircEdge {
	var res int
	if v >= int(g)-1 {
		res = 0
	} else {
		res = v + 1
	}
	return append(ee, CircEdge{from: v, to: res})
}

func (g Circle) Length(v int, e CircEdge) int {
	if e.from != v {
		panic("wrong edge")
	}
	return 1
}

func (g Circle) To(v int, e CircEdge) int {
	if e.from != v {
		panic("wrong edge")
	}
	return e.to
}

func TestCircular(t *testing.T) {
	g := Circle(10)
	from := 5

	for _, to := range []int{-1, 0, 1, 2, 3, 4, 5, 6, 7} {
		path, err := ShortestPath[int, CircEdge, int](g, from, to)
		switch {
		case to < 0:
			if err != ErrNoPath || path != nil {
				t.Error("wrong path found")
			}
			continue
		case to == from:
			if err != nil || len(path) != 0 {
				t.Error("wrong empty path")
			}
		default:
			if err != nil {
				t.Error(err)
			}
			if path[0].from != from || path[len(path)-1].to != to {
				t.Error("wrong path")
			}
		}
		if len(path) != int(to-from+10)%10 {
			t.Errorf("wrong path length %d for %d -> %d", len(path), from, to)
		}
	}
}

type BinaryTree struct{}

func (t BinaryTree) AppendEdges(ee []uint64, v uint64) []uint64 {
	return append(ee, 2*v, 2*v+1)
}

func (t BinaryTree) Length(_ uint64, e uint64) int {
	return 1
}

func (t BinaryTree) To(_ uint64, e uint64) uint64 {
	return e
}

func TestBinaryTree(t *testing.T) {
	g := &BinaryTree{}
	path, err := ShortestPath[uint64, uint64, int](g, 1, 1000)
	if err != nil {
		t.Fatal(err)
	}
	if len(path) != 9 {
		t.Error("wrong path length")
	}
	if path[0] != 3 || path[len(path)-1] != 1000 {
		t.Error("wrong path")
	}
}

func BenchmarkBinaryTree(b *testing.B) {
	g := &BinaryTree{}
	for i := 0; i < b.N; i++ {
		ShortestPath[uint64, uint64, int](g, 1, 1000)
	}
}

type walkPos struct {
	x, y int16
}

type walk struct{}

func (w walk) AppendEdges(ee []walkPos, v walkPos) []walkPos {
	return append(ee,
		walkPos{x: v.x + 1, y: v.y},
		walkPos{x: v.x, y: v.y + 1},
		walkPos{x: v.x - 1, y: v.y},
		walkPos{x: v.x, y: v.y - 1})
}

func (w walk) Length(_ walkPos, e walkPos) int {
	return 1
}

func (w walk) To(_ walkPos, e walkPos) walkPos {
	return e
}

func BenchmarkWalk(b *testing.B) {
	g := &walk{}
	for i := 0; i < b.N; i++ {
		ShortestPath[walkPos, walkPos, int](g, walkPos{x: 0, y: 0}, walkPos{x: 4, y: 6})
	}
}

type FunnyGraph struct{}

type FunnyEdge struct {
	from, to uint32
}

func (g *FunnyGraph) AppendEdges(ee []FunnyEdge, v uint32) []FunnyEdge {
	ee = append(ee, FunnyEdge{from: v, to: v + 1})
	if v > 0 {
		ee = append(ee, FunnyEdge{from: v, to: v - 1})
	}
	if v > 0 && v%2 == 0 {
		ee = append(ee, FunnyEdge{from: v, to: v / 2})
	}
	ee = append(ee, FunnyEdge{from: v, to: 2 * v})
	return ee
}

func (g *FunnyGraph) Length(v uint32, e FunnyEdge) float64 {
	if e.from != v {
		panic("wrong edge")
	}
	return 1 + 1/float64(e.from)
}

func (g *FunnyGraph) To(v uint32, e FunnyEdge) uint32 {
	if e.from != v {
		panic("wrong edge")
	}
	return e.to
}

func TestFunny(t *testing.T) {
	g := &FunnyGraph{}
	path, err := ShortestPath[uint32, FunnyEdge, float64](g, 100, 1000)
	if err != nil {
		t.Fatal(err)
	}
	if len(path) < 2 || path[0].from != 100 || path[len(path)-1].to != 1000 {
		t.Error("wrong path")
	}
	for i := 1; i < len(path); i++ {
		if path[i].from != path[i-1].to {
			t.Error("wrong path")
		}
	}
}
