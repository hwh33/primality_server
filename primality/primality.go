/* Allows users to test primality of integers.

Only logged-in users should have access to this package.

*/

package primality

import "math"

func IsPrime(n, numThreads int) bool {
	return naivePrimalityTest(n)
}

func naiveThreadedTest(n, numThreads int) bool {
	// Test for trivial cases.
	if n < 4 || divides(2, n) {
		return false
	}

	// Now we test every odd number from 3 to n/2. We divide the work up
	// between numThreads threads. It will end up slightly uneven in most
	// cases, but the difference in workload is small.
	var numJobs int
	if (n/2)-3 >= 0 {
		numJobs = (((n / 2) - 3) / 2) + 1
	} else {
		numJobs = 1
	}

	// In this case, it's not worth the overhead.
	if numJobs < (numThreads * 2) {
		return naivePrimalityTest(n)
	}

	jobsPerThread := numJobs / numThreads

	// double check this
	a := 3
	b := a + (jobsPerThread * 2) - 2
	for b < (n / 2) {
		go factorExistsInRange(a, b, n)
		a = b + 2
		b = a + (jobsPerThread * 2) - 2
	}
	// get the rest
	go factorExistsInRange(a, (n / 2), n)

	// make compiler happy for now
	return false
}

func naivePrimalityTest(n int) bool {
	// Test for trivial cases.
	if n < 4 || divides(2, n) {
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

/* Check if any odd numbers from a to b divide n. */
func factorExistsInRange(a, b, n int) bool {
	for i := a; i <= b; i += 2 {
		if divides(i, n) {
			return true
		}
	}
	return false
}

/* Returns true if x divides y. */
func divides(x, y int) bool {
	return math.Mod(float64(y), float64(x)) == 0
}
