package tokenbucket_test

import (
	"github.com/ChimeraCoder/tokenbucket"
	"testing"
	"time"
)

// Test that a bucket that is full does not block execution
func Test_BucketBuffering(t *testing.T) {
	// Create a bucket with capacity 3, that adds tokens every 4 seconds
	const RATE = 4 * time.Second
	const CAPACITY = 3
	const ERROR = 500 * time.Millisecond
	b := tokenbucket.NewBucket(CAPACITY, RATE)

	// Allow the bucket enough time to fill to capacity
	time.Sleep(CAPACITY * RATE)

	// Check that we can empty the bucket without wasting any time
	before := time.Now()
	<-b.SpendToken(1)
	<-b.SpendToken(1)
	<-b.SpendToken(1)
	after := time.Now()

	if diff := after.Sub(before); diff > RATE {
		t.Errorf("Waited %d seconds, though this should have been nearly instantaneous", diff)
	}
}

// Test that a bucket that is empty blocks execution for the correct amount of time
func Test_BucketCreation(t *testing.T) {
	// Create a bucket with capacity 3, that adds tokens every 4 seconds
	const RATE = 4 * time.Second
	const CAPACITY = 3
	const ERROR = 500 * time.Millisecond
	const EXPECTED_DURATION = RATE * CAPACITY

	b := tokenbucket.NewBucket(CAPACITY, RATE)

	// Ensure that the bucket is empty
	<-b.SpendToken(1)
	<-b.SpendToken(1)
	<-b.SpendToken(1)
	<-b.SpendToken(1)

	// Spending three times on an empty bucket should take 12 seconds
	// (Take the average across three, due to imprecision/scheduling)
	before := time.Now()
	<-b.SpendToken(1)
	<-b.SpendToken(1)
	<-b.SpendToken(1)
	after := time.Now()

	lower := EXPECTED_DURATION - ERROR
	upper := EXPECTED_DURATION + ERROR
	if diff := after.Sub(before); diff < lower || diff > upper {
		t.Errorf("Waited %s seconds, though really should have waited between %s and %s", diff.String(), lower.String(), upper.String())
	}
}
