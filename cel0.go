package hjr

/*
type Cell0 struct {
	position  int
	firstByte int
	lastByte  int
}

func parse(data []byte, sep byte) []Cell0 {
	if data == nil || len(data) == 0 {
		return make(Cells, 0)
	}

	var item Cell
	result := make([]Cell0, 0)
	for i, b := range data {
		if b == sep {
			item.position++
			if item.firstByte <= i {
				item.lastByte = i - 1
				if item.lastByte >= item.firstByte {
					result = append(result, item)
				}
			}
			item.firstByte = i + 1
		}
	}
	if item.lastByte = len(data) - 1; item.lastByte >= item.firstByte {
		item.position++
		result = append(result, item)
	}
	return result
}
*/
