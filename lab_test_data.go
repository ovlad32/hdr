package hjr

import "context"

type stitchTest struct {
	name     string
	data     []int
	pos1     int
	pos2     int
	expected []int
}

var stitchTests = []stitchTest{
	{name: "Initializie",
		pos1:     2,
		pos2:     1,
		expected: []int{1, 2},
	},
	{name: "Appending",
		data:     []int{1, 2},
		pos2:     3,
		expected: []int{1, 2, 3},
	},
	{name: "Inserting in the middle",
		data:     []int{1, 4},
		pos2:     3,
		expected: []int{1, 3, 4},
	},
	{name: "Inserting at header",
		data:     []int{3, 4},
		pos2:     2,
		expected: []int{2, 3, 4},
	},
	{name: "Duplicate",
		data:     []int{3, 4},
		pos2:     3,
		expected: []int{3, 4},
	},
}

type parseTest struct {
	name     string
	sep      []byte
	data     []byte
	expected Cells
}

var parseTests = []parseTest{
	{
		name:     "No data",
		sep:      []byte("whatever"),
		data:     []byte{},
		expected: Cells{},
	},
	{
		name: "1 cell",
		sep:  []byte("_"),
		data: []byte("whatever"),

		expected: Cells{
			{value: []byte("whatever") /* hash: 0x9e31996f, */, pos: 1},
		},
	},
	{
		name: "2 cells",
		sep:  []byte("_"),
		data: []byte("value1_value2"),
		expected: Cells{
			{value: []byte("value1"), pos: 1},
			{value: []byte("value2"), pos: 2},
		},
	},
	{
		name: "3 cells, one empty",
		sep:  []byte("_"),
		data: []byte("value1__value3"),
		expected: Cells{
			{value: []byte("value1") /* hash: 1984410657, */, pos: 1},
			{value: []byte("value3") /* hash: 1950855419, */, pos: 3},
		},
	},
	{
		name: "3 cells, 2 empty",
		sep:  []byte("_"),
		data: []byte("_value2_"),
		expected: Cells{
			{value: []byte("value2") /* hash:1934077800, */, pos: 2},
		},
	},
	{
		name: "5 cells, 3 empty",
		sep:  []byte("_"),
		data: []byte("___value4_value5"),
		expected: Cells{
			{value: []byte("value4") /* hash:2034743514, */, pos: 4},
			{value: []byte("value5") /* hash:2051521133, */, pos: 5},
		},
	},
	{
		name: "3 cells, 2 dups in a row",
		sep:  []byte("_"),
		data: []byte("value0_value1_value1"),
		expected: Cells{
			{value: []byte("value0") /*  hash:1967633038,  */, pos: 1},
			{value: []byte("value1") /*  hash:1984410657,  */, pos: 2, poss: []int{2, 3}},
			{value: []byte("value1") /*  hash:1984410657,  */, pos: 0},
		},
	},
	{
		name: "3 cells, 2 dups separately",
		sep:  []byte("_"),
		data: []byte("value0_value1_value0"),
		expected: Cells{
			{value: []byte("value0") /*  hash:1967633038,  */, pos: 1, poss: []int{1, 3}},
			{value: []byte("value0") /*  hash:1967633038,  */, pos: 0},
			{value: []byte("value1") /*  hash:1984410657,  */, pos: 2},
		},
	},
	{
		name: "7 cells, 2+2+3 dups",
		sep:  []byte("_"),
		data: []byte("value0_value1_value0_value2_value1_value2_value0"),
		expected: Cells{
			{value: []byte("value0") /* hash:1967633038, */, pos: 1, poss: []int{1, 3, 7}},
			{value: []byte("value0") /* hash:1967633038, */, pos: 0},
			{value: []byte("value1") /* hash:1984410657, */, pos: 2, poss: []int{2, 5}},
			{value: []byte("value1") /* hash:1984410657, */, pos: 0},
			{value: []byte("value2") /* hash:1934077800, */, pos: 4, poss: []int{4, 6}},
			{value: []byte("value2") /* hash:1934077800, */, pos: 0},
			{value: []byte("value0") /* hash:1967633038, */, pos: 0},
		},
	},
	{
		name: "7 cells, reverse order, 2+2+3 dups",
		sep:  []byte("_"),
		data: []byte("value8_value7_value6_value6_value7_value8_value6"),
		expected: Cells{
			{value: []byte("value6") /* hash:2001188276, */, pos: 3, poss: []int{3, 4, 7}},
			{value: []byte("value6") /* hash:2001188276, */, pos: 0},
			{value: []byte("value6") /* hash:2001188276, */, pos: 0},
			{value: []byte("value7") /* hash:2017965895, */, pos: 2, poss: []int{2, 5}},
			{value: []byte("value8") /* hash:1833412086, */, pos: 1, poss: []int{1, 6}},
			{value: []byte("value8") /* hash:1833412086, */, pos: 0},
			{value: []byte("value7") /* hash:2017965895, */, pos: 0},
		},
	},
}

