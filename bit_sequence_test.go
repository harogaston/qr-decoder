package main

import (
	"reflect"
	"testing"
)

func TestNewBitSeq(t *testing.T) {
	tests := []struct {
		input    uint
		expected bit_seq
	}{
		{0, bit_seq{data: []byte{0}, len: 0}},    // bits.Len(0) is 0
		{1, bit_seq{data: []byte{0x80}, len: 1}}, // 1 -> 1 bit -> 10000000
		{2, bit_seq{data: []byte{0x80}, len: 2}}, // 2 (10) -> 2 bits -> 10000000
		{3, bit_seq{data: []byte{0xC0}, len: 2}}, // 3 (11) -> 2 bits -> 11000000
		{255, bit_seq{data: []byte{0xFF}, len: 8}},
	}

	for _, test := range tests {
		// Note: bits.Len(0) is 0, so NewBitSeq(0) creates empty seq.
		// Wait, createBitSeq with size 0 makes data []byte{}.
		// Let's check implementation of createBitSeq.
		// numBytes := (0 + 7) / 8 = 0. data := make([]byte, 0).
		// So expected data for 0 should be []byte{}.

		got := NewBitSeq(test.input)
		if test.input == 0 {
			if len(got.data) != 0 || got.len != 0 {
				t.Errorf("NewBitSeq(0) = %v, want empty", got)
			}
			continue
		}

		if got.len != test.expected.len {
			t.Errorf("NewBitSeq(%d).len = %d, want %d", test.input, got.len, test.expected.len)
		}
		if !reflect.DeepEqual(got.data, test.expected.data) {
			t.Errorf("NewBitSeq(%d).data = %v, want %v", test.input, got.data, test.expected.data)
		}
	}
}

func TestNewBitSeqWithSize(t *testing.T) {
	tests := []struct {
		input    uint
		size     int
		expected bit_seq
		wantErr  bool
	}{
		{1, 4, bit_seq{data: []byte{0x10}, len: 4}, false}, // 1 -> 0001 -> 00010000 (0x10)
		{3, 4, bit_seq{data: []byte{0x30}, len: 4}, false}, // 3 (11) -> 0011 -> 00110000 (0x30)
		{1, 1, bit_seq{data: []byte{0x80}, len: 1}, false},
		{1, 0, bit_seq{}, true},                            // Size too small (bits.Len(1) is 1)
		{0, 4, bit_seq{data: []byte{0x00}, len: 4}, false}, // 0 -> 0000
	}

	for _, test := range tests {
		got, err := NewBitSeqWithSize(test.input, test.size)
		if (err != nil) != test.wantErr {
			t.Errorf("NewBitSeqWithSize(%d, %d) error = %v, wantErr %v", test.input, test.size, err, test.wantErr)
			continue
		}
		if !test.wantErr {
			if got.len != test.expected.len {
				t.Errorf("NewBitSeqWithSize(%d, %d).len = %d, want %d", test.input, test.size, got.len, test.expected.len)
			}
			if !reflect.DeepEqual(got.data, test.expected.data) {
				t.Errorf("NewBitSeqWithSize(%d, %d).data = %v, want %v", test.input, test.size, got.data, test.expected.data)
			}
		}
	}
}

func TestConcat(t *testing.T) {
	b1, _ := NewBitSeqWithSize(3, 2) // 11
	b2, _ := NewBitSeqWithSize(0, 2) // 00

	// 11 + 00 = 1100 (0xC0)
	got := Concat(b1, b2)
	if got.len != 4 {
		t.Errorf("Concat len = %d, want 4", got.len)
	}
	if got.data[0] != 0xC0 {
		t.Errorf("Concat data = %x, want C0", got.data[0])
	}

	// Test crossing byte boundary
	b3, _ := NewBitSeqWithSize(255, 8) // 11111111
	b4, _ := NewBitSeqWithSize(1, 1)   // 1
	// 11111111 + 1 = 11111111 10000000 (0xFF 0x80)
	got2 := Concat(b3, b4)
	if got2.len != 9 {
		t.Errorf("Concat len = %d, want 9", got2.len)
	}
	expected := []byte{0xFF, 0x80}
	if !reflect.DeepEqual(got2.data, expected) {
		t.Errorf("Concat data = %v, want %v", got2.data, expected)
	}
}

func TestGet(t *testing.T) {
	bs, _ := NewBitSeqWithSize(0xA, 4) // 1010
	if bs.Get(0) != 1 {
		t.Error("Get(0) != 1")
	}
	if bs.Get(1) != 0 {
		t.Error("Get(1) != 0")
	}
	if bs.Get(2) != 1 {
		t.Error("Get(2) != 1")
	}
	if bs.Get(3) != 0 {
		t.Error("Get(3) != 0")
	}
	if bs.Get(4) != 0 {
		t.Error("Get(4) != 0 (out of bounds)")
	}
}

func TestToBytes(t *testing.T) {
	bs, _ := NewBitSeqWithSize(0xFF, 8)
	bytes := bs.ToBytes()
	if len(bytes) != 1 || bytes[0] != 0xFF {
		t.Errorf("ToBytes = %v, want [0xFF]", bytes)
	}
}
