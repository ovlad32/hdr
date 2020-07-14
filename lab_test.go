package hjr

import (
	"hash/fnv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStitch(t *testing.T) {
	assert := assert.New(t)
	for _, test := range stitchTests {
		actual := stitch(test.data, test.pos1, test.pos2)
		assert.Equal(test.expected, actual, test.name)
	}
}

func TestParse(t *testing.T) {
	assert := assert.New(t)
	hf := fnv.New32a()
	_ = hf

	for _, test := range parseTests {
		actual := parse(test.data, test.sep, nil, false)
		assert.Equal(test.expected, actual, test.name)
	}
}

/*
func TestPosBounds(t *testing.T) {
	assert := assert.New(t)
	for _, test := range posBoundsTests {
		low, high := posBounds(test.cells, test.start)
		actual := []int{low, high}
		assert.Equal(test.expected, actual, test.name)
	}
}
*/
func TestRemovePosGaps(t *testing.T) {
	assert := assert.New(t)
	for _, test := range removePosGapsTests {
		actual := removePosGaps(test.data)
		assert.Equal(test.expected, actual, test.name)
	}
}

func TestClusterizeSize(t *testing.T) {
	assert := assert.New(t)

	var clusters Clusters

	for _, test := range clusterSizeTests {
		clusters = make(Clusters, 0)
		clusterizeSize(test.data, test.size, func(c Cluster) error {
			clusters = append(clusters, c)
			return nil
		})
		assert.Equal(test.expected, clusters, test.name)
	}
}

/*
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
*/
