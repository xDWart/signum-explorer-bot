package users

import (
	"fmt"
	"signum-explorer-bot/internal/cross_checker"
	"sort"
)

func (user *User) ProcessCrossing() string {
	user.state = CROSSING_STATE
	return "💽 Please send me a list of your <b>plot file's names</b> separated by new lines," +
		" commas or spaces to check the intersection of nonces:"
}

func (user *User) checkCrossing(message string) string {
	plotsList := cross_checker.CheckPlotsForCrossing(message)

	answer := "📃 <b>Results of cross checking your plots:</b>"
	for account, nonces := range plotsList {
		if account == cross_checker.INVALID_ACCOUNTS {
			continue
		}

		sort.Slice(nonces.ListOfNonces, func(i, j int) bool {
			return nonces.ListOfNonces[i].SharedNonces > nonces.ListOfNonces[j].SharedNonces
		})

		icon := "✅"
		if nonces.AnyError || nonces.SharedNonces > 0 {
			icon = "❌"
		}
		answer += fmt.Sprintf("\n\n%v <b>%v:</b>  <i>%.3f / %.3f</i> TiB",
			icon, account, nonces.PhysicalCapacity-nonces.SharedCapacity, nonces.PhysicalCapacity)
		for _, nonce := range nonces.ListOfNonces {
			var icon string
			var msg string
			if nonce.Error != nil {
				icon = "✖"
				msg = nonce.Error.Error()
			} else if nonce.SharedNonces > 0 {
				icon = "✖"
				msg = fmt.Sprintf("%v shared nonces!", nonce.SharedNonces)
			} else {
				icon = "✔"
				msg = "OK"
			}
			answer += fmt.Sprintf("\n%v %v - %v", icon, nonce.Filename, msg)
		}
	}

	invalidAccounts := plotsList[cross_checker.INVALID_ACCOUNTS]
	if invalidAccounts != nil {
		answer += "\n\n❌ <b>Invalid AccountID:</b>"
		for _, nonce := range invalidAccounts.ListOfNonces {
			answer += "\n" + nonce.Filename
		}
	}

	return answer
}
