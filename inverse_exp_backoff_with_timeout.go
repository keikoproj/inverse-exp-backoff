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
		return nil, errors.New("factor should be between 0 and 1")
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
	leftSec := float64((ieb.timeout - time.Now().Sub(ieb.startedAt)).Nanoseconds())
	if leftSec <= 0 || ieb.nextDelay == 0 {
		return errors.New("no more retries left")
	}
	time.Sleep(time.Duration(ieb.nextDelay))
	leftSec2 := leftSec - ieb.nextDelay
	minNano := float64(ieb.min.Nanoseconds())
	newBackoffTime := ieb.nextDelay * ieb.factor

	if newBackoffTime > minNano && leftSec2 > newBackoffTime {
		ieb.nextDelay = newBackoffTime
	} else if newBackoffTime > minNano && leftSec2 < newBackoffTime {
		ieb.nextDelay = leftSec2 - float64(1*time.Second.Nanoseconds())
	} else if newBackoffTime < minNano && leftSec2 - minNano > 1 {
		ieb.nextDelay = minNano
	} else if newBackoffTime < minNano && leftSec2 - minNano < minNano && leftSec2 > 1 {
		// if the newBackoffTime < minNano and leftSec-minNano < minNano and leftSec> 1 sec we should be able to do a retry one last time before 1sec of timeout.
		ieb.nextDelay = leftSec2 - float64(1*time.Second.Nanoseconds())
		if ieb.nextDelay < 0{
			ieb.nextDelay = 0
		}
	}

return nil
}