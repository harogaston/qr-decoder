qr-decoder
====

# Introduction
There are four technically different, but closely related members of the QR code family, which represent an evolutionary sequence.

- QR code model 1 is the original specification for QR code and is described in AIM ITS 97-001[21].

- QR code model 2 is an enhanced form of the symbology with additional features (primarily, the addition of alignment patterns to assist navigation in larger symbols) and is the basis of the first edition of this document (i.e. ISO/IEC 18004:2000).

- QR code [the basis of the second edition of this document (i.e. ISO/IEC 18004:2006)] is very similar to QR code model 2; its QR code format differs only in the addition of the facility for symbols to appear in a mirror image orientation for reflectance reversal (light symbols on dark backgrounds) and the option for specifying alternative character is set to the default.

- The micro QR code format [also specified in the second edition of this document (i.e. ISO/IEC 18004:2006)], is a variant of QR code with a reduced number of overhead modules and a restricted range of sizes, which enables small to moderate amount of data to be represented in a small symbol, particularly suited to direct marking on parts and components, and to applications where the space available for the symbol is severely restricted.

## Versions
### QR Code
There are 40 sizes from version 1 to version 40. Version 1 measures 21 x 21 modules and each version increases in steps of 4 modules per side up to version 40 which measures 177 x 177 modules.

### Micro QR Code
There are 4 sizes from version M1 to version M4. Version M1 measures 11 x 11 modules and each version increases in steps of 2 modules per side up to version 4 which measures 17 x 17 modules.

## Error correction levels
- L: Low (7%)
- M: Medium (15%)
- Q: Quartile (25%)
- H: High (30%)

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

The third to fifth data bits contain the data mask pattern (000 to 111 bits). Masks are not applied to funcion modules (finder patterns, separator, timing patterns and alignment patterns) only to data modules.

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

The resulting bit string is mapped twice in the QR code, in the corresponding areas reserved in column and row 9. The module (4*V + 9, 8) where V is the version number shall always be a dark module and is not part of the format information.

### Version information
