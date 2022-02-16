package dijkstra

import "container/heap"

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

type Vertex interface {
	comparable
}

type Edge[n Numeric] interface {
	Cost() n
}

func ShortestPath[n Numeric, e Edge[n], v Vertex](start, end v) ([]Edge[n], n) {
	dist := make(map[v]n)
	dist[start] = 0

	cc := &pq[n, e]{}
	heap.Push(cc, start)

	return nil, 0
}

type candidate[n Numeric, e Edge[n], v Vertex] struct {
	vertex v
	from   e
	cost   n
}

type pq[n Numeric, e Edge[n], v Vertex] struct {
	candidates []*candidate[n, e, v]
}

// Len implements heap.Interface
func (s *pq[n, e]) Len() int {
	return len(s.candidates)
}

// Less implements heap.Interface
func (s *pq[n, e]) Less(i, j int) bool {
	cand := s.candidates
	return cand[i].Cost() < cand[j].Cost()
}

// Swap implements heap.Interface
func (s *pq[n, e]) Swap(i, j int) {
	cand := s.candidates
	cand[i], cand[j] = cand[j], cand[i]
}

// Push implements heap.Interface
func (s *pq[n, e]) Push(x interface{}) {
	edge := x.(e)
	s.candidates = append(s.candidates, edge)
}

// Pop implements heap.Interface
func (s *pq[_, e]) Pop() interface{} {
	n := len(s.candidates)
	x := s.candidates[n-1]
	s.candidates = s.candidates[:n-1]
	return x
}
