package main

var error_correction_codes map[errcorr]uint = map[errcorr]uint{
	ERR_CORR_L: 0b01,
	ERR_CORR_M: 0b00,
	ERR_CORR_Q: 0b11,
	ERR_CORR_H: 0b10,
}

func get_err_corr_bits_for_level(err_corr_level errcorr) uint {
	return error_correction_codes[err_corr_level]
}
