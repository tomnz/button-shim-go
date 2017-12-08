package buttonshim

import "time"

// Option allows specifying behavior for the driver.
type Option func(*options)

// WithGamma allows overriding the gamma curve for the LED. Must include 256
// level mappings.
func WithGamma(gamma []byte) Option {
	return func(options *options) {
		if len(gamma) != 256 {
			panic("Must pass 256 gamma levels")
		}
		options.gamma = gamma
	}
}

// WithButtonPollInterval allows modifying the button poll interval. A longer
// interval will use fewer system resources, but risks not registering a fast
// button press. Default 50ms.
func WithButtonPollInterval(interval time.Duration) Option {
	return func(options *options) {
		options.buttonPoll = interval
	}
}

type options struct {
	gamma      []byte
	buttonPoll time.Duration
}

var defaultOptions = options{
	gamma:      defaultGamma,
	buttonPoll: 50 * time.Millisecond,
}
