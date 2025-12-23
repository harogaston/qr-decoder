qr-decoder
====

# Introduction

There are four technically different, but closely related members of the QR code family, which represent an evolutionary sequence.

- QR code model 1 is the original specification for QR code and is described in AIM ITS 97-001 the largest version of which is 14 (73 x 73 modules).

- QR code model 2 is an enhanced form of the symbology with additional features (primarily, the addition of alignment patterns to assist navigation in larger symbols) and is the basis of the first edition of this document (i.e. ISO/IEC 18004:2000).

- QR code is very similar to QR code model 2; its QR code format differs only in the addition of the facility for symbols to appear in a mirror image orientation for reflectance reversal (light symbols on dark backgrounds) and the option for specifying alternative character is set to the default.

- The micro QR code format, is a variant of QR code with a reduced number of overhead modules and a restricted range of sizes, which enables small to moderate amount of data to be represented in a small symbol, particularly suited to direct marking on parts and components, and to applications where the space available for the symbol is severely restricted.

Standard evolution:

- AIM ITS 97-001
- ISO/IEC 18004:2000
- ISO/IEC 18004:2006
- ISO/IEC 18004:2006/Cor 1:2009
- ISO/IEC 18004:2015
- ISO/IEC 18004:2024

## Definitions

1. ***encoding region***

region of the symbol not occupied by *function patterns* and available for encoding data and error
correction codewords, and for *version* and *format information*

1. ***function patterns***

overhead component of the symbol (finder, separator, timing patterns and alignment patterns)
required for location of the symbol or identification of its characteristics to assist in decoding

1. ***version***

size of the symbol represented in terms of its position in the sequence of permissible sizes

1. ***format information***

encoded pattern containing information on symbol characteristics essential to enable the remainder of the
*encoding region* to be decoded

1. ***remainder codeword***

pad codeword, placed after the error correction codewords, used to fill empty codeword positions to
complete the symbol

## Versions (QR code model 2)

### QR Code

There are 40 sizes from version 1 to version 40. Version 1 measures 21 x 21 modules and each version increases in steps of 4 modules per side up to version 40 which measures 177 x 177 modules.

### Micro QR Code

There are 4 sizes from version M1 to version M4. Version M1 measures 11 x 11 modules and each version increases in steps of 2 modules per side up to version 4 which measures 17 x 17 modules.

## Error correction levels

- L: Low (7%)
- M: Medium (15%)
- Q: Quartile (25%)
- H: High (30%)

Note: The error correction level H is not available in micro QR code symbols.

## Encodable character sets

1) numeric data (digits 0 - 9)
2) alphanumeric data (digits; upper case letters A - Z; space, $ % * + - . / :)
3) byte data (default: ISO/IEC 8859-1; or other sets as otherwise defined)
4) Kanji characters

### Additions to data modes and encoding

FNC1 mode for GS1 application identifiers - application specific data. ECI mode for extended channel interpretation - using other character sets different from Latin-1. Structured append mode for splitting data across multiple symbols.

## Structure

### Quiet zone

For QR Code its width shall be equal to the width of 4 modules.
For Micro QR Code its width shall be equal to the width of 2 modules.

### Finder patterns

Three identical finder patterns located in the upper left, upper right and lower left corners. Each finder pattern consist of three superimposed concentric squares constructed of 7 x 7 dark modules, 5 x 5 light modules and 3 x 3 dark modules.

### Separator

Is a one-module wide separator around each finder pattern consisting of all light modules.

### Timing patterns

Timing patterns are one-module wide row and column of alternating dark and light modules commencing and ending with a dark module. In QR Code the lines run on row and column 6. In Micro QR Code they run on row and column 0.

### Alignment patterns

They are only present in QR Code version 2 or higher. They consist of superimposed concentric squares of sizes 5 (dark modules), 3 (light modules) and 1 (dark module). They are positioned symmetrically on either size of the diagonal and spaced as evenly distribued as possible.

### Format information

The format information is a 15-bit sequence containing 5 data bits with 10 error correction bits calculated using the (15, 5) Bose-Chaudhuri-Hocquenghem code. The first two data bits contain the error correction level as per the following table:

| Err corr lvl | Binary code |
| ------------ | ----------- |
| L | 01 |
| M | 00 |
| Q | 11 |
| H | 10 |

The third to fifth data bits contain the data mask pattern (000 to 111 bits). Masks are not applied to function modules (finder patterns, separator, timing patterns and alignment patterns) only to data modules.

| Mask pattern QR Code | Mask pattern Micro QR Code | Condition |
| -------------------- | -------------------------- | --------- |
| 000 | | (i + j) mod 2 = 0 |
| 001 | 00 | i mod 2 = 0 |
| 010 | | j mod 3 = 0 |
| 011 | | (i + j) mod 3 = 0 |
| 100 | 01 | ((i div 2) + ( j div 3)) mod 2 = 0 |
| 101 | | (i j) mod 2 + (i j) mod 3 = 0 |
| 110 | 10 | ((i j) mod 2 + (i j) mod 3) mod 2 = 0 |
| 111 | 11 | ((i+j) mod 2 + (i j) mod 3) mod 2 = 0 |

