package utils

import "sync/atomic"

var (
	// counters holds a slice of all active counters
	counters = map[string]*AtomicCounter{}
)

// AtomicCounter is the counting instance
type AtomicCounter struct {
	Token     string `json:"token"`
	Count     int64  `json:"count"`
	Callbacks []func(*AtomicCounter)
}

// NewAtomicCounter creates a new AtomicCounter instance
func NewAtomicCounter(token string, initial int64) *AtomicCounter {
	return &AtomicCounter{
		Token: token,
		Count: initial,
	}
}

// Add is used to increase/decrease the counter
func (ac *AtomicCounter) Add(count int64) (res int64) {
	res = atomic.AddInt64(&ac.Count, count)
	for _, cb := range ac.Callbacks {
		cb(ac)
	}
	return
}

// AddCallback sets up a callback with every atomic update.
func (ac *AtomicCounter) AddCallback(fn func(*AtomicCounter)) bool {
	l := len(ac.Callbacks)
	ac.Callbacks = append(ac.Callbacks, fn)
	return len(ac.Callbacks) > l
}

// Get returns the countr value
func (ac *AtomicCounter) Get() int64 {
	return atomic.LoadInt64(&ac.Count)
}

// AtomicCount creates or increases/decreases a counter
func AtomicCount(token string, count int64) int64 {
	if counter, ok := counters[token]; ok {
		return counter.Add(count)
	}
	counters[token] = NewAtomicCounter(token, count)
	return counters[token].Count
}

// Atomiccounters handles multiple counters at once
func AtomicCounters(tokens []string, count int64) {
	for _, token := range tokens {
		AtomicCount(token, count)
	}
}

// AddAtomicCallback adds a callback function to an atomic counter
func AddAtomicCallback(token string, fn func(*AtomicCounter)) bool {
	if _, ok := counters[token]; !ok {
		counters[token] = NewAtomicCounter(token, 0)
	}
	return counters[token].AddCallback(fn)
}

// GetAtomicCount returns a counters value
func GetAtomicCount(token string) int64 {
	if counter, ok := counters[token]; ok {
		return counter.Get()
	}
	return 0
}

// GetAtomicCounter returns the AutomicCounter instance
func GetAtomicCounter(token string) (atomic *AtomicCounter, ok bool) {
	atomic, ok = counters[token]
	return
}
