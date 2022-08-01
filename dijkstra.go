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

// E describes the possible edge types.
type E any

// L describes the possible types for edge lengths.
// While signed types are allowed, actual edge lengths must be positive.
type L interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Graph represents a directed graph.
type Graph[vertex V, edge E, length L] interface {
	// Edge returns the outgoing edges of the given vertex.
	Edges(v vertex) []edge

	// Length returns the length of an edge.
	Length(e edge) length

	// To returns the endpoint of an edge.
	To(e edge) vertex
}

// ShortestPath returns the shortest path between two vertices.
func ShortestPath[vertex V, edge E, length L](g Graph[vertex, edge, length], start, end vertex) ([]edge, error) {
	pq := newHeap[vertex, edge, length]()

	var shortestPath *subPath[vertex, edge, length]
	var prevPathLength length

	currentVertex := start
	for currentVertex != end {
		for _, e := range g.Edges(currentVertex) {
			edgeLength := g.Length(e)
			if edgeLength < 0 {
				return nil, ErrInvalidLength
			}
			pathLength := prevPathLength + edgeLength

			v := g.To(e)
			idx, ok := pq.index[v]
			if idx < 0 { // vertex v has already been visited
				continue
			} else if !ok { // vertex seen for the first time
				cand := &subPath[vertex, edge, length]{
					to:        v,
					finalEdge: e,
					total:     pathLength,
					prev:      shortestPath,
				}
				pq.Add(cand)
			} else if cand := pq.candidates[idx]; pathLength < cand.total {
				cand.finalEdge = e
				cand.total = pathLength
				cand.prev = shortestPath
				pq.Update(idx)
			}
		}

		if len(pq.candidates) == 0 {
			return nil, ErrNoPath
		}
		shortestPath = pq.Shortest()
		currentVertex = g.To(shortestPath.finalEdge)
		prevPathLength = shortestPath.total
	}

	var path []edge
	for shortestPath != nil {
		path = append(path, shortestPath.finalEdge)
		shortestPath = shortestPath.prev
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
	ErrNoPath        = errors.New("no path exists")
)