The 15-bit error corrected format information must then be XORed with the mask pattern 1010 1000 0010 010 to ensure a non all-zero data string.

```
format_bits = <err_correction_bits><mask_pattern_bits><bch_bits> ^ format_mask
```

The resulting bit sequence is mapped twice into the QR code, in the corresponding areas reserved in column and row 9. The module (4*V + 9, 8) where V is the version number shall always be a dark module and is not part of the format information.

#### BCH Codes

| n | k | t | Generator polynomial |
| - | - | - | -------------------- |
| 7 | 4 | 1 | 1 011 |
| 15 | 11 | 1 | 10 011 |
| 15 | 7 | 2 | 111 010 001 |
| 15 | 5 | 3 | 10 100 110 111 |

### Version information

Version information is included in QR Code version 7 or higher only. It consists of am 18-bit sequence containing 6 data bits with 12 error correction bits calculated using the (18, 6) Golay code. No version information will result in an all-zero data string since only versions 7 to 40 contain version information. Masking is not applied to version information since this block is only present in QR code versions >= 7 (but global data masking is).

The resulting bis sequence is mapped twice in the QR Code, into the areas reserved for it in the 6 x 3 module block above the timing pattern and immediately to the left or the top right finder pattern separator, and the 3 x 6 module block to the left of the timing pattern and immediately above the lower left finder pattern separator.

# Encoding procedure overview

## Step 1: Data analysis

Identify the different characters to encode. Select the desired error correction level. If no version is specified, select the smallest version that can accomodate the data. This is also when if using mode switching, the mode for each data segment is determined.

## Step 2: Data encoding

Convert the data into a bit stream in accordance with the rules for the selected mode or modes. Split the resulting bit stream into 8-bit codewords. Add pad characters as necessary to fill the number of required data codewords required for the version.

## Step 3: Error correction coding

Divide the codeword sequence into the required number of blocks. Generate the error correction codewords for each block.

## Step 4: Structure final message

Interleave the data and error correction codewords from each block and add reminder bits as necessary. Place modules in matrix (together with finder patterns, separators, timing pattern and possibly alignment patterns).

## Step 5: Data masking

Apply the data masking patterns, evaluate the results and select the pattern which
optimizes the dark/light module balance.

## Step 6: Format and version information

Generate the format information and, where applicable, the version information.

# Encoding procedure

## Data analysis

Analyse the input data and appropiate mode to encode each sequence.

## Modes

The default interpretation for QR Code is ECI (extended channel interpretation) 000003 representing the ISO/IEC 8859-1 character set. A QR Code can contain sequences of data in a combination of any of the modes described here. Special sequences of data are used to signal mode changes.

### Numeric mode

The numeric mode encodes data from the decimal digit set (0 - 9) or \x30 to \x39.

### Alphanumeric mode

The alphanumeric mode encodes data from a set of 45 characters: 10 numeric digits, 26 alphabetic characters (A - Z) and nine symbols (SP, $, %, *, +, -, ., /, :) or \x30 to \x39, \x41 to \x5A and \x20, \x24, \x25, \x2A, \x2B, \x2D to \x2F and \x3A.

### Byte mode

In byte mode, data is encoded at 8 bits per character using Latin-1.

### Kanji mode

The Kanji mode efficiently encodes Kanji characters acoording to the shift JIS system based on JIS X 0208.

### Structured append mode

This mode is used to split data across multiple QR Code symbols. It is not covered by this implementation.

#### Modes indicator table

| Mode              | Binary code |
| ----------------- | ----------- |
| ECI               | 0111        |
| Numeric           | 0001        |
| Alphanumeric      | 0010        |
| Byte              | 0100        |
| Kanji             | 1000        |
| Structured append | 0011        |

(*) The termination (end of message) code is 0000.

#### Character count indicator (number of bits)

|  Version | Numeric mode | Alphanumeric mode | Byte mode |
| -------- | ------------ | ----------------- | --------- |
| M1       | 3            | N/A               | N/A       |
| M2       | 4            | 3                 | N/A       |
| M3       | 5            | 4                 | 4         |
| M4       | 6            | 5                 | 5         |
| 1 to 9   | 10           | 9                 | 8         |
| 10 to 26 | 12           | 11                | 16        |
| 27 to 40 | 14           | 13                | 16        |

# Data masking

For reliable QR code reading, it is preferable for dark and light modules to be arranged in a well-balanced
manner in the symbol. The module pattern 1011101 particularly found in the finder pattern should be
avoided in other areas of the symbol as much as possible. To meet the above conditions, data masking should
be applied following the steps.

a. Do not apply data masking to function patterns.

b. Convert the given module pattern in the encoding region (excluding the format information and the
version information) with multiple matrix patterns successively through the XOR operation. For the XOR
operation, lay the module pattern over each of the data masking matrix patterns in turn and reverse
the modules (from light to dark or vice versa) which correspond to dark modules of the data masking
pattern.

c. Evaluate all the resulting converted patterns by charging penalties for undesirable features on each
conversion result.

d. Select the pattern with the lowest penalty points score.
