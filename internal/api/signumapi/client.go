package signumapi

import (
	"log"
	"os"
	abstract_api_client "signum-explorer-bot/internal/api/abstractclient"
	"signum-explorer-bot/internal/config"
	"sync"
)

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

type FeeType float64

const (
	MIN_FEE      FeeType = 0.00735
	CHEAP_FEE            = 0.0147
	STANDARD_FEE         = 0.02205
	PRIORITY_FEE         = 0.0294
)

type Client struct {
	*abstract_api_client.Client
	localAccountCache      AccountCache
	localTransactionsCache TransactionsCache
	localBlocksCache       BlocksCache
	secretPhrase           string
}

func NewClient() *Client {
	secretPhrase := os.Getenv("SECRET_PHRASE")
	if secretPhrase == "" {
		log.Printf("SECRET_PHRASE does not set")
	}

	return &Client{
		abstract_api_client.NewClient(config.SIGNUM_API.HOSTS, nil),
		AccountCache{sync.RWMutex{}, map[string]*Account{}},
		TransactionsCache{sync.RWMutex{}, map[string]map[TransactionType]map[TransactionSubType]*AccountTransactions{}},
		BlocksCache{sync.RWMutex{}, map[string]*AccountBlocks{}},
		secretPhrase,
	}
}
