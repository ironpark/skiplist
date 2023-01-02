package skiplist

// Options holds Skiplist's options
type Options struct {
	maxLevel    int
	probability float64
	useLock     bool
	usePool     bool
}

// Option is a function used to set Options
type Option func(option *Options)

// WithMutex sets Skiplist goroutine-safety
func WithMutex() Option {
	return func(option *Options) {
		option.useLock = true
	}
}

// WithMaxLevel sets max level of Skiplist
func WithMaxLevel(maxLevel int) Option {
	return func(option *Options) {
		option.maxLevel = maxLevel
	}
}

// WithProbability sets probability of Skiplist
func WithProbability(probability float64) Option {
	return func(option *Options) {
		option.probability = probability
	}
}

// WithPool sets probability of Skiplist
func WithPool() Option {
	return func(option *Options) {
		option.usePool = true
	}
}
