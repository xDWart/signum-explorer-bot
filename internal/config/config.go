package config

import (
	"regexp"
	"time"
)

var ValidAccount = regexp.MustCompile(`^[0-9]{1,}$`)
var ValidAccountRS = regexp.MustCompile(`^(S|BURST)-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{5}$`)

var COMMON = struct {
	MAX_NUM_OF_ACCOUNTS int
}{
	MAX_NUM_OF_ACCOUNTS: 6,
}

var CALCULATOR = struct {
	REINVEST_EVERY_DAYS int
}{
	REINVEST_EVERY_DAYS: 7,
}

var CMC_API = struct {
	ADDRESS                string
	FREE_LIMIT             string
	CACHE_TTL              time.Duration
	LISTENER_PERIOD        time.Duration
	LISTENER_DAYS_QUANTITY uint
}{
	ADDRESS:                "https://pro-api.coinmarketcap.com/v1",
	FREE_LIMIT:             "200",
	CACHE_TTL:              5 * time.Minute,
	LISTENER_PERIOD:        time.Hour,
	LISTENER_DAYS_QUANTITY: 7,
}

var SIGNUM_API = struct {
	HOSTS                     []string
	DEFAULT_AVG_COMMIT        float64
	DEFAULT_BASE_TARGET       float64
	DEFAULT_BLOCK_REWARD      float64
	GET_NETWORK_INFO_TIME     time.Duration
	AVERAGING_DAYS_QUANTITY   uint
	CACHE_TTL                 time.Duration
	NOTIFIER_PERIOD           time.Duration
	NOTIFIER_CHECK_BLOCKS_PER uint
	DEFAULT_DEADLINE          int
}{
	HOSTS: []string{
		"https://europe1.signum.network",
		"https://europe.signum.network",
		"https://europe2.signum.network",
		"https://europe3.signum.network",
		"https://canada.signum.network",
		"https://australia.signum.network",
		"https://brazil.signum.network",
		"https://uk.signum.network",
		"https://wallet.burstcoin.ro",
	},
	DEFAULT_AVG_COMMIT:        2500,
	DEFAULT_BASE_TARGET:       280000,
	DEFAULT_BLOCK_REWARD:      134,
	GET_NETWORK_INFO_TIME:     time.Hour, // per hour
	AVERAGING_DAYS_QUANTITY:   7,         // during 7 days
	CACHE_TTL:                 3 * time.Minute,
	NOTIFIER_PERIOD:           4 * time.Minute,
	NOTIFIER_CHECK_BLOCKS_PER: 3, // 4 min * 3 = per 12 min
	DEFAULT_DEADLINE:          1440,
}
