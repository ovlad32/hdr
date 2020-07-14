package hjr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"sort"
)

type hashVal uint32
type positions []int

func (h0 hashVal) add(h1 hashVal) hashVal {
	return (h0 + h1) & 0xFFFFFFFF
}

// Cell struct {
type Cell struct {
	pos   int
	hash  hashVal
	value []byte
	poss  []int
}

type Cluster struct {
	hashSum hashVal
	size    int
	values  []byte
	poss    []int
}

type Cells = []Cell
type Clusters = []Cluster

const Gap int = 0
const NotFound int = -1
const ClusterRawValueSeparator byte = 0x1

func parse(rawData []byte, sep []byte, h hash.Hash32, orderHash bool) Cells {
	if len(rawData) == 0 {
		return Cells{}
	}
	slices := bytes.Split(rawData, sep)
	result := make(Cells, 0, len(slices))

	for i := range slices {
		if len(slices[i]) > 0 {
			result = append(result, Cell{
				value: slices[i],
				pos:   i + 1,
			})
			if h != nil {
				h.Reset()
				h.Write(slices[i])
				result[len(result)-1].hash = hashVal(h.Sum32())
			}
		}
	}

	sort.Slice(result,
		func(i, j int) bool {
			var rc int
			if result[i].pos != Gap && result[j].pos != Gap {
				if orderHash {
					rc = int(result[i].hash - result[j].hash)
				} else {
					rc = bytes.Compare(result[i].value, result[j].value)
				}
				if rc == 0 {
					// When equal, keep first and shift to the end others
					rc = -1
					if result[i].pos > result[j].pos {
						j, i, rc = i, j, 1
					}
					result[i].poss = stitch(result[i].poss, result[i].pos, result[j].pos)
					result[j].pos = Gap
				}
			} else if result[i].pos != Gap {
				rc = -1
			} else {
				rc = 1
			}
			return rc < 0
		},
	)
	return result
}

func stitch(poss []int, pos1, pos2 int) []int {
	if poss == nil {
		poss = make([]int, 2, 10)
		if pos1 < pos2 {
			poss[0] = pos1
			poss[1] = pos2
		} else {
			poss[0] = pos2
			poss[1] = pos1
		}
	} else {
		f := sort.SearchInts(poss, pos2)
		if f == len(poss) {
			poss = append(poss, pos2)
		} else if poss[f] == pos2 {
		} else {
			poss = append(poss, 0)
			copy(poss[f+1:], poss[f:len(poss)-1])
			poss[f] = pos2
		}
	}
	return poss
}

func removePosGaps(cells Cells) Cells {
	c := 0
	for i := range cells {
		if cells[i].pos != Gap {
			if c != i {
				cells[c] = cells[i]
			}
			c++
		}
	}
	return cells[:c]
}

type clusterAcceptFunc = func(c Cluster) error

func clusterizeSize(cells Cells, size int, acceptFunc clusterAcceptFunc) (err error) {
	if len(cells) == 0 {
		return nil
	}

	if len(cells) < size {
		err = fmt.Errorf("Unexpected parameters: len=%v, size=%v", len(cells), size)
		return err
	}
	upper := size - 1
	for start := 0; (start + upper) < len(cells); start++ {
		cluster := Cluster{
			poss: make([]int, 0, size),
			size: size,
		}
		slice := cells[start:(start + size)]
		valLen := 0
		for j := range slice {
			cluster.hashSum = cluster.hashSum.add(slice[j].hash)
			if slice[j].poss == nil {
				cluster.poss = append(cluster.poss, slice[j].pos)
			} else {
				cluster.poss = append(cluster.poss, slice[j].poss...)
			}
			if slice[j].value != nil {
				valLen = valLen + len(slice[j].value) + 1
			}
		}
		if valLen > 0 {
			cluster.values = make([]byte, 0, valLen)
			for j := range slice {
				cluster.values = append(cluster.values, slice[j].value...)
				cluster.values = append(cluster.values, ClusterRawValueSeparator)
			}
		}
		err = acceptFunc(cluster)
		if err != nil {
			return err
		}
	}

	return nil
}

