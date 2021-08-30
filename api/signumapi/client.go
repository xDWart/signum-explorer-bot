package signumapi

import (
	"github.com/xDWart/signum-explorer-bot/api/abstractapi"
	"sync"
	"time"
)

const DEFAULT_DEADLINE = 1440

type RequestType string

const (
	RT_SEND_MONEY               RequestType = "sendMoney"          // recipient + amountNQT
	RT_SEND_MONEY_MULTI         RequestType = "sendMoneyMulti"     // recipients = <numid1>:<amount1>;<numid2>:<amount2>;<numidN>:<amountN>
	RT_SEND_MONEY_MULTI_SAME    RequestType = "sendMoneyMultiSame" // recipients = <numid1>;<numid2>;<numidN> + amountNQT
	RT_SEND_MESSAGE             RequestType = "sendMessage"
	RT_READ_MESSAGE             RequestType = "readMessage"
	RT_SUGGEST_FEE              RequestType = "suggestFee"
	RT_GET_ACCOUNT              RequestType = "getAccount"
	RT_GET_ACCOUNT_TRANSACTIONS RequestType = "getAccountTransactions"
	RT_GET_MINING_INFO          RequestType = "getMiningInfo"
	RT_GET_REWARD_RECIPIENT     RequestType = "getRewardRecipient"
	RT_SET_REWARD_RECIPIENT     RequestType = "setRewardRecipient"
	RT_ADD_COMMITMENT           RequestType = "addCommitment"
	RT_REMOVE_COMMITMENT        RequestType = "removeCommitment"
	RT_SET_ACCOUNT_INFO         RequestType = "setAccountInfo"
)

type SignumApiClient struct {
	*abstractapi.AbstractApiClient
	localAccountCache      AccountCache
	localTransactionsCache TransactionsCache
	localBlocksCache       BlocksCache
	config                 *Config
}

type Config struct {
	SortingType abstractapi.SortingType
	ApiHosts    []string
	CacheTtl    time.Duration
}

func NewSignumApiClient(logger abstractapi.LoggerI, config *Config) *SignumApiClient {
	abstractConfig := abstractapi.Config{
		SortingType: config.SortingType,
		ApiHosts:    config.ApiHosts,
	}
	return &SignumApiClient{
		AbstractApiClient:      abstractapi.NewAbstractApiClient(logger, &abstractConfig),
		localAccountCache:      AccountCache{sync.RWMutex{}, map[string]*Account{}},
		localTransactionsCache: TransactionsCache{sync.RWMutex{}, map[string]map[TransactionType]map[TransactionSubType]*AccountTransactions{}},
		localBlocksCache:       BlocksCache{sync.RWMutex{}, map[string]*AccountBlocks{}},
		config:                 config,
	}
}
