package crosschecker

import (
	"errors"
	"math"
	"signum-explorer-bot/internal/config"
	"strconv"
	"strings"
)

const INVALID_ACCOUNTS = "INVALID_ACCOUNTS"

type NoncesType struct {
	ListOfNonces     []*NonceType
	AnyError         bool
	TotalNonces      uint64
	PhysicalCapacity float64
	SharedNonces     uint64
	SharedCapacity   float64
}

func (nt *NoncesType) add(nonce *NonceType) {
	nt.ListOfNonces = append(nt.ListOfNonces, nonce)
}

type NonceType struct {
	Filename       string
	Error          error
	StartNonce     uint64
	AmountOfNonces uint64
	FinishOfNonces uint64
	SharedNonces   uint64
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

var NONCE_SIZE = math.Pow(2, 18)
var TIB_BYTES = math.Pow(2, 40)

func CheckPlotsForCrossing(plots string) map[string]*NoncesType {
	plotsList := parsePlots(plots)

	for _, nonces := range plotsList { // range accounts
		for index1, nonce1 := range nonces.ListOfNonces { // range each file
			if nonce1.Error != nil {
				continue
			}
			nonces.TotalNonces += nonce1.AmountOfNonces
			for index2, nonce2 := range nonces.ListOfNonces { // compare with each
				if index1 == index2 || nonce2.Error != nil { // except for myself
					continue
				}

				if nonce1.StartNonce < nonce2.FinishOfNonces && nonce1.FinishOfNonces > nonce2.StartNonce { // they are crossing
					diff := min(nonce1.FinishOfNonces, nonce2.FinishOfNonces) - max(nonce1.StartNonce, nonce2.StartNonce)
					nonce1.SharedNonces += diff
					nonces.SharedNonces += diff
				}
			}
		}
		nonces.SharedNonces /= 2
		nonces.PhysicalCapacity = float64(nonces.TotalNonces) * NONCE_SIZE / TIB_BYTES
		nonces.SharedCapacity = float64(nonces.SharedNonces) * NONCE_SIZE / TIB_BYTES
	}

	return plotsList
}

func getOrCreateNoncesType(plotsList map[string]*NoncesType, accountID string) *NoncesType {
	nonces, ok := plotsList[accountID]
	if !ok {
		nonces = &NoncesType{
			ListOfNonces: make([]*NonceType, 0, 1),
		}
	}
	plotsList[accountID] = nonces
	return nonces
}

func parsePlots(plots string) map[string]*NoncesType {
	var plotsList = make(map[string]*NoncesType)

	plots = strings.TrimSpace(plots)
	plots = strings.Replace(plots, "\n", " ", -1)
	plots = strings.Replace(plots, ",", " ", -1)
	plots = strings.Join(strings.Fields(plots), " ")

	splittedPlots := strings.Split(plots, " ")
	for _, plot := range splittedPlots {
		newNonce := &NonceType{
			Filename: plot,
		}

		splittedPlot := strings.Split(plot, "_")
		accountID := splittedPlot[0]
		if !config.ValidAccount.MatchString(accountID) {
			newNonce.Error = errors.New("invalid AccountID")
			nonces := getOrCreateNoncesType(plotsList, INVALID_ACCOUNTS)
			nonces.add(newNonce)
			nonces.AnyError = true
			continue
		}

		nonces := getOrCreateNoncesType(plotsList, accountID)
		nonces.add(newNonce)

		if len(splittedPlot) != 3 {
			newNonce.Error = errors.New("invalid filename")
			nonces.AnyError = true
			continue
		}

		newNonce.StartNonce, newNonce.Error = strconv.ParseUint(splittedPlot[1], 10, 64)
		if newNonce.Error != nil {
			newNonce.Error = errors.New("invalid syntax")
			nonces.AnyError = true
			continue
		}

		newNonce.AmountOfNonces, newNonce.Error = strconv.ParseUint(splittedPlot[2], 10, 64)
		if newNonce.Error != nil {
			newNonce.Error = errors.New("invalid syntax")
			nonces.AnyError = true
			continue
		}

		newNonce.FinishOfNonces = newNonce.StartNonce + newNonce.AmountOfNonces
	}

	return plotsList
}
