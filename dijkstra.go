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

// Package dijkstra implements Dijkstra's algorithm for shortest paths in directed graphs.
package dijkstra

import (
	"errors"
)

// V describes the possible vertex types.
type V comparable

// E describes the possible edge types
type E any

// L describes the possible types for edge lengths.
type L interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Graph represents a directed graph.
type Graph[vertex V, edge E, length L] interface {
	// Edge returns the outgoing edges from the given vertex.
	Edges(v vertex) []edge

	// Length returns the length of the given edge.
	Length(e edge) length

	// To returns the target vertex of the given edge.
	To(e edge) vertex
}

// ShortestPath returns the shortest path between two vertices.
func ShortestPath[vertex V, edge E, length L](g Graph[vertex, edge, length], start, end vertex) ([]edge, error) {
	pq := &heap[vertex, edge, length]{
		index: make(map[vertex]int),
	}
	done := make(map[vertex]bool)

	var best *candidate[vertex, edge, length]
	to := start
	var prevTotal length
	for to != end {
		done[to] = true

		for _, e := range g.Edges(to) {
			v := g.To(e)
			if done[v] {
				continue
			}

			l := g.Length(e)
			if l < 0 {
				return nil, ErrInvalidLength
			}
			total := prevTotal + l

			idx, ok := pq.index[v]
			if !ok {
				cand := &candidate[vertex, edge, length]{
					to:    v,
					via:   e,
					total: total,
					prev:  best,
				}
				pq.Push(cand)
			} else if cand := pq.candidates[idx]; total < cand.total {
				cand.via = e
				cand.total = total
				cand.prev = best

				pq.Fix(idx)
			}
		}

		if len(pq.candidates) == 0 {
			return nil, ErrNoPath
		}
		best = pq.Pop()
		to = g.To(best.via)
		prevTotal = best.total
	}

	var path []edge
	for best != nil {
		edge := best.via
		path = append(path, edge)
		best = best.prev
	}

	// reverse the path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, nil
}

// Errors returned by ShortestPath.
var (
	ErrInvalidLength = errors.New("negative edge length")
	ErrNoPath        = errors.New("path does not exist")
)
