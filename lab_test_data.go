package hjr

type buildCellsTest struct {
	name     string
	sep      []byte
	data     []byte
	expected []Cell
}

var buildCellsTests = []buildCellsTest{
	{
		name:     "No data",
		sep:      []byte("whatever"),
		data:     []byte{},
		expected: []Cell{},
	},
	{
		name: "1 cell",
		sep:  []byte("_"),
		data: []byte("whatever"),
		expected: []Cell{
			{value: []byte("whatever"), pos: 1},
		},
	},
	{
		name: "2 cells",
		sep:  []byte("_"),
		data: []byte("value1_value2"),
		expected: []Cell{
			{value: []byte("value1"), pos: 1},
			{value: []byte("value2"), pos: 2},
		},
	},
	{
		name: "3 cells, one empty",
		sep:  []byte("_"),
		data: []byte("value1__value3"),
		expected: []Cell{
			{value: []byte("value1"), pos: 1},
			{value: []byte("value3"), pos: 3},
		},
	},
	{
		name: "3 cells, 2 empty",
		sep:  []byte("_"),
		data: []byte("_value2_"),
		expected: []Cell{
			{value: []byte("value2"), pos: 2},
		},
	},
	{
		name: "5 cells, 3 empty",
		sep:  []byte("_"),
		data: []byte("___value4_value5"),
		expected: []Cell{
			{value: []byte("value4"), pos: 4},
			{value: []byte("value5"), pos: 5},
		},
	},
}

type buildCellClustersTest struct {
	name     string
	cells    []Cell
	expected []CellCluster
}

var buildCellClustersTests = []buildCellClustersTest{
	{
		name:     "No data",
		cells:    []Cell{},
		expected: []CellCluster{},
	},
	{
		name: "2 cells",
		cells: []Cell{
			{
				pos:   1,
				value: []byte("ABC"),
			},
			{
				pos:   2,
				value: []byte("ZY"),
			},
		},
		expected: []CellCluster{
			{
				pos:   []int{2},
				value: []byte("ZY"),
			},
			{
				pos:   []int{1},
				value: []byte("ABC"),
			},
		},
	},
	{
		name: "3 cells",
		cells: []Cell{
			{
				pos:   1,
				value: []byte("ABC"),
			},
			{
				pos:   2,
				value: []byte("ZY"),
			},

			{
				pos:   3,
				value: []byte("0"),
			},
		},
		expected: []CellCluster{
			{
				pos:   []int{3},
				value: []byte("0"),
			},
			{
				pos:   []int{2},
				value: []byte("ZY"),
			},
			{
				pos:   []int{1},
				value: []byte("ABC"),
			},
		},
	},
	{
		name: "3 + 2 cells",
		cells: []Cell{
			{
				pos:   1,
				value: []byte("AA"),
			},
			{
				pos:   2,
				value: []byte("BBB"),
			},
			{
				pos:   3,
				value: []byte("BBB"),
			},
			{
				pos:   4,
				value: []byte("BBB"),
			},
			{
				pos:   5,
				value: []byte("AA"),
			},
		},
		expected: []CellCluster{
			{
				pos:   []int{1, 5},
				value: []byte("AA"),
			},
			{
				pos:   []int{2, 3, 4},
				value: []byte("BBB"),
			},
		},
	},
}

type appendCellClusterTest struct {
	name         string
	cellClusters []CellCluster
	cells        []Cell
	expected     []CellCluster
}

var appendCellClusterTests = []appendCellClusterTest{
	{
		name:         "No cells",
		cellClusters: []CellCluster{},
		cells:        []Cell{},
		expected:     []CellCluster{},
	},
	{
		name:         "3 different cells",
		cellClusters: []CellCluster{},
		cells: []Cell{
			{
				pos:   1,
				value: []byte("CD"),
			},
			{
				pos:   2,
				value: []byte("BD"),
			},
			{
				pos:   3,
				value: []byte("AF"),
			},
		},
		expected: []CellCluster{
			{
				value: []byte("AF"),
				pos:   []int{3},
			},
			{
				value: []byte("BD"),
				pos:   []int{2},
			},
			{
				value: []byte("CD"),
				pos:   []int{1},
			},
		},
	},
	{
		name:         "3 same cells",
		cellClusters: []CellCluster{},
		cells: []Cell{
			{
				pos:   1,
				value: []byte("AA"),
			},
			{
				pos:   2,
				value: []byte("AA"),
			},
			{
				pos:   3,
				value: []byte("AA"),
			},
		},
		expected: []CellCluster{
			{
				value: []byte("AA"),
				pos:   []int{1, 2, 3},
			},
		},
	},
	{
		name:         "2 same cells",
		cellClusters: []CellCluster{},
		cells: []Cell{
			{
				pos:   1,
				value: []byte("AA"),
			},
			{
				pos:   2,
				value: []byte("AB"),
			},
			{
				pos:   3,
				value: []byte("AA"),
			},
		},
		expected: []CellCluster{
			{
				value: []byte("AA"),
				pos:   []int{1, 3},
			},
			{
				value: []byte("AB"),
				pos:   []int{2},
			},
		},
	},
}

