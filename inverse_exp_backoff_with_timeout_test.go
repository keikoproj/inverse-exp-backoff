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
	var ietest1, ietest2 *IEBWithTimeout
	var ietest3 *IEBWithTimeout
	var err error
	for ietest1, err = NewIEBWithTimeout(30*time.Second, 5*time.Second, 1*time.Minute, 0.5, time.Now()); err == nil; err = ietest1.Next() {
		appErr := sampleFuncWithTimeout()
		if appErr == nil {
			break
		}
	}
	g.Expect(err).To(gomega.BeNil())

	var err2 error
	for ietest2, err2 = NewIEBWithTimeout(4*time.Second, 1*time.Second, 8*time.Second, 0.4, time.Now()); err2 == nil; err2 = ietest2.Next() {
	}
	g.Expect(err2).NotTo(gomega.BeNil())

	var err3 error
	for ietest3, err3 = NewIEBWithTimeout(5*time.Second, 1*time.Second, 7*time.Second, 0.65, time.Now()); err3 == nil; err3 = ietest3.Next() {
	}
	g.Expect(err3).NotTo(gomega.BeNil())
}
