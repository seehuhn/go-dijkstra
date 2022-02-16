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

// Package dijkstra implements Dijkstra's algorithm for finding shortest paths in directed graphs.
package dijkstra

import (
	"container/heap"
	"errors"
)

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

type Vertex[edge any] interface {
	comparable
	Edges() []edge
}

type Edge[cost Numeric, vertex any] interface {
	From() vertex
	To() vertex
	Cost() cost
}

type candidate[n Numeric, edge Edge[n, vertex], vertex Vertex[edge]] struct {
	to   vertex
	via  edge
	cost n
	prev *candidate[n, edge, vertex]
}

// ShortestPath returns the shortest path from the start vertex to the end vertex.
func ShortestPath[n Numeric, edge Edge[n, vertex], vertex Vertex[edge]](start, end vertex) ([]edge, error) {
	pq := &priorityQueue[n, edge, vertex]{
		index: make(map[vertex]int),
	}
	done := make(map[vertex]bool)

	best := &candidate[n, edge, vertex]{to: start}
	for best.to != end {
		done[best.to] = true

		for _, e := range best.to.Edges() {
			v := e.To()
			if done[v] {
				continue
			}

			cost := e.Cost()
			if cost <= 0 {
				return nil, ErrInvalidCost
			}
			total := best.cost + cost

			idx, ok := pq.index[v]
			if !ok {
				cand := &candidate[n, edge, vertex]{
					to:   v,
					via:  e,
					cost: total,
					prev: best,
				}
				heap.Push(pq, cand)
			} else if cand := pq.candidates[idx]; total < cand.cost {
				cand.via = e
				cand.cost = total
				cand.prev = best
				heap.Fix(pq, idx)
			}
		}

		if pq.Len() == 0 {
			return nil, ErrNoPath
		}
		best = heap.Pop(pq).(*candidate[n, edge, vertex])
	}

	var path []edge
	for best.prev != nil {
		edge := best.via
		path = append(path, edge)
		best = best.prev
	}

	reverse(path)

	return path, nil
}

// https://stackoverflow.com/a/28058324/648741
func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

type priorityQueue[n Numeric, edge Edge[n, vertex], vertex Vertex[edge]] struct {
	candidates []*candidate[n, edge, vertex]
	index      map[vertex]int
}

// Len implements heap.Interface
func (s *priorityQueue[n, e, v]) Len() int {
	return len(s.candidates)
}

// Less implements heap.Interface
func (s *priorityQueue[n, e, v]) Less(i, j int) bool {
	cand := s.candidates
	return cand[i].cost < cand[j].cost
}

// Swap implements heap.Interface
func (s *priorityQueue[n, e, v]) Swap(i, j int) {
	cand := s.candidates
	cand[i], cand[j] = cand[j], cand[i]
	s.index[cand[i].to] = i
	s.index[cand[j].to] = j
}

// Push implements heap.Interface
func (s *priorityQueue[n, edge, vertex]) Push(x interface{}) {
	cand := x.(*candidate[n, edge, vertex])
	s.index[cand.to] = len(s.candidates)
	s.candidates = append(s.candidates, cand)
}

// Pop implements heap.Interface
func (s *priorityQueue[n, e, v]) Pop() interface{} {
	length := len(s.candidates)
	x := s.candidates[length-1]
	s.candidates = s.candidates[:length-1]
	delete(s.index, x.to)
	return x
}

var (
	// ErrInvalidCost is returned when an edge has a negative or zero cost.
	ErrInvalidCost = errors.New("edge cost not strictly positive")

	// ErrNoPath is returned when no path exists from the start to the end vertex.
	ErrNoPath = errors.New("no path found")
)
