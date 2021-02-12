package template

import (
	"sort"

	"github.com/drgo/core/errors"
)

func topoSort(m map[string][]string) []string {
	var sorted []string
	seen := make(map[string]bool)
	var visitAll func(items []string)

	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				sorted = append(sorted, item)
			}
		}
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	visitAll(keys)
	return sorted
}

// KahnSort performs a topolocial sort on a map of nodes to their adjacency list
func KahnSort(graph map[string][]string) ([]string, []string, error) {
	indegrees := make(map[string]int, len(graph)*2) //prealloc reasonably-sized map to hold # edges/node
	// calculate indegree (# of incoming (parent) edges) of all nodes
	for node, edges := range graph {
		// add this node with default indegree=0 (leaf) if it is not already added
		if _, seen := indegrees[node]; !seen {
			indegrees[node] = 0
		}
			for _, v := range edges{
				indegrees[v]++
			}
	}

	// move all nodes with indegree==0 (leaf nodes) into a queue
	var queue []string
	for node, ndeg := range indegrees{
	  if ndeg == 0 {
			queue = append(queue, node)
		}
	}

	var sorted []string
	for len(queue) > 0 {
		// move a leaf from the queue into sorted 
		node := queue[len(queue)-1]
		queue = queue[:(len(queue) - 1)]
		sorted = append(sorted, node)

		// delete all incoming edgees to a node (by decreasing its indegree). 
		// If indegree becomes 0, add to queue to be eventually moved 
		// to sorted
		for _, v := range graph[node] {
			indegrees[v]--
			if indegrees[v] == 0 {
				queue = append(queue, v)
			}
		}
	}
  // Any node with indegree > 0 has edges that cannot be removed ie cyclic.
	var cyclic []string
	for node, indegree := range indegrees {
		if indegree > 0 {
		   cyclic = append(cyclic, node)
		}
	}
	var err error
	if len(cyclic)> 0 {
    err = errors.Errorf("cycles detected")
	}
	return sorted, cyclic, err
}
