package iebackoff

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/onsi/gomega"
)

// TestFactor tests that the factor is only allowed to be >0.0 and < 1.0
func TestFactor(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	var factorsToTest = []float64{1.5, 1.0, 0.0}
	for _, factor := range factorsToTest {
		_, err := NewIEBackoff(4*time.Minute, 30*time.Second, factor, 10)
		g.Expect(err).NotTo(gomega.BeNil())
	}

	factor := 0.9
	_, err := NewIEBackoff(4*time.Minute, 30*time.Second, factor, 10)
	g.Expect(err).To(gomega.BeNil())
}

func sampleFunc() error {
	fmt.Println(time.Now().String())
	var src = rand.NewSource(time.Now().UnixNano())
	var r = rand.New(src)
	if r.Intn(2) != 0 {
		return errors.New("some random error")
	}

	return nil
}

func TestIEBackoff(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	// The sample test runs a IEBackoff with 5 max retries and a randomized
	// failure rate. sampleFunc is expected to fail 50% of the times. So, within
	// 5 retries, at least one retry should succeed and the error returned should
	// be nil
	var ietest *IEBackoff
	var err error
	for ietest, err = NewIEBackoff(30*time.Second, 5*time.Second, 0.5, 5); err == nil; err = ietest.Next() {
		appErr := sampleFunc()
		if appErr == nil {
			break
		}
	}
	g.Expect(err).To(gomega.BeNil())

	var err2 error
	// This test goes through all the retries. ietest.Next() should throw an error at the end.
	for ietest, err2 = NewIEBackoff(5*time.Second, 1*time.Second, 0.5, 3); err2 == nil; err2 = ietest.Next() {
	}
	g.Expect(err2).NotTo(gomega.BeNil())
}
