package hjr

import (
	"bytes"
	"fmt"
	"sort"
)

// Cell struct {
type Cell struct {
	pos   int
	value []byte
}

// CellCluster struct {
type CellCluster struct {
	pos   []int
	value []byte
}

func buildCells(data []byte, sep []byte) []Cell {
	if data == nil || len(data) == 0 {
		return make([]Cell, 0)
	}
	values := bytes.Split(data, sep)
	result := make([]Cell, 0, len(values))
	for i := range values {
		if len(values[i]) > 0 {
			result = append(result, Cell{
				value: values[i],
				pos:   i + 1,
			})
		}
	}
	return result
}

func buildCellClusters(cells []Cell) []CellCluster {
	byLength := make(map[int][]Cell)
	allLengths := make([]int, 0, len(cells))

	for i := range cells {
		l := len(cells[i].value)
		var ok bool
		var ca []Cell
		if ca, ok = byLength[l]; !ok {
			ca = make([]Cell, 0)
			allLengths = append(allLengths, l)
		}
		byLength[l] = append(ca, cells[i])
	}
	sort.Ints(allLengths)

	result := make([]CellCluster, 0, len(cells))
	for _, l := range allLengths {
		result = appendCellCluster(result, byLength[l])
	}

	return result
}

func appendCellCluster(result []CellCluster, cells []Cell) []CellCluster {
	switch len(cells) {
	case 0:
		return result
	case 1:
		return append(result, CellCluster{
			pos:   []int{cells[0].pos},
			value: cells[0].value,
		})
	}

	var dupFound bool
	sort.Slice(cells,
		func(i, j int) bool {
			rc := bytes.Compare(cells[i].value, cells[j].value)
			dupFound = dupFound || rc == 0
			return rc < 0
		},
	)

	result = append(result, CellCluster{
		value: cells[0].value,
		pos:   []int{cells[0].pos},
	})
	rest := cells[1:]
	for i := range rest {
		last := len(result) - 1
		if dupFound && 0 == bytes.Compare(result[last].value, rest[i].value) {
			result[last].pos = append(result[last].pos, rest[i].pos)
		} else {
			result = append(result, CellCluster{
				value: rest[i].value,
				pos:   []int{rest[i].pos},
			})
		}
	}
	return result
}

func sliceCellClusters(clusters []CellCluster, start, size int) ([]CellCluster, bool) {
	if size <= 0 {
		panic(fmt.Sprintf("Unexpected size value %v", size))
	}

	if start < 0 {
		panic(fmt.Sprintf("Unexpected start value %v", start))
	}

	if start == 0 && len(clusters) <= size {
		return clusters, true
	}
	finish := start + size
	return clusters[start:finish], finish >= len(clusters)
}

//CellClusterIndexWriter interface
type CellClusterIndexWriter interface {
	Write(clusters []CellCluster, line int) error
}

// CellClusterReader interface {
type CellClusterReader interface {
	Read(cluster []CellCluster) error
}
