package hjr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildCells(t *testing.T) {
	assert := assert.New(t)

	for _, test := range buildCellsTests {
		actual := buildCells(test.data, test.sep)
		assert.Equal(test.expected, actual, test.name)
	}
}
func Test_buildCellClusters(t *testing.T) {
	assert := assert.New(t)
	for _, test := range buildCellClustersTests {
		actual := buildCellClusters(test.cells)
		assert.Equal(test.expected, actual, test.name)
	}
}

func Test_appendCellCluster(t *testing.T) {
	assert := assert.New(t)
	for _, test := range appendCellClusterTests {
		actual := appendCellCluster(test.cellClusters, test.cells)
		assert.Equal(test.expected, actual, test.name)
	}
}

func Test_sliceCellClusters(t *testing.T) {
	assert := assert.New(t)

	for _, test := range sliceCellClustersTests {
		for _, wanted := range test.slices {
			if test.panic {
				assert.Panics(
					func() {
						sliceCellClusters(test.clusters, wanted.start, test.size)
					},
					fmt.Sprintf("%v, when start = %v", test.name, wanted.start),
				)
			} else {
				actual1, actual2 := sliceCellClusters(test.clusters, wanted.start, test.size)
				assert.Equal(wanted.expected1, actual1, fmt.Sprintf("%v, when start = %v", test.name, wanted.start))
				assert.Equal(wanted.expected2, actual2, fmt.Sprintf("%v, when start = %v", test.name, wanted.start))

			}

		}
	}
}
