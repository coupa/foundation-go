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

//Get gets the client object created by the client factory.
func Get() *Statsd {
	if client == nil {
		if clientFactory == nil {
			panic("Cannot initialize statsd client without the client factory")
		}
		client = clientFactory()
	}
	return client
}

func SetFactory(f func() *Statsd) {
	clientFactory = f
}

func Set(c *Statsd) {
	client = c
}

//Statsd is a wrapper around a statsd client to provide microservice standard support
//The general usage also hides the need to clone the client for specifying tags.
//Tags are now received as string hashes, to prevent usage mistakes when using statsd.Tag
type Statsd struct {
	Client *statsd.Client
}

//NewStatsd creates a Statsd object with the given parameters.
//address is the address of telegraf agent
//Sample rate sets what percentage of data points would actually be collected.
//E.g., sampleRate = 1.0 means 100% of the stats are collected;
//sampleRate = 0.1 means 10% of the stats are collected.
//0 < sampleRate <= 1.0
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

//Clone clones this Statsd wrapper and its statsd client objects.
func (s *Statsd) Clone(tags ...map[string]string) *Statsd {
	return &Statsd{
		Client: s.clone(tags...),
	}
}

//WithSampleRate returns a clone which will send metrics with the specified sample
//rate and tags.
func (s *Statsd) WithSampleRate(rate float32, tags ...map[string]string) *Statsd {
	flattened := flatten(tags...)
	return &Statsd{
		Client: s.Client.Clone(
			statsd.Tags(flattened...),
			statsd.SampleRate(rate),
		),
	}
}

//Count adds n to the EVENTS measurement with the specific name tag
func (s *Statsd) Count(name string, value interface{}, tags ...map[string]string) {
	tags = addNameTag(name, tags...)
	s.Clone(tags...).Count(EVENTS, value)
}

//Increment increment the EVENTS measurement with the specific name tag. It is equivalent to Count(name, 1).
func (s *Statsd) Increment(name string, tags ...map[string]string) {
	tags = addNameTag(name, tags...)
	s.clone(tags...).Increment(EVENTS)
}

//Gauge records an absolute value for the EVENTS measurement with the specific name tag.
func (s *Statsd) Gauge(name string, value interface{}, tags ...map[string]string) {
	tags = addNameTag(name, tags...)
	s.Clone(tags...).Gauge(EVENTS, value)
}

//Timing sends a timing value to the TRANSACTIONS measurement with the specific name tag.
func (s *Statsd) Timing(name string, value interface{}, tags ...map[string]string) {
	tags = addNameTag(name, tags...)
	s.clone(tags...).Timing(TRANSACTIONS, value)
}

//NewTiming creates a StatsdTiming object with the current time. Afterwards,
//the caller will need to call Send() on the StatsdTiming to emit the measurement
func (s *Statsd) NewTiming(name string, tags ...map[string]string) *StatsdTiming {
	tags = addNameTag(name, tags...)
	return &StatsdTiming{
		Timing: s.clone(tags...).NewTiming(),
	}
}

//MeasureTiming measures and emit metric on the time it takes for functionToMeasure to execute.
func (s *Statsd) MeasureTiming(name string, functionToMeasure func(), tags ...map[string]string) {
	defer s.NewTiming(name, tags...).Send()
	functionToMeasure()
}

//Histogram sends an histogram value to an EVENTS measurement with the specific name tag.
func (s *Statsd) Histogram(name string, value interface{}, tags ...map[string]string) {
	tags = addNameTag(name, tags...)
	s.Clone(tags...).Histogram(EVENTS, value)
}

//Unique sends the given value to a set bucket.
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

//StatsdTiming when created, it holds a timestamp. When Send() function is called,
//the timing metric is emitted.
type StatsdTiming struct {
	Timing statsd.Timing
}

//Send takes the time difference between now and when the StatsdTiming was created
//and emit the metrics.
func (t *StatsdTiming) Send() {
	t.Timing.Send(TRANSACTIONS)
}

func (s *Statsd) clone(tags ...map[string]string) *statsd.Client {
	flattened := flatten(tags...)
	return s.Client.Clone(
		statsd.Tags(flattened...),
	)
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
