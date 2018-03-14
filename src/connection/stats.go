/*

MIT License

Copyright (c) 2017 Peter Bjorklund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package connection

import (
	"bytes"
	"fmt"
)

type RunningStatsDirection struct {
	packets uint

	octets         uint
	octetsOverhead uint
}

type RunningStats struct {
	Received RunningStatsDirection
	Sent     RunningStatsDirection
}

type StatsDirection struct {
	packetsPerSecond               float32
	octetsWithoutOverheadPerSecond float32
	octetsPerSecond                float32
}

type Stats struct {
	received StatsDirection
	sent     StatsDirection
}

func (s *RunningStatsDirection) AddPackets(count uint, octetCount uint) {
	s.packets += count
	s.octets += octetCount
	const ipv4UDPOverhead = 28
	s.octetsOverhead += count * ipv4UDPOverhead
}

func (s *RunningStatsDirection) Reset() {
	s.octets = 0
	s.packets = 0
	s.octetsOverhead = 0
}

func (s *StatsDirection) SetFromRunningStats(rs *RunningStatsDirection, millisecondsDuration uint) {
	s.octetsPerSecond = float32((rs.octets+rs.octetsOverhead)*1000) / float32(millisecondsDuration)
	s.octetsWithoutOverheadPerSecond = float32((rs.octets)*1000) / float32(millisecondsDuration)
	s.packetsPerSecond = float32(rs.packets*1000) / float32(millisecondsDuration)
	rs.Reset()
}

func (s *Stats) SetFromRunningStats(rs *RunningStats, millisecondsDuration uint) {
	s.received.SetFromRunningStats(&rs.Received, millisecondsDuration)
	s.sent.SetFromRunningStats(&rs.Sent, millisecondsDuration)
}

type Rate int

const (
	B  Rate = 1
	KB Rate = 1000 * B
	MB Rate = 1000 * KB
)

func (b Rate) Megabytes() float64 {
	return float64(b) / float64(MB)
}

func (b Rate) Kilobytes() float64 {
	return float64(b) / float64(KB)
}

func (b Rate) RatePerSecond() string {
	switch {
	case b > (KB * 100):
		return fmt.Sprintf("%.1f mbps", b.Megabytes())
	case b > KB:
		return fmt.Sprintf("%.0f kbps", b.Kilobytes())
	default:
		return fmt.Sprintf("%d bps", b)
	}
}

func (s *StatsDirection) String() string {
	megaBitsPerSecond := Rate(s.octetsPerSecond * 8).RatePerSecond()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	efficiency := float32(0.0)
	const EPSILON = 0.00001
	if s.octetsPerSecond > EPSILON {
		efficiency = s.octetsWithoutOverheadPerSecond / s.octetsPerSecond
	}
	buffer.WriteString(fmt.Sprintf("%v %0.1f packets/s (efficiency:%.0f %%)", megaBitsPerSecond, s.packetsPerSecond, (efficiency * 100.0)))
	buffer.WriteString("]")
	return buffer.String()
}

func (s *Stats) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("[stats ")
	buffer.WriteString(fmt.Sprintf("r:%v s:%v", &s.received, &s.sent))
	buffer.WriteString("]")
	return buffer.String()
}
