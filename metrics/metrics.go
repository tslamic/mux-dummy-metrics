package metrics

import (
	"math"
	"time"
)

type holder struct {
	cnt uint64
	avg time.Duration
}

// add calculates cumulative moving average for this holder.
func (h *holder) add(duration time.Duration) time.Duration {
	cnt := float64(h.cnt)
	avg := float64(h.avg)
	dur := float64(duration)

	newAvg := avg + ((dur - avg) / (cnt + 1))
	h.cnt += 1
	h.avg = time.Duration(math.Round(newAvg))

	return h.avg
}

type Metrics interface {
	Put(key string, d time.Duration) time.Duration
	Get(key string) time.Duration
}

type metrics struct {
	endpoints map[string]*holder
}

func NewMetrics() *metrics {
	endpoints := make(map[string]*holder)
	return &metrics{endpoints: endpoints}
}

func (m *metrics) Put(key string, d time.Duration) time.Duration {
	h, ok := m.endpoints[key]
	if ok {
		return h.add(d)
	}
	m.endpoints[key] = &holder{cnt: 1, avg: d}
	return d
}

func (m *metrics) Get(key string) time.Duration {
	h, ok := m.endpoints[key]
	if ok {
		return h.avg
	}
	return time.Duration(-1)
}
