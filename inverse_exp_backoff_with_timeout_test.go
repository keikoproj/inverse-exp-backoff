package iebackoff

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/onsi/gomega"
)

// TestFactorWithTimeout tests that the factor is only allowed to be >0.0 and < 1.0
func TestFactorWithTimeout(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	var factorsToTest = []float64{1.5, 1.0, 0.0}
	//startedAt = time.Now()
	for _, factor := range factorsToTest {
		_, err := NewIEBWithTimeout(4*time.Minute, 30*time.Second, 5*time.Minute, factor, time.Now())
		g.Expect(err).NotTo(gomega.BeNil())
	}

	factor := 0.9
	_, err := NewIEBWithTimeout(4*time.Minute, 30*time.Second, 5*time.Minute, factor, time.Now())
	g.Expect(err).To(gomega.BeNil())
}

func sampleFuncWithTimeout() error {
	fmt.Println(time.Now().String())
	var src = rand.NewSource(time.Now().UnixNano())
	var r = rand.New(src)
	if r.Intn(2) != 0 {
		return errors.New("some random error")
	}

	return nil
}

func TestIEBackoffWithTimeout(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	// The sample test runs a IEBackoff with 5 max retries and a randomized
	// failure rate. sampleFunc is expected to fail 50% of the times. So, within
	// 5 retries, at least one retry should succeed and the error returned should
	// be nil
	var ietest1 *IEBWithTimeout
	var err error
	for ietest1, err = NewIEBWithTimeout(30*time.Second, 1*time.Second, 3*time.Minute, 0.5, time.Now()); err == nil; err = ietest1.Next() {
		appErr := sampleFuncWithTimeout()
		if appErr == nil {
			break
		}
	}
	g.Expect(err).To(gomega.BeNil())

	var err2 error
	// This test goes through all the retries. ietest1.Next() should throw an error at the end.
	for ietest1, err2 = NewIEBWithTimeout(4*time.Second, 1*time.Second, 8*time.Second, 0.4, time.Now()); err2 == nil; err2 = ietest1.Next() {
	}
	g.Expect(err2).NotTo(gomega.BeNil())
}
