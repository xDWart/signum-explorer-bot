package abstractapi

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type SortingType byte

const (
	UNSORTED SortingType = iota
	RANGING
	RANDOM
)

type Config struct {
	sync.Mutex
	apiHostsLatencies []time.Duration
	SortingType       SortingType
	Debug             bool
	ApiHosts          []string
	StaticHeaders     map[string]string
}

func (c *Config) getNextHost(prevIndex int) (string, int) {
	switch c.SortingType {
	case UNSORTED:
		nextIndex := prevIndex + 1
		return c.ApiHosts[nextIndex], nextIndex
	case RANDOM:
		nextIndex := rand.Intn(len(c.ApiHosts))
		return c.ApiHosts[nextIndex], nextIndex
	case RANGING:
		c.Lock()
		defer c.Unlock()

		min := 0
		for i, lat := range c.apiHostsLatencies {
			if lat < c.apiHostsLatencies[min] {
				min = i
			}
		}

		return c.ApiHosts[min], min
	default:
		log.Fatal("Unknown API SortingType: %v", c.SortingType)
		return "", 0
	}
}

const averagingFactor = 20

func (c *Config) appendLatencyToHost(lat time.Duration, i int) {
	if c.SortingType != RANGING {
		return // no latencies for other types
	}

	c.Lock()
	defer c.Unlock()

	c.apiHostsLatencies[i] = (c.apiHostsLatencies[i]*(averagingFactor-1) + lat) / averagingFactor
}

func (c *Config) penaltyTheHost(i int) {
	if c.SortingType != RANGING {
		return // no penalties for other types
	}

	c.Lock()
	defer c.Unlock()

	// if the first host is bad and all next are zeros too
	if c.apiHostsLatencies[0] == 0 {
		c.apiHostsLatencies[0] = 200 * time.Millisecond
		return
	}

	var max time.Duration
	for _, lat := range c.apiHostsLatencies {
		if lat > max {
			max = lat
		}
	}

	// make it worse than the worst
	c.apiHostsLatencies[i] = max * 2
}
