package hjr

/*
func TestParse(t *testing.T) {
	a := assert.New(t)
	const separator byte = 1
	tests := []struct {
		name     string
		data     []byte
		expected []Cell0
	}{
		{
			name: "2 columns both are not empty",
			data: []byte{^separator, separator, ^separator},
			expected: []Cell0{
				{
					position:  1,
					firstByte: 0,
					lastByte:  0,
				},
				{
					position:  2,
					firstByte: 2,
					lastByte:  2,
				},
			},
		},
		{
			name: "2 columns: empty, not empty",
			data: []byte{separator, ^separator},
			expected: []Cell0{
				{position: 2,
					firstByte: 1,
					lastByte:  1,
				},
			},
		},
		{
			name: "2 columns: not empty, empty",
			data: []byte{^separator, separator},
			expected: []Cell0{
				{position: 1,
					firstByte: 0,
					lastByte:  0,
				},
			},
		},
		{
			name:     "2 columns:  empty, empty",
			data:     []byte{separator},
			expected: []Cell0{},
		},
		{
			name:     "3 columns:  empty, empty, empty",
			data:     []byte{separator, separator},
			expected: []Cell0{},
		},
		{
			name:     "no columns:",
			data:     []byte{},
			expected: []Cell0{},
		},
		{
			name:     "nil columns:",
			data:     nil,
			expected: []Cell0{},
		},
		{
			name: "3 columns: not empty, empty, empty",
			data: []byte{^separator, separator, separator},
			expected: []Cell0{
				{position: 1,
					firstByte: 0,
					lastByte:  0,
				},
			},
		},
		{
			name: "3 columns: not empty, empty, empty",
			data: []byte{separator, ^separator, separator},
			expected: []Cell0{
				{position: 2,
					firstByte: 1,
					lastByte:  1,
				},
			},
		},
		{
			name: "3 columns:  empty,  2x not empty, empty",
			data: []byte{separator, ^separator, ^separator, separator},
			expected: []Cell0{
				{position: 2,
					firstByte: 1,
					lastByte:  2,
				},
			},
		},
		{
			name: "3 columns:  2x not empty,  2x not empty, 3x not empty",
			data: []byte{^separator, ^separator, separator, ^separator, ^separator, separator, ^separator, ^separator},
			expected: []Cell0{
				{position: 1,
					firstByte: 0,
					lastByte:  1,
				},
				{position: 2,
					firstByte: 3,
					lastByte:  4,
				},
				{position: 3,
					firstByte: 6,
					lastByte:  7,
				},
			},
		},
		{
			name: "5 columns:  empty,empty, 2x not empty,  2x not empty, 3x not empty",
			data: []byte{separator, separator, ^separator, ^separator, separator, ^separator, ^separator, separator, ^separator, ^separator},
			expected: []Cell0{
				{position: 3,
					firstByte: 2,
					lastByte:  3,
				},
				{position: 4,
					firstByte: 5,
					lastByte:  6,
				},
				{position: 5,
					firstByte: 8,
					lastByte:  9,
				},
			},
		},
		{
			name: "5 columns:  empty,empty, 2x not empty,  2x not empty, 3x not empty",
			data: []byte{separator, separator, ^separator, ^separator, separator, separator, ^separator, ^separator, separator, ^separator, ^separator},
			expected: []Cell0{
				{position: 3,
					firstByte: 2,
					lastByte:  3,
				},
				{position: 5,
					firstByte: 6,
					lastByte:  7,
				},
				{position: 6,
					firstByte: 9,
					lastByte:  10,
				},
			},
		},
	}
	for _, test := range tests {
		a.Equal(test.expected, parse(test.data, separator), test.name)
	}
}
*/