/*
type posBoundsTest struct {
	name     string
	cells    Cells
	start    int
	gap      bool
	expected []int
}


var posBoundsTests = []posBoundsTest{
	{
		name: "No data",
		cells: Cells{},
		expected: []int{NotFound, NotFound},
	},
	{
		name: "All items set",
		cells: Cells{
			{pos: 1},
			{pos: 2},
			{pos: 3},
		},
		expected: []int{0, 2},
	},
	{
		name: "3 items set",
		cells: Cells{
			{pos: 1},
			{pos: 2},
			{pos: 3},
			{pos: 0},
		},
		expected: []int{0, 2},
	},
}
*/

type removePosGapsTest struct {
	name     string
	data     Cells
	expected Cells
}

var removePosGapsTests = []removePosGapsTest{
	{name: "empty", data: Cells{}, expected: Cells{}},
	{
		name: "No gaps",
		data: Cells{
			{pos: 1}, {pos: 2}, {pos: 3},
		},
		expected: Cells{
			{pos: 1}, {pos: 2}, {pos: 3},
		},
	},
	{
		name: "All gaps",
		data: Cells{
			{pos: Gap}, {pos: Gap}, {pos: Gap},
		},
		expected: Cells{},
	},
	{
		name: "1 gap in array of 3",
		data: Cells{
			{pos: 1}, {pos: Gap}, {pos: 3},
		},
		expected: Cells{{pos: 1}, {pos: 3}},
	},
	{
		name: "1 item in array of 3",
		data: Cells{
			{pos: Gap}, {pos: 1}, {pos: Gap},
		},
		expected: Cells{{pos: 1}},
	},
	{
		name: "3 item in array of 6",
		data: Cells{
			{pos: Gap}, {pos: 1}, {pos: 2}, {pos: Gap}, {pos: Gap}, {pos: 3},
		},
		expected: Cells{{pos: 1}, {pos: 2}, {pos: 3}},
	},
}

type clusterSizeTest struct {
	name     string
	ctx      context.Context
	data     Cells
	size     int
	expected Clusters
}

var clusterSizeTests = []clusterSizeTest{
	{
		name:     "No Cells - No Clusters",
		ctx:      context.Background(),
		data:     Cells{},
		expected: Clusters{},
	},
	{
		name: "2 Cells - 1 Cluster",
		ctx:  context.Background(),
		data: Cells{
			{
				pos:  1,
				hash: 11,
			},
			{
				pos:  2,
				hash: 22,
			},
		},
		size: 2,
		expected: Clusters{
			{poss: []int{1, 2}, hashSum: 33, size: 2},
		},
	},

	{
		name: "3 Cells - 2 Clusters",
		ctx:  context.Background(),
		data: Cells{
			{
				pos:   1,
				hash:  11,
				value: []byte("111"),
			},
			{
				pos:   2,
				hash:  22,
				value: []byte("222")},

			{
				pos:   3,
				hash:  33,
				value: []byte("333")},
		},
		size: 2,
		expected: Clusters{
			{poss: []int{1, 2}, hashSum: 33, size: 2,
				values: append(
					append([]byte("111"), ClusterRawValueSeparator),
					append([]byte("222"), ClusterRawValueSeparator)...,
				),
			},
			{poss: []int{2, 3}, hashSum: 55, size: 2,
				values: append(
					append([]byte("222"), ClusterRawValueSeparator),
					append([]byte("333"), ClusterRawValueSeparator)...,
				),
			},
		},
	},
}

/*
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
*/
