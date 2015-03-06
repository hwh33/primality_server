/* Allows users to test primality of integers.

Only logged-in users should have access to this package.

*/

package primality

import "math"

func IsPrime(n int) bool {
	return naivePrimalityTest(n)
}

func naivePrimalityTest(n int) bool {
	// Test if n is even.
	if divides(2, n) {
		return false
	}

	// Test every odd number from 3 to n / 2.
	for i := 3; i < (n / 2); i += 2 {
		if divides(i, n) {
			return false
		}
	}

	// At this point, it must be prime.
	return true
}

/* Returns true if x divides y. */
func divides(x, y int) bool {
	return math.Mod(float64(y), float64(x)) == 0
}
