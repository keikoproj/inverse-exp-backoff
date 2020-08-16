package iebackoff

import (
	"errors"
	"time"
)

// IEBackoff is the basic struct used for Inverse Exponential Backoff.
type IEBWithTimeout struct {
	// max and min are the starting and minimum durations for the backoff.
	max, min time.Duration
	// timeout is the maximum time upto which retires will be performed. After timeout it will be failed.
	timeout time.Duration
	// factor is the multiplication factor
	factor float64
	//start time
	startedAt time.Time
	// nextDelay is the delay that would be used in the next iteration.
	nextDelay float64
}

// NewIEBWithTimeout creates and returns a new IEBackoff object
func NewIEBWithTimeout(max time.Duration, min time.Duration, timeout time.Duration, factor float64, startedAt time.Time) (*IEBWithTimeout, error) {
	if factor >= 1.0 || factor <= 0.0 {
		return nil, errors.New("Factor should be between 0 and 1")
	}

	ieb := IEBWithTimeout{
		max:       max,
		min:       min,
		timeout:   timeout,
		factor:    factor,
		startedAt: startedAt,
		nextDelay: float64(max.Nanoseconds()),
	}
	return &ieb, nil
}

// Next is the main method for the inverse exponential backoff. It takes a function pointer
// and the arguments required for that function as parameters. The function passed as argument
// is expected to return a golang error object.
func (ieb *IEBWithTimeout) Next() error {
	// Confirm there are retries left.
	if  time.Now().Sub(ieb.startedAt) > ieb.timeout {
		return errors.New("No more retries left")
	}

	// Actually sleep for the given delay.
	time.Sleep(time.Duration(ieb.nextDelay))

	// Calculate the delay for the next iteration.
	minNano := float64(ieb.min.Nanoseconds())
	newBackoffTime := ieb.nextDelay * ieb.factor
	if newBackoffTime > minNano {
		ieb.nextDelay = newBackoffTime
	} else {
		// if the newBackoffTime > minNano we should be able to do a retry one last time before 1sec of timeout.
		lastRetry := float64((time.Now().Sub(ieb.startedAt) - ieb.timeout).Nanoseconds())
		if lastRetry <= minNano {
			ieb.nextDelay = lastRetry - 1
		}
	}

	return nil
}