func clusterizeRange(cells Cells, min, max int, acceptFunc clusterAcceptFunc) error {
	for size := min; size <= max; size++ {
		err := clusterizeSize(cells, size, acceptFunc)
		if err != nil {
			return err
		}
	}
	return nil
}

func serialize(c Cluster, line int, w io.Writer) (err error) {
	err = binary.Write(w, binary.BigEndian, uint64(line))
	if err != nil {
		err = fmt.Errorf("writting line number: %v", err)
		return err
	}
	err = binary.Write(w, binary.BigEndian, uint16(len(c.poss)))
	if err != nil {
		err = fmt.Errorf("writting position sequence length: %v", err)
		return err
	}
	for _, pos := range c.poss {
		err = binary.Write(w, binary.BigEndian, uint16(pos))
		if err != nil {
			err = fmt.Errorf("writting a position: %v", err)
			return err
		}
	}
	return nil
}

type flushFunc = func() error
type lookuper interface {
	lookup(hashVal) (io.Writer, flushFunc, error)
}

/*
type cellBuilderOption struct {
	sep []byte
}
type cellBuilderOptions = []cellBuilderOption

type cellBuilder struct {
	opts cellBuilderOptions
}
func (cb) Parse

type clusterBuilderOption struct {

}
*/
type IndexingOptions struct {
	hxx     hash.Hash32
	sep     []byte
	storage lookuper
}

func NewIndexOptions() *IndexingOptions {
	return &IndexingOptions{}
}
func (iop *IndexingOptions) SetHashing(hxx hash.Hash32) *IndexingOptions {
	iop.hxx = hxx
	return iop
}
func (iop *IndexingOptions) SetSeparator(sep []byte) *IndexingOptions {
	iop.sep = sep
	return iop
}

func (iop *IndexingOptions) SetStorage(l lookuper) *IndexingOptions {
	iop.storage = l
	return iop
}
func NewIndexProc(opts *IndexingOptions) *IndexProc {
	return &IndexProc{
		*opts,
	}
}

type IndexProc struct {
	IndexingOptions
}

func (ip IndexProc) ParseFunc() AcceptFunc {
	return ip.Index
}

func (ip IndexProc) Index(line int, raw []byte) error {
	af := func(c Cluster) error {
		w, flush, err := ip.storage.lookup(c.hashSum)
		if err != nil {
			return err
		}
		err = serialize(c, line, w)
		if err != nil {
			return err
		}
		flush()
		return nil
	}
	//log.Printf("idx = %v",len(raw))
	cells := parse(raw, ip.sep, ip.hxx, false)
	cells = removePosGaps(cells)
	clusterizeRange(cells, 2, 3, af)
	return nil
}

func NewMemStorage() lookuper {
	return &memStorage{
		mem: make(map[hashVal]*bytes.Buffer),
	}
}

type memStorage struct {
	mem map[hashVal]*bytes.Buffer
}

func (s *memStorage) lookup(h hashVal) (io.Writer, flushFunc, error) {
	var ok bool
	var b *bytes.Buffer

	flushFunc := func() error {
		s.mem[h] = b
		//log.Printf("l = %v",len(s.mem))
		return nil
	}

	if b, ok = s.mem[h]; !ok {
		b = bytes.NewBuffer([]byte{})
	}
	return b, flushFunc, nil
}

/*
func(w memStorage) Write(c context.Context, id string, line int, raw []byte) error {
	cells := parse(raw, w.sep)
	var eol bool
	for i:=0; !eol; i++ {
		cluster, eol = createCluster(cluster,i,3)
		hash := calcCellClusterHash(cluster)
		var ok bool
		var b *bytes.Buffer
		if b, ok = w.mp[hash]; !ok {
			b = bytes.NewBuffer([]byte{})
		}
		if err := writeCellClusters(id, line, slicedCc, b); err != nil {
			return fmt.Errorf("writting a cluster to buffer: %v",err)
		}
		w.mp[hash] = b
	}
	return nil
}*/

