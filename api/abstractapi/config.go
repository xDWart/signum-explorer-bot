package abstractapi

import (
	"sync"
	"time"
)

type Config struct {
	sync.Mutex
	apiHostsLatencies []time.Duration
	Debug             bool
	ApiHosts          []string
	StaticHeaders     map[string]string
}

func (c *Config) getBestHost() (string, int) {
	c.Lock()
	defer c.Unlock()

	min := 0
	for i, lat := range c.apiHostsLatencies {
		if lat < c.apiHostsLatencies[min] {
			min = i
		}
	}

	return c.ApiHosts[min], min
}

const averagingFactor = 20

func (c *Config) appendLatencyToHost(lat time.Duration, i int) {
	c.Lock()
	defer c.Unlock()

	c.apiHostsLatencies[i] = (c.apiHostsLatencies[i]*(averagingFactor-1) + lat) / averagingFactor
}

func (c *Config) penaltyTheHost(i int) {
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
