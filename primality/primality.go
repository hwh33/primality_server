/* Allows users to test primality of integers.

Only logged-in users should have access to this package.

*/

package primality

import "math"

func IsPrime(n int) bool {
	// return naivePrimalityTest(n)
	return naiveThreadedPrimalityTest(n)
}

func naiveThreadedPrimalityTest(n int) bool {

	// This was chosen arbitrarily, but it wouldn't take too much testing to
	// determine the true optimum number.
	// IMPORTANT: this number should be even.
	const operationsPerThread = 50

	// Test for trivial cases.
	if n < 4 || divides(2, n) {
		return false
	}

	// We check to see if any number from 3 to (n/2) is a factor of n. We do
	// this by giving each thread operationsPerThread numbers to check until we
	// run out.
	numThreads := (((n / 2) - 3) / 2) + 1
	results := make([]chan bool, numThreads)
	currentChannelNumber := 0
	startRange := 3
	endRange := startRange + operationsPerThread
	for endRange < (n / 2) {
		results[currentChannelNumber] = make(chan bool)
		go factorExistsWrapper(
			startRange, endRange, n, results[currentChannelNumber])
		currentChannelNumber++
		startRange := endRange + 2
		endRange := startRange + operationsPerThread
	}
	// Hit the rest of the range.
	results[currentChannelNumber] = make(chan bool)
	go factorExistsWrapper(
		startRange, (n / 2), n, results[currentChannelNumber])

	// Now wait for all of our goroutines to finish and get the result.
	// If any of the results are true, then we know that there exists at least
	// one factor of n. So our result is the converse of the logical OR of all
	// results.
	totalResult := false
	for _, resultChan := range results {
		result := <-resultChan
		totalResult = totalResult || result
	}
	return !totalResult
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

/* Check if any odd numbers from a to b divide n. Assumes a and b are odd. */
func factorExistsInRange(a, b, n int) bool {
	for i := a; i <= b; i += 2 {
		if divides(i, n) {
			return true
		}
	}
	return false
}

/* A wrapper for factorExistsInRange which allows us to store the result in a
 * channel.
 */
func factorExistsWrapper(a, b, n int, channel chan bool) {
	channel <- factorExistsInRange(a, b, n)
}

/* Returns true if x divides y. */
func divides(x, y int) bool {
	return math.Mod(float64(y), float64(x)) == 0
}
