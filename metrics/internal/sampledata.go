// Copyright 2023 The go-zenanet Authors
// This file is part of the go-zenanet library.
//
// The go-zenanet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-zenanet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-zenanet library. If not, see <http://www.gnu.org/licenses/>.

package internal

import (
	"bytes"
	"encoding/gob"
	metrics2 "runtime/metrics"
	"time"

	"github.com/zenanet-network/go-zenanet/metrics"
)

// ExampleMetrics returns an ordered registry populated with a sample of metrics.
func ExampleMetrics() metrics.Registry {
	var registry = metrics.NewOrderedRegistry()

	metrics.NewRegisteredCounterFloat64("test/counter", registry).Inc(12345)
	metrics.NewRegisteredCounterFloat64("test/counter_float64", registry).Inc(54321.98)
	metrics.NewRegisteredGauge("test/gauge", registry).Update(23456)
	metrics.NewRegisteredGaugeFloat64("test/gauge_float64", registry).Update(34567.89)
	metrics.NewRegisteredGaugeInfo("test/gauge_info", registry).Update(
		metrics.GaugeInfoValue{
			"version":           "1.10.18-unstable",
			"arch":              "amd64",
			"os":                "linux",
			"commit":            "7caa2d8163ae3132c1c2d6978c76610caee2d949",
			"protocol_versions": "64 65 66",
		})

	{
		s := metrics.NewUniformSample(3)
		s.Update(1)
		s.Update(2)
		s.Update(3)
		//metrics.NewRegisteredHistogram("test/histogram", registry, metrics.NewSampleSnapshot(3, []int64{1, 2, 3}))
		metrics.NewRegisteredHistogram("test/histogram", registry, s)
	}
	registry.Register("test/meter", metrics.NewInactiveMeter())
	{
		timer := metrics.NewRegisteredResettingTimer("test/resetting_timer", registry)
		timer.Update(10 * time.Millisecond)
		timer.Update(11 * time.Millisecond)
		timer.Update(12 * time.Millisecond)
		timer.Update(120 * time.Millisecond)
		timer.Update(13 * time.Millisecond)
		timer.Update(14 * time.Millisecond)
	}
	{
		timer := metrics.NewRegisteredTimer("test/timer", registry)
		timer.Update(20 * time.Millisecond)
		timer.Update(21 * time.Millisecond)
		timer.Update(22 * time.Millisecond)
		timer.Update(120 * time.Millisecond)
		timer.Update(23 * time.Millisecond)
		timer.Update(24 * time.Millisecond)
		timer.Stop()
	}
	registry.Register("test/empty_resetting_timer", metrics.NewResettingTimer().Snapshot())

	{ // go runtime metrics
		var sLatency = "7\xff\x81\x03\x01\x01\x10Float64Histogram\x01\xff\x82\x00\x01\x02\x01\x06Counts\x01\xff\x84\x00\x01\aBuckets\x01\xff\x86\x00\x00\x00\x16\xff\x83\x02\x01\x01\b[]uint64\x01\xff\x84\x00\x01\x06\x00\x00\x17\xff\x85\x02\x01\x01\t[]float64\x01\xff\x86\x00\x01\b\x00\x00\xfe\x06T\xff\x82\x01\xff\xa2\x00\xfe\r\xef\x00\x01\x02\x02\x04\x05\x04\b\x15\x17 B?6.L;$!2) \x1a? \x190aH7FY6#\x190\x1d\x14\x10\x1b\r\t\x04\x03\x01\x01\x00\x03\x02\x00\x03\x05\x05\x02\x02\x06\x04\v\x06\n\x15\x18\x13'&.\x12=H/L&\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\xff\xa3\xfe\xf0\xff\x00\xf8\x95\xd6&\xe8\v.q>\xf8\x95\xd6&\xe8\v.\x81>\xf8\xdfA:\xdc\x11ŉ>\xf8\x95\xd6&\xe8\v.\x91>\xf8:\x8c0\xe2\x8ey\x95>\xf8\xdfA:\xdc\x11ř>\xf8\x84\xf7C֔\x10\x9e>\xf8\x95\xd6&\xe8\v.\xa1>\xf8:\x8c0\xe2\x8ey\xa5>\xf8\xdfA:\xdc\x11ũ>\xf8\x84\xf7C֔\x10\xae>\xf8\x95\xd6&\xe8\v.\xb1>\xf8:\x8c0\xe2\x8ey\xb5>\xf8\xdfA:\xdc\x11Ź>\xf8\x84\xf7C֔\x10\xbe>\xf8\x95\xd6&\xe8\v.\xc1>\xf8:\x8c0\xe2\x8ey\xc5>\xf8\xdfA:\xdc\x11\xc5\xc9>\xf8\x84\xf7C֔\x10\xce>\xf8\x95\xd6&\xe8\v.\xd1>\xf8:\x8c0\xe2\x8ey\xd5>\xf8\xdfA:\xdc\x11\xc5\xd9>\xf8\x84\xf7C֔\x10\xde>\xf8\x95\xd6&\xe8\v.\xe1>\xf8:\x8c0\xe2\x8ey\xe5>\xf8\xdfA:\xdc\x11\xc5\xe9>\xf8\x84\xf7C֔\x10\xee>\xf8\x95\xd6&\xe8\v.\xf1>\xf8:\x8c0\xe2\x8ey\xf5>\xf8\xdfA:\xdc\x11\xc5\xf9>\xf8\x84\xf7C֔\x10\xfe>\xf8\x95\xd6&\xe8\v.\x01?\xf8:\x8c0\xe2\x8ey\x05?\xf8\xdfA:\xdc\x11\xc5\t?\xf8\x84\xf7C֔\x10\x0e?\xf8\x95\xd6&\xe8\v.\x11?\xf8:\x8c0\xe2\x8ey\x15?\xf8\xdfA:\xdc\x11\xc5\x19?\xf8\x84\xf7C֔\x10\x1e?\xf8\x95\xd6&\xe8\v.!?\xf8:\x8c0\xe2\x8ey%?\xf8\xdfA:\xdc\x11\xc5)?\xf8\x84\xf7C֔\x10.?\xf8\x95\xd6&\xe8\v.1?\xf8:\x8c0\xe2\x8ey5?\xf8\xdfA:\xdc\x11\xc59?\xf8\x84\xf7C֔\x10>?\xf8\x95\xd6&\xe8\v.A?\xf8:\x8c0\xe2\x8eyE?\xf8\xdfA:\xdc\x11\xc5I?\xf8\x84\xf7C֔\x10N?\xf8\x95\xd6&\xe8\v.Q?\xf8:\x8c0\xe2\x8eyU?\xf8\xdfA:\xdc\x11\xc5Y?\xf8\x84\xf7C֔\x10^?\xf8\x95\xd6&\xe8\v.a?\xf8:\x8c0\xe2\x8eye?\xf8\xdfA:\xdc\x11\xc5i?\xf8\x84\xf7C֔\x10n?\xf8\x95\xd6&\xe8\v.q?\xf8:\x8c0\xe2\x8eyu?\xf8\xdfA:\xdc\x11\xc5y?\xf8\x84\xf7C֔\x10~?\xf8\x95\xd6&\xe8\v.\x81?\xf8:\x8c0\xe2\x8ey\x85?\xf8\xdfA:\xdc\x11ŉ?\xf8\x84\xf7C֔\x10\x8e?\xf8\x95\xd6&\xe8\v.\x91?\xf8:\x8c0\xe2\x8ey\x95?\xf8\xdfA:\xdc\x11ř?\xf8\x84\xf7C֔\x10\x9e?\xf8\x95\xd6&\xe8\v.\xa1?\xf8:\x8c0\xe2\x8ey\xa5?\xf8\xdfA:\xdc\x11ũ?\xf8\x84\xf7C֔\x10\xae?\xf8\x95\xd6&\xe8\v.\xb1?\xf8:\x8c0\xe2\x8ey\xb5?\xf8\xdfA:\xdc\x11Ź?\xf8\x84\xf7C֔\x10\xbe?\xf8\x95\xd6&\xe8\v.\xc1?\xf8:\x8c0\xe2\x8ey\xc5?\xf8\xdfA:\xdc\x11\xc5\xc9?\xf8\x84\xf7C֔\x10\xce?\xf8\x95\xd6&\xe8\v.\xd1?\xf8:\x8c0\xe2\x8ey\xd5?\xf8\xdfA:\xdc\x11\xc5\xd9?\xf8\x84\xf7C֔\x10\xde?\xf8\x95\xd6&\xe8\v.\xe1?\xf8:\x8c0\xe2\x8ey\xe5?\xf8\xdfA:\xdc\x11\xc5\xe9?\xf8\x84\xf7C֔\x10\xee?\xf8\x95\xd6&\xe8\v.\xf1?\xf8:\x8c0\xe2\x8ey\xf5?\xf8\xdfA:\xdc\x11\xc5\xf9?\xf8\x84\xf7C֔\x10\xfe?\xf8\x95\xd6&\xe8\v.\x01@\xf8:\x8c0\xe2\x8ey\x05@\xf8\xdfA:\xdc\x11\xc5\t@\xf8\x84\xf7C֔\x10\x0e@\xf8\x95\xd6&\xe8\v.\x11@\xf8:\x8c0\xe2\x8ey\x15@\xf8\xdfA:\xdc\x11\xc5\x19@\xf8\x84\xf7C֔\x10\x1e@\xf8\x95\xd6&\xe8\v.!@\xf8:\x8c0\xe2\x8ey%@\xf8\xdfA:\xdc\x11\xc5)@\xf8\x84\xf7C֔\x10.@\xf8\x95\xd6&\xe8\v.1@\xf8:\x8c0\xe2\x8ey5@\xf8\xdfA:\xdc\x11\xc59@\xf8\x84\xf7C֔\x10>@\xf8\x95\xd6&\xe8\v.A@\xf8:\x8c0\xe2\x8eyE@\xf8\xdfA:\xdc\x11\xc5I@\xf8\x84\xf7C֔\x10N@\xf8\x95\xd6&\xe8\v.Q@\xf8:\x8c0\xe2\x8eyU@\xf8\xdfA:\xdc\x11\xc5Y@\xf8\x84\xf7C֔\x10^@\xf8\x95\xd6&\xe8\v.a@\xf8:\x8c0\xe2\x8eye@\xf8\xdfA:\xdc\x11\xc5i@\xf8\x84\xf7C֔\x10n@\xf8\x95\xd6&\xe8\v.q@\xf8:\x8c0\xe2\x8eyu@\xf8\xdfA:\xdc\x11\xc5y@\xf8\x84\xf7C֔\x10~@\xf8\x95\xd6&\xe8\v.\x81@\xf8:\x8c0\xe2\x8ey\x85@\xf8\xdfA:\xdc\x11ŉ@\xf8\x84\xf7C֔\x10\x8e@\xf8\x95\xd6&\xe8\v.\x91@\xf8:\x8c0\xe2\x8ey\x95@\xf8\xdfA:\xdc\x11ř@\xf8\x84\xf7C֔\x10\x9e@\xf8\x95\xd6&\xe8\v.\xa1@\xf8:\x8c0\xe2\x8ey\xa5@\xf8\xdfA:\xdc\x11ũ@\xf8\x84\xf7C֔\x10\xae@\xf8\x95\xd6&\xe8\v.\xb1@\xf8:\x8c0\xe2\x8ey\xb5@\xf8\xdfA:\xdc\x11Ź@\xf8\x84\xf7C֔\x10\xbe@\xf8\x95\xd6&\xe8\v.\xc1@\xf8:\x8c0\xe2\x8ey\xc5@\xf8\xdfA:\xdc\x11\xc5\xc9@\xf8\x84\xf7C֔\x10\xce@\xf8\x95\xd6&\xe8\v.\xd1@\xf8:\x8c0\xe2\x8ey\xd5@\xf8\xdfA:\xdc\x11\xc5\xd9@\xf8\x84\xf7C֔\x10\xde@\xf8\x95\xd6&\xe8\v.\xe1@\xf8:\x8c0\xe2\x8ey\xe5@\xf8\xdfA:\xdc\x11\xc5\xe9@\xf8\x84\xf7C֔\x10\xee@\xf8\x95\xd6&\xe8\v.\xf1@\xf8:\x8c0\xe2\x8ey\xf5@\xf8\xdfA:\xdc\x11\xc5\xf9@\xf8\x84\xf7C֔\x10\xfe@\xf8\x95\xd6&\xe8\v.\x01A\xfe\xf0\x7f\x00"
		var gcPauses = "7\xff\x81\x03\x01\x01\x10Float64Histogram\x01\xff\x82\x00\x01\x02\x01\x06Counts\x01\xff\x84\x00\x01\aBuckets\x01\xff\x86\x00\x00\x00\x16\xff\x83\x02\x01\x01\b[]uint64\x01\xff\x84\x00\x01\x06\x00\x00\x17\xff\x85\x02\x01\x01\t[]float64\x01\xff\x86\x00\x01\b\x00\x00\xfe\x06R\xff\x82\x01\xff\xa2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x01\x01\x00\x01\x00\x00\x00\x01\x01\x01\x01\x01\x01\x01\x00\x02\x00\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\xff\xa3\xfe\xf0\xff\x00\xf8\x95\xd6&\xe8\v.q>\xf8\x95\xd6&\xe8\v.\x81>\xf8\xdfA:\xdc\x11ŉ>\xf8\x95\xd6&\xe8\v.\x91>\xf8:\x8c0\xe2\x8ey\x95>\xf8\xdfA:\xdc\x11ř>\xf8\x84\xf7C֔\x10\x9e>\xf8\x95\xd6&\xe8\v.\xa1>\xf8:\x8c0\xe2\x8ey\xa5>\xf8\xdfA:\xdc\x11ũ>\xf8\x84\xf7C֔\x10\xae>\xf8\x95\xd6&\xe8\v.\xb1>\xf8:\x8c0\xe2\x8ey\xb5>\xf8\xdfA:\xdc\x11Ź>\xf8\x84\xf7C֔\x10\xbe>\xf8\x95\xd6&\xe8\v.\xc1>\xf8:\x8c0\xe2\x8ey\xc5>\xf8\xdfA:\xdc\x11\xc5\xc9>\xf8\x84\xf7C֔\x10\xce>\xf8\x95\xd6&\xe8\v.\xd1>\xf8:\x8c0\xe2\x8ey\xd5>\xf8\xdfA:\xdc\x11\xc5\xd9>\xf8\x84\xf7C֔\x10\xde>\xf8\x95\xd6&\xe8\v.\xe1>\xf8:\x8c0\xe2\x8ey\xe5>\xf8\xdfA:\xdc\x11\xc5\xe9>\xf8\x84\xf7C֔\x10\xee>\xf8\x95\xd6&\xe8\v.\xf1>\xf8:\x8c0\xe2\x8ey\xf5>\xf8\xdfA:\xdc\x11\xc5\xf9>\xf8\x84\xf7C֔\x10\xfe>\xf8\x95\xd6&\xe8\v.\x01?\xf8:\x8c0\xe2\x8ey\x05?\xf8\xdfA:\xdc\x11\xc5\t?\xf8\x84\xf7C֔\x10\x0e?\xf8\x95\xd6&\xe8\v.\x11?\xf8:\x8c0\xe2\x8ey\x15?\xf8\xdfA:\xdc\x11\xc5\x19?\xf8\x84\xf7C֔\x10\x1e?\xf8\x95\xd6&\xe8\v.!?\xf8:\x8c0\xe2\x8ey%?\xf8\xdfA:\xdc\x11\xc5)?\xf8\x84\xf7C֔\x10.?\xf8\x95\xd6&\xe8\v.1?\xf8:\x8c0\xe2\x8ey5?\xf8\xdfA:\xdc\x11\xc59?\xf8\x84\xf7C֔\x10>?\xf8\x95\xd6&\xe8\v.A?\xf8:\x8c0\xe2\x8eyE?\xf8\xdfA:\xdc\x11\xc5I?\xf8\x84\xf7C֔\x10N?\xf8\x95\xd6&\xe8\v.Q?\xf8:\x8c0\xe2\x8eyU?\xf8\xdfA:\xdc\x11\xc5Y?\xf8\x84\xf7C֔\x10^?\xf8\x95\xd6&\xe8\v.a?\xf8:\x8c0\xe2\x8eye?\xf8\xdfA:\xdc\x11\xc5i?\xf8\x84\xf7C֔\x10n?\xf8\x95\xd6&\xe8\v.q?\xf8:\x8c0\xe2\x8eyu?\xf8\xdfA:\xdc\x11\xc5y?\xf8\x84\xf7C֔\x10~?\xf8\x95\xd6&\xe8\v.\x81?\xf8:\x8c0\xe2\x8ey\x85?\xf8\xdfA:\xdc\x11ŉ?\xf8\x84\xf7C֔\x10\x8e?\xf8\x95\xd6&\xe8\v.\x91?\xf8:\x8c0\xe2\x8ey\x95?\xf8\xdfA:\xdc\x11ř?\xf8\x84\xf7C֔\x10\x9e?\xf8\x95\xd6&\xe8\v.\xa1?\xf8:\x8c0\xe2\x8ey\xa5?\xf8\xdfA:\xdc\x11ũ?\xf8\x84\xf7C֔\x10\xae?\xf8\x95\xd6&\xe8\v.\xb1?\xf8:\x8c0\xe2\x8ey\xb5?\xf8\xdfA:\xdc\x11Ź?\xf8\x84\xf7C֔\x10\xbe?\xf8\x95\xd6&\xe8\v.\xc1?\xf8:\x8c0\xe2\x8ey\xc5?\xf8\xdfA:\xdc\x11\xc5\xc9?\xf8\x84\xf7C֔\x10\xce?\xf8\x95\xd6&\xe8\v.\xd1?\xf8:\x8c0\xe2\x8ey\xd5?\xf8\xdfA:\xdc\x11\xc5\xd9?\xf8\x84\xf7C֔\x10\xde?\xf8\x95\xd6&\xe8\v.\xe1?\xf8:\x8c0\xe2\x8ey\xe5?\xf8\xdfA:\xdc\x11\xc5\xe9?\xf8\x84\xf7C֔\x10\xee?\xf8\x95\xd6&\xe8\v.\xf1?\xf8:\x8c0\xe2\x8ey\xf5?\xf8\xdfA:\xdc\x11\xc5\xf9?\xf8\x84\xf7C֔\x10\xfe?\xf8\x95\xd6&\xe8\v.\x01@\xf8:\x8c0\xe2\x8ey\x05@\xf8\xdfA:\xdc\x11\xc5\t@\xf8\x84\xf7C֔\x10\x0e@\xf8\x95\xd6&\xe8\v.\x11@\xf8:\x8c0\xe2\x8ey\x15@\xf8\xdfA:\xdc\x11\xc5\x19@\xf8\x84\xf7C֔\x10\x1e@\xf8\x95\xd6&\xe8\v.!@\xf8:\x8c0\xe2\x8ey%@\xf8\xdfA:\xdc\x11\xc5)@\xf8\x84\xf7C֔\x10.@\xf8\x95\xd6&\xe8\v.1@\xf8:\x8c0\xe2\x8ey5@\xf8\xdfA:\xdc\x11\xc59@\xf8\x84\xf7C֔\x10>@\xf8\x95\xd6&\xe8\v.A@\xf8:\x8c0\xe2\x8eyE@\xf8\xdfA:\xdc\x11\xc5I@\xf8\x84\xf7C֔\x10N@\xf8\x95\xd6&\xe8\v.Q@\xf8:\x8c0\xe2\x8eyU@\xf8\xdfA:\xdc\x11\xc5Y@\xf8\x84\xf7C֔\x10^@\xf8\x95\xd6&\xe8\v.a@\xf8:\x8c0\xe2\x8eye@\xf8\xdfA:\xdc\x11\xc5i@\xf8\x84\xf7C֔\x10n@\xf8\x95\xd6&\xe8\v.q@\xf8:\x8c0\xe2\x8eyu@\xf8\xdfA:\xdc\x11\xc5y@\xf8\x84\xf7C֔\x10~@\xf8\x95\xd6&\xe8\v.\x81@\xf8:\x8c0\xe2\x8ey\x85@\xf8\xdfA:\xdc\x11ŉ@\xf8\x84\xf7C֔\x10\x8e@\xf8\x95\xd6&\xe8\v.\x91@\xf8:\x8c0\xe2\x8ey\x95@\xf8\xdfA:\xdc\x11ř@\xf8\x84\xf7C֔\x10\x9e@\xf8\x95\xd6&\xe8\v.\xa1@\xf8:\x8c0\xe2\x8ey\xa5@\xf8\xdfA:\xdc\x11ũ@\xf8\x84\xf7C֔\x10\xae@\xf8\x95\xd6&\xe8\v.\xb1@\xf8:\x8c0\xe2\x8ey\xb5@\xf8\xdfA:\xdc\x11Ź@\xf8\x84\xf7C֔\x10\xbe@\xf8\x95\xd6&\xe8\v.\xc1@\xf8:\x8c0\xe2\x8ey\xc5@\xf8\xdfA:\xdc\x11\xc5\xc9@\xf8\x84\xf7C֔\x10\xce@\xf8\x95\xd6&\xe8\v.\xd1@\xf8:\x8c0\xe2\x8ey\xd5@\xf8\xdfA:\xdc\x11\xc5\xd9@\xf8\x84\xf7C֔\x10\xde@\xf8\x95\xd6&\xe8\v.\xe1@\xf8:\x8c0\xe2\x8ey\xe5@\xf8\xdfA:\xdc\x11\xc5\xe9@\xf8\x84\xf7C֔\x10\xee@\xf8\x95\xd6&\xe8\v.\xf1@\xf8:\x8c0\xe2\x8ey\xf5@\xf8\xdfA:\xdc\x11\xc5\xf9@\xf8\x84\xf7C֔\x10\xfe@\xf8\x95\xd6&\xe8\v.\x01A\xfe\xf0\x7f\x00"

		var secondsToNs = float64(time.Second)

		dserialize := func(data string) *metrics2.Float64Histogram {
			var res metrics2.Float64Histogram
			if err := gob.NewDecoder(bytes.NewReader([]byte(data))).Decode(&res); err != nil {
				panic(err)
			}
			return &res
		}
		cpuSchedLatency := metrics.RuntimeHistogramFromData(secondsToNs, dserialize(sLatency))
		registry.Register("system/cpu/schedlatency", cpuSchedLatency)

		memPauses := metrics.RuntimeHistogramFromData(secondsToNs, dserialize(gcPauses))
		registry.Register("system/memory/pauses", memPauses)
	}
	return registry
}
