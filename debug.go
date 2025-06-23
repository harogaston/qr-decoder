package main

import (
	"fmt"
	"strings"
)

func print_bit_seq(bits bit_seq) {
	fmt.Println(uint_to_string(bits.value, bits.size))
}

func print_uint_as_string(bits uint) {
	fmt.Println(uint_to_string(bits, 0))
}

func print_uint_as_modules(bits uint) {
	m := uint_to_modules(bits, 0)
	print_module_slice(m)
}

func print_module_slice(mods []module) {
	b := strings.Builder{}
	for _, m := range mods {
		b.WriteString(m.Char())
	}
	fmt.Println(b.String())
}
