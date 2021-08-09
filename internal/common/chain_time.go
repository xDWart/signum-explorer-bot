package common

import "time"

const GENESIS_BLOCK_TIME = 1407722400

func ChainTimeToTime(chainTime int64) time.Time {
	return time.Unix(GENESIS_BLOCK_TIME+chainTime, 0)
}

func FormatChainTimeToStringDatetimeUTC(chainTime int64) string {
	return ChainTimeToTime(chainTime).UTC().Format("2006-01-02 15:04")
}

func FormatChainTimeToStringTimeUTC(chainTime int64) string {
	return ChainTimeToTime(chainTime).UTC().Format("15:04")
}

func FormatChainTimeToStringDate(chainTime int64) string {
	return ChainTimeToTime(chainTime).UTC().Format("2006-01-02")
}
