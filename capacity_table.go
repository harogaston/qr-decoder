package main

import (
	"github.com/harogaston/qr-decoder/version"
)

// BlockGroup represents a group of blocks with the same characteristics
type BlockGroup struct {
	NumBlocks      int
	TotalCodewords int
	DataCodewords  int
}

// ECInfo stores error correction information for a specific level
type ECInfo struct {
	TotalECCodewords int
	BlockGroups      []BlockGroup
}

var microCapacityData = map[int]struct {
	totalCodewords int
	ecInfo         map[errcorr]ECInfo
}{
	1: {
		totalCodewords: 5,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 2,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 5, DataCodewords: 3},
				},
			},
		},
	},
	2: {
		totalCodewords: 10,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 5,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 10, DataCodewords: 5},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 6,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 10, DataCodewords: 4},
				},
			},
		},
	},
	3: {
		totalCodewords: 17,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 6,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 17, DataCodewords: 11},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 8,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 17, DataCodewords: 9},
				},
			},
		},
	},
	4: {
		totalCodewords: 24,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 8,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 24, DataCodewords: 16},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 10,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 24, DataCodewords: 14},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 14,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 14, DataCodewords: 10},
				},
			},
		},
	},
}

// FIXME: Table values are wrong from level 16 and above
var capacityData = map[int]struct {
	totalCodewords int
	ecInfo         map[errcorr]ECInfo
}{
	1: {
		totalCodewords: 26,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 7,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 26, DataCodewords: 19},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 10,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 26, DataCodewords: 16},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 13,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 26, DataCodewords: 13},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 17,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 26, DataCodewords: 9},
				},
			},
		},
	},
	2: {
		totalCodewords: 44,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 10,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 44, DataCodewords: 34},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 16,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 44, DataCodewords: 28},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 22,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 44, DataCodewords: 22},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 44, DataCodewords: 16},
				},
			},
		},
	},
	3: {
		totalCodewords: 70,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 15,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 70, DataCodewords: 55},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 26,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 70, DataCodewords: 44},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 36,
				BlockGroups: []BlockGroup{
					{NumBlocks: 2, TotalCodewords: 35, DataCodewords: 17},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 44,
				BlockGroups: []BlockGroup{
					{NumBlocks: 2, TotalCodewords: 35, DataCodewords: 13},
				},
			},
		},
	},
	4: {
		totalCodewords: 100,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 20,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 100, DataCodewords: 80},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 36,
				BlockGroups: []BlockGroup{
					{NumBlocks: 2, TotalCodewords: 50, DataCodewords: 32},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 52,
				BlockGroups: []BlockGroup{
					{NumBlocks: 2, TotalCodewords: 50, DataCodewords: 24},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 64,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 25, DataCodewords: 9},
				},
			},
		},
	},
	5: {
		totalCodewords: 134,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 26,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 134, DataCodewords: 108},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 48,
				BlockGroups: []BlockGroup{
					{NumBlocks: 2, TotalCodewords: 67, DataCodewords: 43},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 72,
				BlockGroups: []BlockGroup{
					{NumBlocks: 2, TotalCodewords: 35, DataCodewords: 15},
					{NumBlocks: 2, TotalCodewords: 34, DataCodewords: 16},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 88,
				BlockGroups: []BlockGroup{
					{NumBlocks: 2, TotalCodewords: 33, DataCodewords: 11},
					{NumBlocks: 2, TotalCodewords: 34, DataCodewords: 12},
				},
			},
		},
	},
	6: {
		totalCodewords: 172,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 36,
				BlockGroups: []BlockGroup{
					{NumBlocks: 2, TotalCodewords: 86, DataCodewords: 68},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 64,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 43, DataCodewords: 27},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 96,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 43, DataCodewords: 19},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 112,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 43, DataCodewords: 15},
				},
			},
		},
	},
	7: {
		totalCodewords: 196,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 40,
				BlockGroups: []BlockGroup{
					{NumBlocks: 2, TotalCodewords: 98, DataCodewords: 78},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 72,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 49, DataCodewords: 31},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 108,
				BlockGroups: []BlockGroup{
					{NumBlocks: 2, TotalCodewords: 32, DataCodewords: 14},
					{NumBlocks: 4, TotalCodewords: 33, DataCodewords: 15},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 130,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 39, DataCodewords: 13},
					{NumBlocks: 1, TotalCodewords: 40, DataCodewords: 14},
				},
			},
		},
	},
	8: {
		totalCodewords: 242,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 48,
				BlockGroups: []BlockGroup{
					{NumBlocks: 2, TotalCodewords: 121, DataCodewords: 97},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 88,
				BlockGroups: []BlockGroup{
					{NumBlocks: 2, TotalCodewords: 60, DataCodewords: 38},
					{NumBlocks: 2, TotalCodewords: 61, DataCodewords: 39},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 132,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 40, DataCodewords: 18},
					{NumBlocks: 2, TotalCodewords: 41, DataCodewords: 19},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 156,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 40, DataCodewords: 14},
					{NumBlocks: 2, TotalCodewords: 41, DataCodewords: 15},
				},
			},
		},
	},
	9: {
		totalCodewords: 292,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 60,
				BlockGroups: []BlockGroup{
					{NumBlocks: 2, TotalCodewords: 146, DataCodewords: 116},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 110,
				BlockGroups: []BlockGroup{
					{NumBlocks: 3, TotalCodewords: 58, DataCodewords: 36},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 160,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 36, DataCodewords: 16},
					{NumBlocks: 4, TotalCodewords: 37, DataCodewords: 17},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 192,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 36, DataCodewords: 12},
					{NumBlocks: 4, TotalCodewords: 37, DataCodewords: 13},
				},
			},
		},
	},
	10: {
		totalCodewords: 346,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 72,
				BlockGroups: []BlockGroup{
					{NumBlocks: 2, TotalCodewords: 86, DataCodewords: 68},
					{NumBlocks: 2, TotalCodewords: 87, DataCodewords: 69},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 130,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 69, DataCodewords: 43},
					{NumBlocks: 1, TotalCodewords: 70, DataCodewords: 44},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 192,
				BlockGroups: []BlockGroup{
					{NumBlocks: 6, TotalCodewords: 43, DataCodewords: 19},
					{NumBlocks: 2, TotalCodewords: 44, DataCodewords: 20},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 224,
				BlockGroups: []BlockGroup{
					{NumBlocks: 6, TotalCodewords: 43, DataCodewords: 15},
					{NumBlocks: 2, TotalCodewords: 44, DataCodewords: 16},
				},
			},
		},
	},
	11: {
		totalCodewords: 404,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 80,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 101, DataCodewords: 81},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 150,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 80, DataCodewords: 50},
					{NumBlocks: 4, TotalCodewords: 81, DataCodewords: 51},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 224,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 50, DataCodewords: 22},
					{NumBlocks: 4, TotalCodewords: 51, DataCodewords: 23},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 264,
				BlockGroups: []BlockGroup{
					{NumBlocks: 3, TotalCodewords: 36, DataCodewords: 12},
					{NumBlocks: 8, TotalCodewords: 37, DataCodewords: 13},
				},
			},
		},
	},
	12: {
		totalCodewords: 466,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 96,
				BlockGroups: []BlockGroup{
					{NumBlocks: 2, TotalCodewords: 116, DataCodewords: 92},
					{NumBlocks: 2, TotalCodewords: 117, DataCodewords: 93},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 176,
				BlockGroups: []BlockGroup{
					{NumBlocks: 6, TotalCodewords: 58, DataCodewords: 36},
					{NumBlocks: 2, TotalCodewords: 59, DataCodewords: 37},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 260,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 46, DataCodewords: 20},
					{NumBlocks: 6, TotalCodewords: 47, DataCodewords: 21},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 308,
				BlockGroups: []BlockGroup{
					{NumBlocks: 7, TotalCodewords: 42, DataCodewords: 14},
					{NumBlocks: 4, TotalCodewords: 43, DataCodewords: 15},
				},
			},
		},
	},
	13: {
		totalCodewords: 532,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 104,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 133, DataCodewords: 107},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 198,
				BlockGroups: []BlockGroup{
					{NumBlocks: 8, TotalCodewords: 59, DataCodewords: 37},
					{NumBlocks: 1, TotalCodewords: 60, DataCodewords: 38},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 288,
				BlockGroups: []BlockGroup{
					{NumBlocks: 8, TotalCodewords: 44, DataCodewords: 20},
					{NumBlocks: 4, TotalCodewords: 45, DataCodewords: 21},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 352,
				BlockGroups: []BlockGroup{
					{NumBlocks: 12, TotalCodewords: 33, DataCodewords: 11},
					{NumBlocks: 4, TotalCodewords: 34, DataCodewords: 12},
				},
			},
		},
	},
	14: {
		totalCodewords: 581,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 120,
				BlockGroups: []BlockGroup{
					{NumBlocks: 3, TotalCodewords: 145, DataCodewords: 115},
					{NumBlocks: 1, TotalCodewords: 146, DataCodewords: 116},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 216,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 64, DataCodewords: 40},
					{NumBlocks: 5, TotalCodewords: 65, DataCodewords: 41},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 320,
				BlockGroups: []BlockGroup{
					{NumBlocks: 11, TotalCodewords: 36, DataCodewords: 16},
					{NumBlocks: 5, TotalCodewords: 37, DataCodewords: 17},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 384,
				BlockGroups: []BlockGroup{
					{NumBlocks: 11, TotalCodewords: 36, DataCodewords: 12},
					{NumBlocks: 5, TotalCodewords: 37, DataCodewords: 13},
				},
			},
		},
	},
	15: {
		totalCodewords: 655,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 132,
				BlockGroups: []BlockGroup{
					{NumBlocks: 5, TotalCodewords: 109, DataCodewords: 87},
					{NumBlocks: 1, TotalCodewords: 110, DataCodewords: 88},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 240,
				BlockGroups: []BlockGroup{
					{NumBlocks: 5, TotalCodewords: 65, DataCodewords: 41},
					{NumBlocks: 5, TotalCodewords: 66, DataCodewords: 42},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 360,
				BlockGroups: []BlockGroup{
					{NumBlocks: 5, TotalCodewords: 54, DataCodewords: 24},
					{NumBlocks: 7, TotalCodewords: 55, DataCodewords: 25},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 432,
				BlockGroups: []BlockGroup{
					{NumBlocks: 11, TotalCodewords: 36, DataCodewords: 12},
					{NumBlocks: 7, TotalCodewords: 37, DataCodewords: 13},
				},
			},
		},
	},
	16: {
		totalCodewords: 733,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 144,
				BlockGroups: []BlockGroup{
					{NumBlocks: 5, TotalCodewords: 122, DataCodewords: 98},
					{NumBlocks: 1, TotalCodewords: 123, DataCodewords: 99},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 280,
				BlockGroups: []BlockGroup{
					{NumBlocks: 7, TotalCodewords: 73, DataCodewords: 45},
					{NumBlocks: 3, TotalCodewords: 74, DataCodewords: 46},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 408,
				BlockGroups: []BlockGroup{
					{NumBlocks: 15, TotalCodewords: 43, DataCodewords: 19},
					{NumBlocks: 2, TotalCodewords: 44, DataCodewords: 20},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 480,
				BlockGroups: []BlockGroup{
					{NumBlocks: 3, TotalCodewords: 45, DataCodewords: 15},
					{NumBlocks: 13, TotalCodewords: 46, DataCodewords: 16},
				},
			},
		},
	},
	17: {
		totalCodewords: 815,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 5, TotalCodewords: 136, DataCodewords: 108},
					{NumBlocks: 1, TotalCodewords: 135, DataCodewords: 107},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 5, TotalCodewords: 152, DataCodewords: 122},
					{NumBlocks: 1, TotalCodewords: 55, DataCodewords: 25},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 5, TotalCodewords: 145, DataCodewords: 115},
					{NumBlocks: 2, TotalCodewords: 45, DataCodewords: 15},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 75, DataCodewords: 47},
					{NumBlocks: 10, TotalCodewords: 74, DataCodewords: 46},
				},
			},
		},
	},
	18: {
		totalCodewords: 901,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 151, DataCodewords: 121},
					{NumBlocks: 5, TotalCodewords: 150, DataCodewords: 120},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 26,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 70, DataCodewords: 44},
					{NumBlocks: 9, TotalCodewords: 69, DataCodewords: 43},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 24,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 64, DataCodewords: 40},
					{NumBlocks: 15, TotalCodewords: 43, DataCodewords: 19},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 51, DataCodewords: 23},
					{NumBlocks: 17, TotalCodewords: 50, DataCodewords: 22},
				},
			},
		},
	},
	19: {
		totalCodewords: 991,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 142, DataCodewords: 114},
					{NumBlocks: 3, TotalCodewords: 141, DataCodewords: 113},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 26,
				BlockGroups: []BlockGroup{
					{NumBlocks: 11, TotalCodewords: 71, DataCodewords: 45},
					{NumBlocks: 3, TotalCodewords: 70, DataCodewords: 44},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 3, TotalCodewords: 147, DataCodewords: 117},
					{NumBlocks: 10, TotalCodewords: 55, DataCodewords: 25},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 26,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 48, DataCodewords: 22},
					{NumBlocks: 17, TotalCodewords: 47, DataCodewords: 21},
				},
			},
		},
	},
	20: {
		totalCodewords: 1085,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 5, TotalCodewords: 136, DataCodewords: 108},
					{NumBlocks: 3, TotalCodewords: 135, DataCodewords: 107},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 26,
				BlockGroups: []BlockGroup{
					{NumBlocks: 13, TotalCodewords: 68, DataCodewords: 42},
					{NumBlocks: 3, TotalCodewords: 67, DataCodewords: 41},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 3, TotalCodewords: 147, DataCodewords: 117},
					{NumBlocks: 14, TotalCodewords: 46, DataCodewords: 16},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 5, TotalCodewords: 55, DataCodewords: 25},
					{NumBlocks: 15, TotalCodewords: 54, DataCodewords: 24},
				},
			},
		},
	},
	21: {
		totalCodewords: 1156,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 145, DataCodewords: 117},
					{NumBlocks: 4, TotalCodewords: 144, DataCodewords: 116},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 26,
				BlockGroups: []BlockGroup{
					{NumBlocks: 8, TotalCodewords: 132, DataCodewords: 106},
					{NumBlocks: 2, TotalCodewords: 50, DataCodewords: 24},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 26,
				BlockGroups: []BlockGroup{
					{NumBlocks: 17, TotalCodewords: 68, DataCodewords: 42},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 7, TotalCodewords: 73, DataCodewords: 45},
					{NumBlocks: 15, TotalCodewords: 43, DataCodewords: 15},
				},
			},
		},
	},
	22: {
		totalCodewords: 1258,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 7, TotalCodewords: 140, DataCodewords: 112},
					{NumBlocks: 2, TotalCodewords: 139, DataCodewords: 111},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 17, TotalCodewords: 74, DataCodewords: 46},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 16, TotalCodewords: 55, DataCodewords: 25},
					{NumBlocks: 7, TotalCodewords: 54, DataCodewords: 24},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 24,
				BlockGroups: []BlockGroup{
					{NumBlocks: 34, TotalCodewords: 37, DataCodewords: 13},
				},
			},
		},
	},
	23: {
		totalCodewords: 1364,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 5, TotalCodewords: 152, DataCodewords: 122},
					{NumBlocks: 4, TotalCodewords: 151, DataCodewords: 121},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 14, TotalCodewords: 76, DataCodewords: 48},
					{NumBlocks: 4, TotalCodewords: 75, DataCodewords: 47},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 14, TotalCodewords: 55, DataCodewords: 25},
					{NumBlocks: 11, TotalCodewords: 54, DataCodewords: 24},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 14, TotalCodewords: 46, DataCodewords: 16},
					{NumBlocks: 16, TotalCodewords: 45, DataCodewords: 15},
				},
			},
		},
	},
	24: {
		totalCodewords: 1474,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 148, DataCodewords: 118},
					{NumBlocks: 6, TotalCodewords: 147, DataCodewords: 117},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 14, TotalCodewords: 74, DataCodewords: 46},
					{NumBlocks: 6, TotalCodewords: 73, DataCodewords: 45},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 6, TotalCodewords: 146, DataCodewords: 116},
					{NumBlocks: 13, TotalCodewords: 46, DataCodewords: 16},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 16, TotalCodewords: 55, DataCodewords: 25},
					{NumBlocks: 11, TotalCodewords: 54, DataCodewords: 24},
				},
			},
		},
	},
	25: {
		totalCodewords: 1588,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 26,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 133, DataCodewords: 107},
					{NumBlocks: 8, TotalCodewords: 132, DataCodewords: 106},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 13, TotalCodewords: 76, DataCodewords: 48},
					{NumBlocks: 8, TotalCodewords: 75, DataCodewords: 47},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 3, TotalCodewords: 146, DataCodewords: 116},
					{NumBlocks: 25, TotalCodewords: 46, DataCodewords: 16},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 22, TotalCodewords: 55, DataCodewords: 25},
					{NumBlocks: 7, TotalCodewords: 54, DataCodewords: 24},
				},
			},
		},
	},
	26: {
		totalCodewords: 1706,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 2, TotalCodewords: 143, DataCodewords: 115},
					{NumBlocks: 10, TotalCodewords: 142, DataCodewords: 114},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 75, DataCodewords: 47},
					{NumBlocks: 19, TotalCodewords: 74, DataCodewords: 46},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 6, TotalCodewords: 51, DataCodewords: 23},
					{NumBlocks: 28, TotalCodewords: 50, DataCodewords: 22},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 47, DataCodewords: 17},
					{NumBlocks: 33, TotalCodewords: 46, DataCodewords: 16},
				},
			},
		},
	},
	27: {
		totalCodewords: 1828,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 153, DataCodewords: 123},
					{NumBlocks: 8, TotalCodewords: 152, DataCodewords: 122},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 3, TotalCodewords: 74, DataCodewords: 46},
					{NumBlocks: 22, TotalCodewords: 73, DataCodewords: 45},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 26, TotalCodewords: 54, DataCodewords: 24},
					{NumBlocks: 8, TotalCodewords: 53, DataCodewords: 23},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 10, TotalCodewords: 54, DataCodewords: 24},
					{NumBlocks: 28, TotalCodewords: 46, DataCodewords: 16},
				},
			},
		},
	},
	28: {
		totalCodewords: 1921,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 10, TotalCodewords: 148, DataCodewords: 118},
					{NumBlocks: 3, TotalCodewords: 147, DataCodewords: 117},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 23, TotalCodewords: 74, DataCodewords: 46},
					{NumBlocks: 3, TotalCodewords: 73, DataCodewords: 45},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 6, TotalCodewords: 146, DataCodewords: 116},
					{NumBlocks: 19, TotalCodewords: 55, DataCodewords: 25},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 5, TotalCodewords: 145, DataCodewords: 115},
					{NumBlocks: 26, TotalCodewords: 46, DataCodewords: 16},
				},
			},
		},
	},
	29: {
		totalCodewords: 2051,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 7, TotalCodewords: 147, DataCodewords: 117},
					{NumBlocks: 7, TotalCodewords: 146, DataCodewords: 116},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 26, TotalCodewords: 76, DataCodewords: 48},
					{NumBlocks: 1, TotalCodewords: 75, DataCodewords: 47},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 7, TotalCodewords: 74, DataCodewords: 46},
					{NumBlocks: 21, TotalCodewords: 73, DataCodewords: 45},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 19, TotalCodewords: 74, DataCodewords: 46},
					{NumBlocks: 15, TotalCodewords: 43, DataCodewords: 15},
				},
			},
		},
	},
	30: {
		totalCodewords: 2185,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 10, TotalCodewords: 146, DataCodewords: 116},
					{NumBlocks: 5, TotalCodewords: 145, DataCodewords: 115},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 10, TotalCodewords: 142, DataCodewords: 114},
					{NumBlocks: 15, TotalCodewords: 51, DataCodewords: 23},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 10, TotalCodewords: 76, DataCodewords: 48},
					{NumBlocks: 19, TotalCodewords: 75, DataCodewords: 47},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 25, TotalCodewords: 55, DataCodewords: 25},
					{NumBlocks: 15, TotalCodewords: 54, DataCodewords: 24},
				},
			},
		},
	},
	31: {
		totalCodewords: 2323,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 3, TotalCodewords: 146, DataCodewords: 116},
					{NumBlocks: 13, TotalCodewords: 145, DataCodewords: 115},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 29, TotalCodewords: 75, DataCodewords: 47},
					{NumBlocks: 2, TotalCodewords: 74, DataCodewords: 46},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 55, DataCodewords: 25},
					{NumBlocks: 42, TotalCodewords: 54, DataCodewords: 24},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 28, TotalCodewords: 46, DataCodewords: 16},
					{NumBlocks: 23, TotalCodewords: 45, DataCodewords: 15},
				},
			},
		},
	},
	32: {
		totalCodewords: 2465,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 17, TotalCodewords: 145, DataCodewords: 115},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 23, TotalCodewords: 75, DataCodewords: 47},
					{NumBlocks: 10, TotalCodewords: 74, DataCodewords: 46},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 5, TotalCodewords: 152, DataCodewords: 122},
					{NumBlocks: 31, TotalCodewords: 55, DataCodewords: 25},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 35, TotalCodewords: 55, DataCodewords: 25},
					{NumBlocks: 10, TotalCodewords: 54, DataCodewords: 24},
				},
			},
		},
	},
	33: {
		totalCodewords: 2611,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 1, TotalCodewords: 146, DataCodewords: 116},
					{NumBlocks: 17, TotalCodewords: 145, DataCodewords: 115},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 21, TotalCodewords: 75, DataCodewords: 47},
					{NumBlocks: 14, TotalCodewords: 74, DataCodewords: 46},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 6, TotalCodewords: 151, DataCodewords: 121},
					{NumBlocks: 31, TotalCodewords: 55, DataCodewords: 25},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 5, TotalCodewords: 145, DataCodewords: 115},
					{NumBlocks: 41, TotalCodewords: 46, DataCodewords: 16},
				},
			},
		},
	},
	34: {
		totalCodewords: 2761,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 6, TotalCodewords: 146, DataCodewords: 116},
					{NumBlocks: 13, TotalCodewords: 145, DataCodewords: 115},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 3, TotalCodewords: 135, DataCodewords: 107},
					{NumBlocks: 31, TotalCodewords: 76, DataCodewords: 48},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 23, TotalCodewords: 75, DataCodewords: 47},
					{NumBlocks: 14, TotalCodewords: 74, DataCodewords: 46},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 7, TotalCodewords: 55, DataCodewords: 25},
					{NumBlocks: 44, TotalCodewords: 54, DataCodewords: 24},
				},
			},
		},
	},
	35: {
		totalCodewords: 2876,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 17, TotalCodewords: 152, DataCodewords: 122},
					{NumBlocks: 2, TotalCodewords: 146, DataCodewords: 116},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 26, TotalCodewords: 76, DataCodewords: 48},
					{NumBlocks: 12, TotalCodewords: 75, DataCodewords: 47},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 152, DataCodewords: 122},
					{NumBlocks: 42, TotalCodewords: 54, DataCodewords: 24},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 5, TotalCodewords: 152, DataCodewords: 122},
					{NumBlocks: 46, TotalCodewords: 46, DataCodewords: 16},
				},
			},
		},
	},
	36: {
		totalCodewords: 3034,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 14, TotalCodewords: 152, DataCodewords: 122},
					{NumBlocks: 6, TotalCodewords: 151, DataCodewords: 121},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 20, TotalCodewords: 147, DataCodewords: 117},
					{NumBlocks: 2, TotalCodewords: 47, DataCodewords: 17},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 17, TotalCodewords: 152, DataCodewords: 122},
					{NumBlocks: 10, TotalCodewords: 45, DataCodewords: 15},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 34, TotalCodewords: 76, DataCodewords: 48},
					{NumBlocks: 6, TotalCodewords: 75, DataCodewords: 47},
				},
			},
		},
	},
	37: {
		totalCodewords: 3196,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 153, DataCodewords: 123},
					{NumBlocks: 17, TotalCodewords: 152, DataCodewords: 122},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 14, TotalCodewords: 75, DataCodewords: 47},
					{NumBlocks: 29, TotalCodewords: 74, DataCodewords: 46},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 151, DataCodewords: 121},
					{NumBlocks: 48, TotalCodewords: 54, DataCodewords: 24},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 10, TotalCodewords: 55, DataCodewords: 25},
					{NumBlocks: 49, TotalCodewords: 54, DataCodewords: 24},
				},
			},
		},
	},
	38: {
		totalCodewords: 3362,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 18, TotalCodewords: 153, DataCodewords: 123},
					{NumBlocks: 4, TotalCodewords: 152, DataCodewords: 122},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 19, TotalCodewords: 148, DataCodewords: 118},
					{NumBlocks: 10, TotalCodewords: 55, DataCodewords: 25},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 32, TotalCodewords: 75, DataCodewords: 47},
					{NumBlocks: 13, TotalCodewords: 74, DataCodewords: 46},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 14, TotalCodewords: 55, DataCodewords: 25},
					{NumBlocks: 48, TotalCodewords: 54, DataCodewords: 24},
				},
			},
		},
	},
	39: {
		totalCodewords: 3532,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 4, TotalCodewords: 148, DataCodewords: 118},
					{NumBlocks: 20, TotalCodewords: 147, DataCodewords: 117},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 19, TotalCodewords: 148, DataCodewords: 118},
					{NumBlocks: 16, TotalCodewords: 45, DataCodewords: 15},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 14, TotalCodewords: 152, DataCodewords: 122},
					{NumBlocks: 26, TotalCodewords: 54, DataCodewords: 24},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 7, TotalCodewords: 76, DataCodewords: 48},
					{NumBlocks: 40, TotalCodewords: 75, DataCodewords: 47},
				},
			},
		},
	},
	40: {
		totalCodewords: 3706,
		ecInfo: map[errcorr]ECInfo{
			ERR_CORR_L: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 6, TotalCodewords: 149, DataCodewords: 119},
					{NumBlocks: 19, TotalCodewords: 148, DataCodewords: 118},
				},
			},
			ERR_CORR_M: {
				TotalECCodewords: 28,
				BlockGroups: []BlockGroup{
					{NumBlocks: 31, TotalCodewords: 76, DataCodewords: 48},
					{NumBlocks: 18, TotalCodewords: 75, DataCodewords: 47},
				},
			},
			ERR_CORR_Q: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 34, TotalCodewords: 55, DataCodewords: 25},
					{NumBlocks: 34, TotalCodewords: 54, DataCodewords: 24},
				},
			},
			ERR_CORR_H: {
				TotalECCodewords: 30,
				BlockGroups: []BlockGroup{
					{NumBlocks: 61, TotalCodewords: 46, DataCodewords: 16},
					{NumBlocks: 20, TotalCodewords: 45, DataCodewords: 15},
				},
			},
		},
	},
}

func getTotalCodewords(v version.QRVersion) int {
	switch v.Format {
	case version.FORMAT_MICRO_QR:
		data := microCapacityData[v.Number]
		return data.totalCodewords
	case version.FORMAT_QR:
		fallthrough
	case version.FORMAT_QR_MODEL_2:
		data := capacityData[v.Number]
		return data.totalCodewords
	}
	return 0
}

func getTotalECCodewords(v version.QRVersion, ecLevel errcorr) int {
	var ecInfo ECInfo
	switch v.Format {
	case version.FORMAT_MICRO_QR:
		capacityData := microCapacityData[v.Number]
		ecInfo = capacityData.ecInfo[ecLevel]
	case version.FORMAT_QR:
		fallthrough
	case version.FORMAT_QR_MODEL_2:
		capacityData := capacityData[v.Number]
		ecInfo = capacityData.ecInfo[ecLevel]
	}

	return ecInfo.TotalECCodewords
}

func getTotalDataCodewords(v version.QRVersion, ecLevel errcorr) int {
	return getTotalCodewords(v) - getTotalECCodewords(v, ecLevel)
}