type sliceCellClusterResult struct {
	start     int
	expected1 []CellCluster
	expected2 bool
}
type sliceCellClustersTest struct {
	name     string
	clusters []CellCluster
	size     int
	slices   []sliceCellClusterResult
	panic    bool
}

var sliceCellClustersTests = []sliceCellClustersTest{
	{
		name: "4 clusters, 2 slices",
		clusters: []CellCluster{
			{pos: []int{1, 2}, value: []byte("123")},
			{pos: []int{3, 4}, value: []byte("456")},
			{pos: []int{5, 6}, value: []byte("789")},
			{pos: []int{7, 8}, value: []byte("ABC")},
		},
		size: 2,
		slices: []sliceCellClusterResult{
			{
				start: 0,
				expected1: []CellCluster{
					{pos: []int{1, 2}, value: []byte("123")},
					{pos: []int{3, 4}, value: []byte("456")},
				},
				expected2: false,
			},
			{
				start: 1,
				expected1: []CellCluster{
					{pos: []int{3, 4}, value: []byte("456")},
					{pos: []int{5, 6}, value: []byte("789")},
				},
				expected2: false,
			},
			{
				start: 2,
				expected1: []CellCluster{
					{pos: []int{5, 6}, value: []byte("789")},
					{pos: []int{7, 8}, value: []byte("ABC")},
				},
				expected2: true,
			},
		},
	},
	{
		name: "4 clusters, 3 slices",
		clusters: []CellCluster{
			{pos: []int{1, 2}, value: []byte("123")},
			{pos: []int{3, 4}, value: []byte("456")},
			{pos: []int{5, 6}, value: []byte("789")},
			{pos: []int{7, 8}, value: []byte("ABC")},
		},
		size: 3,
		slices: []sliceCellClusterResult{
			{
				start: 0,
				expected1: []CellCluster{
					{pos: []int{1, 2}, value: []byte("123")},
					{pos: []int{3, 4}, value: []byte("456")},
					{pos: []int{5, 6}, value: []byte("789")},
				},
				expected2: false,
			},
			{
				start: 1,
				expected1: []CellCluster{
					{pos: []int{3, 4}, value: []byte("456")},
					{pos: []int{5, 6}, value: []byte("789")},
					{pos: []int{7, 8}, value: []byte("ABC")},
				},
				expected2: true,
			},
		},
	},
	{
		name: "3 clusters, 1 slices",
		clusters: []CellCluster{
			{pos: []int{1, 2}, value: []byte("123")},
			{pos: []int{3, 4}, value: []byte("456")},
			{pos: []int{5, 6}, value: []byte("789")},
		},
		size: 1,
		slices: []sliceCellClusterResult{
			{
				start: 0,
				expected1: []CellCluster{
					{pos: []int{1, 2}, value: []byte("123")},
				},
				expected2: false,
			},
			{
				start: 1,
				expected1: []CellCluster{
					{pos: []int{3, 4}, value: []byte("456")},
				},
				expected2: false,
			},
			{
				start: 2,
				expected1: []CellCluster{
					{pos: []int{5, 6}, value: []byte("789")},
				},
				expected2: true,
			},
		},
	},
	{
		name: "3 clusters, 4 slices",
		clusters: []CellCluster{
			{pos: []int{1, 2}, value: []byte("123")},
			{pos: []int{3, 4}, value: []byte("456")},
			{pos: []int{5, 6}, value: []byte("789")},
		},
		size: 4,
		slices: []sliceCellClusterResult{
			{
				start: 0,
				expected1: []CellCluster{
					{pos: []int{1, 2}, value: []byte("123")},
					{pos: []int{3, 4}, value: []byte("456")},
					{pos: []int{5, 6}, value: []byte("789")},
				},
				expected2: true,
			},
		},
	},
	{
		name:     "panics when negative or zero size value",
		panic:    true,
		size:     0,
		clusters: []CellCluster{},
		slices: []sliceCellClusterResult{
			{start: 0, expected1: []CellCluster{}},
		},
	},
	{
		name:     "panics when negative start ",
		panic:    true,
		size:     1,
		clusters: []CellCluster{},
		slices: []sliceCellClusterResult{
			{start: -1, expected1: []CellCluster{}},
		},
	},
}
