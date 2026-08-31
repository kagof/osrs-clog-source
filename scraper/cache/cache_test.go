package cache

import (
	"errors"
	"scraper/log"
	"testing"
)

func TestCache(t *testing.T) {
	log.Writer = t.Output()
	c := NewCache[int, string](2)

	c.Put(0, "asdf")
	assertEq(t, 1, c.Len())
	assertEq(t, 0, c.oldest)
	c.Put(0, "abcd")
	assertEq(t, 1, c.Len())
	assertEq(t, 0, c.oldest)
	g0, g0ok := c.Get(0)
	assertEq(t, true, g0ok)
	assertEq(t, "abcd", g0)

	c.Put(1, "efgh")
	assertEq(t, 2, c.Len())
	_, g0ok = c.Get(0)
	assertEq(t, true, g0ok)
	_, g1ok := c.Get(1)
	assertEq(t, true, g1ok)

	c.Put(2, "ijkl")
	assertEq(t, 2, c.Len())
	_, g0ok = c.Get(0)
	assertEq(t, false, g0ok)
	_, g1ok = c.Get(1)
	assertEq(t, true, g1ok)
	g2, g2ok := c.Get(2)
	assertEq(t, true, g2ok)
	assertEq(t, "ijkl", g2)
	assertEq(t, 1, c.oldest)

	c.Put(3, "mnop")
	assertEq(t, 2, c.Len())
	_, g0ok = c.Get(0)
	assertEq(t, false, g0ok)
	_, g1ok = c.Get(1)
	assertEq(t, false, g1ok)
	_, g2ok = c.Get(2)
	assertEq(t, true, g2ok)
	g3, g3ok := c.Get(3)
	assertEq(t, true, g3ok)
	assertEq(t, "mnop", g3)
	assertEq(t, 0, c.oldest)
}

func TestCacheMemoize(t *testing.T) {
	log.Writer = t.Output()
	c := NewCache[int, int](2)

	errorToReturn := errors.New("kaboom")
	numCalls := 0

	m := c.Memoized(func(k int) (int, error) {
		numCalls++
		if k < 0 {
			return 0, errorToReturn
		}
		return numCalls, nil
	})

	result, err := m(0)
	assertEq(t, nil, err)
	assertEq(t, 1, numCalls)
	assertEq(t, 1, result)

	result, err = m(0)
	assertEq(t, nil, err)
	assertEq(t, 1, numCalls) // numCalls should not increment
	assertEq(t, 1, result)

	result, err = m(1)
	assertEq(t, nil, err)
	assertEq(t, 2, numCalls)
	assertEq(t, 2, result)

	result, err = m(0)
	assertEq(t, nil, err)
	assertEq(t, 2, numCalls) // numCalls should not increment
	assertEq(t, 1, result)   // cached result

	result, err = m(2)
	assertEq(t, nil, err)
	assertEq(t, 3, numCalls)
	assertEq(t, 3, result)

	result, err = m(1)
	assertEq(t, nil, err)
	assertEq(t, 3, numCalls) // numCalls should not increment
	assertEq(t, 2, result)   // cached result

	result, err = m(0) // no longer should be cached
	assertEq(t, nil, err)
	assertEq(t, 4, numCalls)
	assertEq(t, 4, result)

	result, err = m(-1)
	assertEq(t, errorToReturn, err)
	assertEq(t, 5, numCalls)

	result, err = m(-1)
	assertEq(t, errorToReturn, err)
	assertEq(t, 6, numCalls) // not cached

	result, err = m(2)
	assertEq(t, nil, err)
	assertEq(t, 6, numCalls) // numCalls should not increment
	assertEq(t, 3, result)   // cached result

	result, err = m(3)
	assertEq(t, nil, err)
	assertEq(t, 7, numCalls)
	assertEq(t, 7, result)
}

func assertEq(t *testing.T, exp, got any) {
	t.Helper()
	if exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}
