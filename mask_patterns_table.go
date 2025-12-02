package main

const format_information_mask uint = 0b101010000010010

// For a given mask any modules for which the condition is true
// is defined as dark
type mask struct {
	bits      int
	condition func(i, j int) Bit
}

var mask_patterns []mask = []mask{
	{bits: 0b000, condition: func(i, j int) Bit {
		return Bit((i + j) % 2)
	}},
	{bits: 0b001, condition: func(i, j int) Bit {
		return Bit(i % 2)
	}},
	{bits: 0b010, condition: func(i, j int) Bit {
		return Bit(j % 3)
	}},
	{bits: 0b011, condition: func(i, j int) Bit {
		return Bit((i + j) % 3)
	}},
	{bits: 0b100, condition: func(i, j int) Bit {
		return Bit(((i / 2) + (j / 3)) % 2)
	}},
	{bits: 0b101, condition: func(i, j int) Bit {
		return Bit((i*j)%2 + (i*j)%3)
	}},
	{bits: 0b110, condition: func(i, j int) Bit {
		return Bit(((i*j)%2 + (i*j)%3) % 2)
	}},
	{bits: 0b111, condition: func(i, j int) Bit {
		return Bit(((i+j)%2 + (i*j)%3) % 2)
	}},
}

func get_mask_pattern_for_mask(mask int) mask {
	return mask_patterns[mask]
}
