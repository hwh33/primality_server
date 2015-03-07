/* Allows users to test primality of integers.

Only logged-in users should have access to this package.

*/

package primality

import (
	"fmt"
	"math"
	"strconv" //debugging
)

func IsPrime(n uint64) bool {
	// return naivePrimalityTest(n)
	return naiveThreadedPrimalityTest(n)
}

func naiveThreadedPrimalityTest(n uint64) bool {

	// This was chosen arbitrarily, but it wouldn't take too much testing to
	// determine the true optimum number.
	// IMPORTANT: this number should be even.
	const operationsPerThread = 50

	// Test for trivial cases.
	if n < 4 {
		return true
	} else if divides(2, n) {
		return false
	}

	// We check to see if any number from 3 to (n/2) is a factor of n. We do
	// this by giving each thread operationsPerThread numbers to check until we
	// run out.
	fmt.Println("Beginning primality test for n = " + strconv.FormatUint(n, 10))
	// numThreads := int(math.Ceil(((float64(n) / 2) - 3) / (2 * 50)))
	results := make([]chan bool, 1)
	// fmt.Println(strconv.Itoa(numThreads) + " channels created for results")
	var startRange uint64 = 3
	var endRange uint64 = startRange + operationsPerThread
	var currentChannelNumber uint64 = 0
	results[currentChannelNumber] = make(chan bool)
	for endRange < (n / 2) {
		// fmt.Println("Creating goroutine " + string(currentChannelNumber))
		go factorExistsWrapper(
			startRange, endRange, n, results[currentChannelNumber])
		currentChannelNumber++
		startRange = endRange + 2
		endRange = startRange + (operationsPerThread * 2)
		results = append(results, make(chan bool))
	}
	// Hit the rest of the range.
	fmt.Println("Creating goroutine " + strconv.FormatUint(currentChannelNumber, 10))
	results[currentChannelNumber] = make(chan bool)
	go factorExistsWrapper(
		startRange, (n / 2), n, results[currentChannelNumber])

	// Now wait for all of our goroutines to finish and get the result.
	// If any of the results are true, then we know that there exists at least
	// one factor of n. So our result is the converse of the logical OR of all
	// results.
	totalResult := false
	fmt.Println("Collecting results from " + strconv.Itoa(len(results)) + " goroutines")
	for _, resultChan := range results {
		result := <-resultChan
		totalResult = totalResult || result
		// fmt.Println("Result " + strconv.FormatUint(uint64(i), 10) + " received")
	}
	fmt.Println("Done")
	return !totalResult
}

func naivePrimalityTest(n uint64) bool {
	// Test for trivial cases.
	if n < 4 || divides(2, n) {
		return false
	}

	// Test every odd number from 3 to n / 2.
	var i uint64
	for i = 3; i < (n / 2); i += 2 {
		if divides(i, n) {
			return false
		}
	}

	// At this point, it must be prime.
	return true
}

/* Check if any odd numbers from a to b divide n. Assumes a and b are odd. */
func factorExistsInRange(a, b, n uint64) bool {
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
func factorExistsWrapper(a, b, n uint64, channel chan bool) {
	channel <- factorExistsInRange(a, b, n)
}

/* Returns true if x divides y. */
func divides(x, y uint64) bool {
	return math.Mod(float64(y), float64(x)) == 0
}
