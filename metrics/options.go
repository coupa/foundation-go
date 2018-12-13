package metrics

import (
	"gopkg.in/alexcesaro/statsd.v2"
	"time"
)

// Address sets the address of the StatsD daemon.
//
// By default, ":8125" is used. This option is ignored in Client.Clone().
func Address(addr string) Option {
	return statsd.Address(addr)
}

// ErrorHandler sets the function called when an error happens when sending
// metrics (e.g. the StatsD daemon is not listening anymore).
//
// By default, these errors are ignored.  This option is ignored in
// Client.Clone().
func ErrorHandler(h func(error)) Option {
	return statsd.ErrorHandler(h)
}

// FlushPeriod sets how often the Client's buffer is flushed. If p is 0, the
// goroutine that periodically flush the buffer is not lauched and the buffer
// is only flushed when it is full.
//
// By default, the flush period is 100 ms.  This option is ignored in
// Client.Clone().
func FlushPeriod(p time.Duration) Option {
	return statsd.FlushPeriod(p)
}

// MaxPacketSize sets the maximum packet size in bytes sent by the Client.
//
// By default, it is 1440 to avoid IP fragmentation. This option is ignored in
// Client.Clone().
func MaxPacketSize(n int) Option {
	return statsd.MaxPacketSize(n)
}

// Network sets the network (udp, tcp, etc) used by the client. See the
// net.Dial documentation (https://golang.org/pkg/net/#Dial) for the available
// network options.
//
// By default, network is udp. This option is ignored in Client.Clone().
func Network(network string) Option {
	return statsd.Network(network)
}

// Mute sets whether the Client is muted. All methods of a muted Client do
// nothing and return immedialtly.
//
// This option can be used in Client.Clone() only if the parent Client is not
// muted. The clones of a muted Client are always muted.
func Mute(b bool) Option {
	return statsd.Mute(b)
}

// SampleRate sets the sample rate of the Client. It allows sending the metrics
// less often which can be useful for performance intensive code paths.
func SampleRate(rate float32) Option {
	return statsd.SampleRate(rate)
}

// Prefix appends the prefix that will be used in every bucket name.
//
// Note that when used in cloned, the prefix of the parent Client is not
// replaced but is prepended to the given prefix.
func Prefix(p string) Option {
	return statsd.Prefix(p)
}

func Tags(tags ...string) Option {
	return statsd.Tags(tags...)
}

func TagsFormat(tf TagFormat) Option {
	return statsd.TagsFormat(tf)
}
