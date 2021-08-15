package signumapi

type Transaction struct {
	TransactionID string             `json:"transaction"`
	Type          TransactionType    `json:"type"`
	Subtype       TransactionSubType `json:"subtype"`
	Timestamp     int64              `json:"timestamp"`
	RecipientRS   string             `json:"recipientRS"`
	AmountNQT     float64            `json:"amountNQT,string"`
	FeeNQT        float64            `json:"feeNQT,string"`
	Sender        string             `json:"sender"`
	SenderRS      string             `json:"senderRS"`
	Attachment    struct {
		Recipients RecipientsType `json:"recipients"`
		// VersionMultiOutCreation          byte           `json:"version.MultiOutCreation"`
		// VersionCommitmentAdd             byte           `json:"version.CommitmentAdd"`
		// VersionRewardRecipientAssignment byte           `json:"version.RewardRecipientAssignment"`
		// VersionPublicKeyAnnouncement     byte           `json:"version.PublicKeyAnnouncement"`
		// VersionMessage                   byte           `json:"version.Message"`
		// AmountNQT                        uint64         `json:"amountNQT"`
		// Message                          string         `json:"message"`
		// RecipientPublicKey               string         `json:"recipientPublicKey"`
		// MessageIsText                    bool           `json:"messageIsText"`
	} `json:"attachment"`
	// Signature       string             `json:"signature"`
	// SignatureHash   string             `json:"signatureHash"`
	// FullHash        string             `json:"fullHash"`
	// Deadline        uint64             `json:"deadline"`
	// SenderPublicKey string             `json:"senderPublicKey"`
	// Recipient       string             `json:"recipient"`
	// Height         uint64 `json:"height"`
	// Version        uint64 `json:"version"`
	// EcBlockId      uint64 `json:"ecBlockId,string"`
	// EcBlockHeight  uint64 `json:"ecBlockHeight"`
	// Block          uint64 `json:"block,string"`
	// Confirmations  uint64 `json:"confirmations"`
	// BlockTimestamp uint64 `json:"blockTimestamp"`
}

type TransactionRequest struct {
	SecretPhrase string  // is the secret passphrase of the account (optional, but transaction neither signed nor broadcast if omitted)
	FeeNQT       float64 // is the fee (in NQT) for the transaction
	Deadline     int     // deadline (in minutes) for the transaction to be confirmed, 1440 minutes maximum
	// PublicKey    string  // is the public key of the account (optional if secretPhrase provided)
	//Broadcast                    bool    // is set to false to prevent broadcasting the transaction to the network (optional)
	//Message                      string  // is either UTF-8 text or a string of hex digits (perhaps previously encoded using an arbitrary algorithm) to be converted into a bytecode with a maximum length of one kilobyte
	//MessageIsText                bool    // is false if the message is a hex string, otherwise the message is text (optional)
	//MessageToEncrypt             string  // is either UTF-8 text or a string of hex digits to be compressed and converted into a bytecode with a maximum length of one kilobyte, then encrypted using AES (optional)
	//MessageToEncryptIsText       bool    // is false if the message to encrypt is a hex string, otherwise, the message to encrypt is text (optional)
	//EncryptedMessageData         string  // is already encrypted data which overrides messageToEncrypt if provided (optional)
	//EncryptedMessageNonce        string  // is a unique 32-byte number which cannot be reused (optional unless encryptedMessageData is provided)
	//MessageToEncryptToSelf       string  // is either UTF-8 text or a string of hex digits to be compressed and converted into a one-kilobyte maximum bytecode then encrypted with AES, then sent to the sending account (optional)
	//MessageToEncryptToSelfIsText bool    // is false if the message to self-encrypt is a hex string, otherwise the message to encrypt is text (optional)
	//EncryptToSelfMessageData     string  // is already encrypted data which overrides messageToEncryptToSelf if provided (optional)
	//EncryptToSelfMessageNonce    string  // is a unique 32-byte number which cannot be reused (optional unless encryptToSelfMessageData is provided)
	//RecipientPublicKey           string  //is the public key of the receiving account (optional, enhances the security of a new account)
}

type TransactionResponse struct {
	SignatureHash            string      // is the SHA-256 hash of the transaction signature
	UnsignedTransactionBytes string      // is the unsigned transaction bytes
	TransactionJSON          Transaction // is the transaction object (refer to Get Transaction for details)
	Broadcasted              bool        // is the transaction was broadcasted or not
	RequestProcessingTime    int         // is the API request processing time (in millisec)
	TransactionBytes         string      // is the signed transaction bytes
	FullHash                 string      // is the full hash of the signed transaction
	Transaction              string      // is the ID of the newly created transaction
	ErrorDescription         string
}
