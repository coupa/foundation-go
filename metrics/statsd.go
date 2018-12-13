package metrics

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/alexcesaro/statsd.v2"
)

const (
	EVENTS       = "events"
	TRANSACTIONS = "transactions"
)

//- Package level client and factory

var client *Statsd
var clientFactory func() *Statsd

type Option = statsd.Option
type TagFormat = statsd.TagFormat

type ClientFactory interface {
	NewClient() *Statsd
}

func Get() *Statsd {
	if client == nil && clientFactory != nil {
		client = clientFactory()
	}
	return client
}

func Set(c *Statsd) {
	client = c
}

func SetFactory(f func() *Statsd) {
	clientFactory = f
}

//- Statsd wrapper to provide microservice standard support

type Statsd struct {
	Client *statsd.Client
}

func NewStatsd(address, prefix, version, appName string, sampleRate float32) *Statsd {
	client, err := statsd.New(
		statsd.Address(address),
		statsd.Prefix(prefix),
		statsd.SampleRate(sampleRate),
		statsd.TagsFormat(statsd.InfluxDB),
		statsd.Tags("version", version, "app", appName),
	)
	if err != nil {
		log.Errorf("Error creating a statsd client. Statsd will not be functioning. Error: %v", err)
	}
	return &Statsd{
		Client: client,
	}
}

func (s *Statsd) clone(tags ...map[string]string) *statsd.Client {
	flattened := flatten(tags...)
	return s.Client.Clone(
		statsd.Tags(flattened...),
	)
}

func (s *Statsd) cloneWithSampleRate(rate float32, tags ...map[string]string) *statsd.Client {
	flattened := flatten(tags...)
	return s.Client.Clone(
		statsd.Tags(flattened...),
		statsd.SampleRate(rate),
	)
}

func (s *Statsd) Clone(tags ...map[string]string) *Statsd {
	return &Statsd{
		Client: s.clone(tags...),
	}
}

func (s *Statsd) WithSampleRate(rate float32, tags ...map[string]string) *Statsd {
	return &Statsd{
		Client: s.cloneWithSampleRate(rate, tags...),
	}
}

func (s *Statsd) Count(name string, value interface{}, tags ...map[string]string) {
	tags = addNameTag(name, tags...)
	s.Clone(tags...).Count(EVENTS, value)
}

func (s *Statsd) Increment(name string, tags ...map[string]string) {
	tags = addNameTag(name, tags...)
	s.clone(tags...).Increment(EVENTS)
}

// Gauge records an absolute value for the given bucket.
func (s *Statsd) Gauge(name string, value interface{}, tags ...map[string]string) {
	tags = addNameTag(name, tags...)
	s.Clone(tags...).Gauge(EVENTS, value)
}

func (s *Statsd) Timing(name string, value interface{}, tags ...map[string]string) {
	tags = addNameTag(name, tags...)
	s.clone(tags...).Timing(TRANSACTIONS, value)
}

//The caller will need to call Send() to omit the measurement
func (s *Statsd) NewTiming(name string, tags ...map[string]string) *StatsdTiming {
	tags = addNameTag(name, tags...)
	return &StatsdTiming{
		Timing: s.clone(tags...).NewTiming(),
	}
}

func (s *Statsd) MeasureTiming(name string, functionToMeasure func(), tags ...map[string]string) {
	defer s.NewTiming(name, tags...).Send()
	functionToMeasure()
}

// Histogram sends an histogram value to a bucket.
func (s *Statsd) Histogram(name string, value interface{}, tags ...map[string]string) {
	tags = addNameTag(name, tags...)
	s.Clone(tags...).Histogram(EVENTS, value)
}

// Unique sends the given value to a set bucket.
func (s *Statsd) Unique(name string, value interface{}, tags ...map[string]string) {
	tags = addNameTag(name, tags...)
	s.Clone(tags...).Unique(name, value)
}

// Flush flushes the Client's buffer.
func (s *Statsd) Flush() {
	s.Client.Flush()
}

// Close flushes the Client's buffer and releases the associated ressources. The
// Client and all the cloned Clients must not be used afterward.
func (s *Statsd) Close() {
	s.Client.Close()
}

type StatsdTiming struct {
	Timing statsd.Timing
}

func (t *StatsdTiming) Send() {
	t.Timing.Send(TRANSACTIONS)
}

func addNameTag(name string, tags ...map[string]string) []map[string]string {
	nameTag := map[string]string{"name": name}
	return append(tags, nameTag)
}

func flatten(tags ...map[string]string) []string {
	flattened := []string{}
	for _, ts := range tags {
		for k, v := range ts {
			flattened = append(flattened, k, v)
		}
	}
	return flattened
}
