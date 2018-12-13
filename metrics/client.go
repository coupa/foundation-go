package metrics

func Clone(tags ...map[string]string) *Statsd {
	c := Get()
	if c == nil {
		return nil
	}
	return c.Clone(tags...)
}

func Count(name string, n interface{}, tags ...map[string]string) {
	Get().Count(name, n, tags...)
}

// Increment increment the given name. It is equivalent to Count(name, 1).
func Increment(name string, tags ...map[string]string) {
	Get().Increment(name, tags...)
}

// Gauge records an absolute value for the given name.
func Gauge(name string, value interface{}, tags ...map[string]string) {
	Get().Gauge(name, value, tags...)
}

// Timing sends a timing value to a name.
func Timing(name string, value interface{}, tags ...map[string]string) {
	Get().Timing(name, value, tags...)
}

// Histogram sends an histogram value to a name.
func Histogram(name string, value interface{}, tags ...map[string]string) {
	Get().Histogram(name, value, tags...)
}

// Unique sends the given value to a set name.
func Unique(name string, value string, tags ...map[string]string) {
	Get().Unique(name, value, tags...)
}

// Flush flushes the Client's buffer.
func Flush() {
	Get().Flush()
}

// Close flushes the Client's buffer and releases the associated ressources. The
// Client and all the cloned Clients must not be used afterward.
func Close() {
	Get().Close()
}

// WithSampleRate returns a clone with specified sample rate and tags. Then you
// can call all the metric functions.
func WithSampleRate(rate float32, tags ...map[string]string) *Statsd {
	return Get().WithSampleRate(rate, tags...)
}

// NewTiming returns a StatsdTiming and later you can call Send() on it to finish
// measuring and emit metric.
func NewTiming(name string, tags ...map[string]string) *StatsdTiming {
	return Get().NewTiming(name, tags...)
}

//MeasureTiming measures how much time functionToMeasure takes to run.
func MeasureTiming(name string, functionToMeasure func(), tags ...map[string]string) {
	Get().MeasureTiming(name, functionToMeasure, tags...)
}
