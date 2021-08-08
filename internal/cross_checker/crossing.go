package cross_checker

type AccountID string

type PlotsList map[AccountID]struct {
	ListOfNonces     []NonceType
	PhysicalCapacity float64
	SharedCapacity   float64
}

type NonceType struct {
	Filename       string
	Error          error
	StartNonce     uint64
	AmountOfNonces uint64
	SharedNonces   uint64
}

func CheckPlotsForCrossing(plots string) (*PlotsList, error) {
	//parcePlots
	//	truncate spaces
	//	replace enter to space
	//	split by space
	//	map[accountID]nonces{}
	//
	//for account, nonces := range plots {
	//	print Account:
	//		checkAccount()
	//			for range nonces
	//				print nonce
	//}
	return nil, nil
}
