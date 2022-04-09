package signumapi

import "github.com/xDWart/signum-explorer-bot/api/abstractapi"

type ATDetails struct {
	AT          string `json:"at"`
	MachineData string `json:"machineData"`
	BalanceNQT  uint64 `json:"balanceNQT,string"`
	NextBlock   uint64 `json:"nextBlock"`
	//PrevBalanceNQT        string `json:"prevBalanceNQT"`
	//Frozen                bool   `json:"frozen"`
	//Running               bool   `json:"running"`
	//Stopped               bool   `json:"stopped"`
	//Finished              bool   `json:"finished"`
	//Dead                  bool   `json:"dead"`
	//MachineCodeHashId     string `json:"machineCodeHashId"`
	//AtVersion             int    `json:"atVersion"`
	//AtRS                  string `json:"atRS"`
	//Name                  string `json:"name"`
	//Description           string `json:"description"`
	//Creator               string `json:"creator"`
	//CreatorRS             string `json:"creatorRS"`
	//MachineCode           string `json:"machineCode"`
	//MinActivation         string `json:"minActivation"`
	//CreationBlock         int    `json:"creationBlock"`
	//RequestProcessingTime int    `json:"requestProcessingTime"`
	ErrorDescription string `json:"errorDescription"`
}

func (a *ATDetails) GetError() string {
	return a.ErrorDescription
}

func (a *ATDetails) ClearError() {
	a.ErrorDescription = ""
}

func (c *SignumApiClient) GetATDetails(logger abstractapi.LoggerI, at string) (*ATDetails, error) {
	atDetails := &ATDetails{}
	_, err := c.doJsonReq(logger, "GET", "/burst",
		map[string]string{"requestType": string(RT_GET_AT_DETAILS), "at": at},
		nil,
		atDetails)
	return atDetails, err
}
