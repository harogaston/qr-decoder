package main

// For a given mask any modules for which the condition is true
// is defined as dark
type mask struct {
	bits   int
	maskFn func(i, j int) bool
}

var mask_patterns []mask = []mask{
	{bits: 0b000, maskFn: func(i, j int) bool {
		return ((i + j) % 2) == 0
	}},
	{bits: 0b001, maskFn: func(i, j int) bool {
		return (i % 2) == 0
	}},
	{bits: 0b010, maskFn: func(i, j int) bool {
		return (j % 3) == 0
	}},
	{bits: 0b011, maskFn: func(i, j int) bool {
		return ((i + j) % 3) == 0
	}},
	{bits: 0b100, maskFn: func(i, j int) bool {
		return (((i / 2) + (j / 3)) % 2) == 0
	}},
	{bits: 0b101, maskFn: func(i, j int) bool {
		return ((i*j)%2 + (i*j)%3) == 0
	}},
	{bits: 0b110, maskFn: func(i, j int) bool {
		return (((i*j)%2 + (i*j)%3) % 2) == 0
	}},
	{bits: 0b111, maskFn: func(i, j int) bool {
		return (((i+j)%2 + (i*j)%3) % 2) == 0
	}},
}

func get_mask_pattern_for_mask(mask int) mask {
	return mask_patterns[mask]
}