/*
func (r rawLineProc) accept(c context.Context, id string, line int, raw []byte)  error {
	r.lookup
}


func lookupWriter(h hashVal, table string) (w io.Writer, flush func()) {
	return
}

type words struct {
	c []CellCluster
	line int
}

func (ws words) hash() (hash uint64) {
	for i := range ws.c {
		hash = hash + ws.c[i].hash
	}
	return 0xFFFFFFFF & hash
}

func (ws words) serializeLocation(w io.Writer) (err error) {
	for i := range ws.c {
		err = binary.Write(w, binary.BigEndian, uint64(ws.line))
		if err != nil {
			err = fmt.Errorf("writting line number: %v", err)
			return err
		}
		err = binary.Write(w, binary.BigEndian, uint16(len(ws.c[i].pos)))
		if err != nil {
			err = fmt.Errorf("writting position sequence length: %v", err)
			return err
		}
		for _, pos := range ws.c[i].pos {
			err = binary.Write(w, binary.BigEndian, uint16(pos))
			if err != nil {
				err = fmt.Errorf("writting a position: %v", err)
				return err
			}
		}
	}
	return nil
}

/
/*
//CellClusterLineWriter interface
type CellClusterLineWriter interface {
	Write(ctx context.Context,  source string, line int, clusters []CellCluster,) error
}

// CellClusterReader interface {
type CellClusterReader interface {
	Read(cluster []CellCluster) error
}
*/

type indexWriter struct {
	wc io.WriteCloser
}

/*

func posBounds(cells Cells, start int) (low, high int) {
	low, high = NotFound, NotFound
	if len(cells) == 0 {
		return
	}
	gap := cells[start].pos == 0
	for index := start; index < len(cells); index++ {
		if low != NotFound && ((cells[index].pos != Gap) == gap) {
			high = index - 1
			return
		}
		if low == NotFound && ((cells[index].pos == Gap) == gap) {
			low = index
		}
	}
	if low != NotFound {
		return low, len(cells) - 1
	}
	return NotFound, NotFound
}

*/
/*
func (i indexWriter)  Write(ctx context.Context, clusters []CellCluster, line int) error {

}*/

/*

type fileWriter struct {
	sep []byte
}


func(fw fileWriter) Write(c context.Context, id string, line int, raw []byte) error {
	cells := buildCells(raw,fw.sep)
	cc := buildCellClusters(cells)
	var eol bool
	for i:=0; !eol; i++ {
		var slicedCc []CellCluster
		slicedCc, eol = sliceCellClusters(cc,i,3)
		err := flushCellClustersToTableFile(id, line, slicedCc)
		if err != nil {
			return err
		}
	}
	return nil
}


func indexPath(hash uint64) string {
	const indexRoot string = ".index"
	return path.Join(indexRoot,strconv.FormatUint(hash,36))
}

func flushCellClustersToTableFile(table string, line int, c RowValues) error {
	hash :=  calcCellClusterHash(c)
	var b bytes.Buffer
	if err := writeCellClusters(table, line, c, &b); err != nil {
		return fmt.Errorf("writting a cluster to buffer: %v",err)
	}

	pi := indexPath(hash) + "." + table
	ifl, err := os.OpenFile(pi, os.O_APPEND + os.O_CREATE, 0600 )
	if err != nil {
		return fmt.Errorf("Opening file %v: %v", pi, err)
	}
	defer ifl.Close()

	_, err =  b.WriteTo(ifl)
	if err != nil {
		return fmt.Errorf("Writting buffer to file: %v",err)
	}
	return nil

}*/

/*
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
*/
