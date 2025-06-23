package main

var (
	alignment_patterns_table = [][]int{
		{}, // version 1
		{6, 18},
		{6, 22},
		{6, 26},
		{6, 30},
		{6, 34},
		{6, 22, 38}, // version 7
		{6, 24, 42},
		{6, 26, 46},
		{6, 28, 50},
		{6, 30, 54},
		{6, 32, 58},
		{6, 34, 62},
		{6, 26, 46, 66}, // version 14
		{6, 26, 48, 70},
		{6, 26, 50, 74},
		{6, 30, 54, 78},
		{6, 30, 56, 82},
		{6, 30, 58, 86},
		{6, 34, 62, 90},
		{6, 28, 50, 72, 94}, // version 21
		{6, 26, 50, 74, 98},
		{6, 30, 54, 78, 102},
		{6, 28, 54, 80, 106},
		{6, 32, 58, 84, 110},
		{6, 30, 58, 86, 114},
		{6, 34, 62, 90, 118},
		{6, 26, 50, 74, 98, 122}, // version 28
	}
)

func get_alignment_patterns_for_version(version int) [][]int {
	if version < 2 {
		return [][]int{}
	}

	coords := alignment_patterns_table[version-1]

	var pos [][]int

	for _, c1 := range coords {
		for _, c2 := range coords {
			pos = append(pos, []int{c1, c2})
		}
	}
	return pos
}
