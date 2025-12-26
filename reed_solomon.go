package main

// Galois Field GF(256) arithmetic for QR Codes
// Primitive polynomial: x^8 + x^4 + x^3 + x^2 + 1 (0x11D)

const primitivePoly = 0x11D
const gfSize = 256

var expTable [gfSize]int
var logTable [gfSize]int

func init() {
	initTables()
}

func initTables() {
	x := 1
	for i := range gfSize - 1 {
		expTable[i] = x
		logTable[x] = i
		x <<= 1
		if x >= gfSize {
			x ^= primitivePoly
		}
	}
	// logTable[0] is undefined, but often set to 0 or handled separately
}

// gfMul multiplies two numbers in GF(256)
func gfMul(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return expTable[(logTable[a]+logTable[b])%(gfSize-1)]
}

// gfDiv divides two numbers in GF(256)
func gfDiv(a, b int) int {
	if b == 0 {
		panic("division by zero")
	}
	if a == 0 {
		return 0
	}
	return expTable[(logTable[a]-logTable[b]+(gfSize-1))%(gfSize-1)]
}

// gfPow calculates a^b in GF(256)
func gfPow(a, b int) int {
	if b == 0 {
		return 1
	}
	if a == 0 {
		return 0
	}
	return expTable[(logTable[a]*b)%(gfSize-1)]
}

// Polynomial arithmetic

type polynomial struct {
	coeffs []int // Coefficients from highest degree to lowest (or vice-versa? usually highest to lowest in RS implementation)
	// Let's stick to: index 0 is x^n, index len-1 is x^0.
	// Wait, standard implementation often uses index 0 as x^0 or x^n.
	// Let's use: coeffs[0] is the coefficient of x^(len-1), ..., coeffs[len-1] is x^0.
}

func newPolynomial(coeffs []int) *polynomial {
	// Remove leading zeros
	start := 0
	for start < len(coeffs)-1 && coeffs[start] == 0 {
		start++
	}
	return &polynomial{coeffs: coeffs[start:]}
}

func (p *polynomial) degree() int {
	return len(p.coeffs) - 1
}

func polyMul(p1, p2 *polynomial) *polynomial {
	if p1.degree() == -1 || p2.degree() == -1 { // Zero polynomial check (if we represented 0 as empty or [0])
		return newPolynomial([]int{0})
	}
	if len(p1.coeffs) == 1 && p1.coeffs[0] == 0 {
		return newPolynomial([]int{0})
	}
	if len(p2.coeffs) == 1 && p2.coeffs[0] == 0 {
		return newPolynomial([]int{0})
	}

	coeffs := make([]int, len(p1.coeffs)+len(p2.coeffs)-1)
	for i := 0; i < len(p1.coeffs); i++ {
		for j := 0; j < len(p2.coeffs); j++ {
			coeffs[i+j] ^= gfMul(p1.coeffs[i], p2.coeffs[j])
		}
	}
	return newPolynomial(coeffs)
}

func polyMod(dividend, divisor *polynomial) *polynomial {
	if len(dividend.coeffs) < len(divisor.coeffs) {
		return dividend
	}

	// Make a copy of dividend coefficients
	remainder := make([]int, len(dividend.coeffs))
	copy(remainder, dividend.coeffs)

	for len(remainder) >= len(divisor.coeffs) {
		// Leading coefficient of remainder
		coeff := remainder[0]
		// Leading coefficient of divisor is usually 1 for RS generator, but let's be generic
		// If divisor is monic (leading coeff 1), we just multiply by coeff.
		// If not, we'd need to divide. For RS generator poly, it is monic.
		// Let's assume monic for now or implement full division.
		// factor = remainder[0] / divisor[0]
		factor := gfDiv(coeff, divisor.coeffs[0])

		// Subtract (XOR) divisor * factor from remainder
		for i := 0; i < len(divisor.coeffs); i++ {
			remainder[i] ^= gfMul(divisor.coeffs[i], factor)
		}

		// Remove leading zeros (at least one should be removed)
		start := 0
		for start < len(remainder) && remainder[start] == 0 {
			start++
		}
		remainder = remainder[start:]
	}
	if len(remainder) == 0 {
		return newPolynomial([]int{0})
	}
	return newPolynomial(remainder)
}

// Reed-Solomon Encoding

func generateGeneratorPolynomial(numECCodewords int) *polynomial {
	p := newPolynomial([]int{1})
	for i := range numECCodewords {
		// (x - 2^i) -> (x + 2^i) in GF(2^8) since add/sub are XOR
		// coeffs: [1, 2^i]
		term := newPolynomial([]int{1, gfPow(2, i)})
		p = polyMul(p, term)
	}
	return p
}

func reedSolomonEncode(data []byte, numECCodewords int) []byte {
	// Convert data to polynomial coefficients
	dataInts := make([]int, len(data))
	for i, b := range data {
		dataInts[i] = int(b)
	}

	// Message polynomial M(x) * x^n
	// We append n zeros to the data
	paddedData := make([]int, len(data)+numECCodewords)
	copy(paddedData, dataInts)
	messagePoly := newPolynomial(paddedData)

	generatorPoly := generateGeneratorPolynomial(numECCodewords)

	remainderPoly := polyMod(messagePoly, generatorPoly)

	// The result is the remainder coefficients.
	// We need to ensure it has length numECCodewords, padding with leading zeros if necessary.
	remainder := remainderPoly.coeffs
	if len(remainder) < numECCodewords {
		padding := make([]int, numECCodewords-len(remainder))
		remainder = append(padding, remainder...)
	}

	// Convert back to bytes
	ecBytes := make([]byte, numECCodewords)
	for i, v := range remainder {
		ecBytes[i] = byte(v)
	}
	return ecBytes
}
