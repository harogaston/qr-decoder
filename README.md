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

### Timing patterns

### Alignment patterns

### Format information

### Version information
