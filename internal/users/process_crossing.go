package users

import (
	"fmt"
	"github.com/xDWart/signum-explorer-bot/internal/crosschecker"
	"sort"
)

func (user *User) ProcessCrossing() string {
	user.state = CROSSING_STATE
	return "💽 Please send me a list of your <b>plot file names</b> separated by new lines, " +
		"commas or spaces to check the crossing of nonces:"
}

func (user *User) checkCrossing(message string) string {
	plotsList := crosschecker.CheckPlotsForCrossing(message)

	var anyError bool
	answer := "💽 <b>Results of cross checking your plots:</b>"
	for account, nonces := range plotsList {
		if account == crosschecker.INVALID_ACCOUNTS {
			continue
		}

		sort.Slice(nonces.ListOfNonces, func(i, j int) bool {
			return nonces.ListOfNonces[i].SharedNonces > nonces.ListOfNonces[j].SharedNonces
		})

		icon := "✅"
		if nonces.AnyError || nonces.SharedNonces > 0 {
			icon = "❌"
			anyError = true
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

	invalidAccounts := plotsList[crosschecker.INVALID_ACCOUNTS]
	if invalidAccounts != nil {
		anyError = true
		answer += "\n\n❌ <b>Invalid AccountID:</b>"
		for _, nonce := range invalidAccounts.ListOfNonces {
			answer += "\n" + nonce.Filename
		}
	}

	if anyError {
		answer += "\n\n🚫 <b>Attention: your plots should not overlap to maximize mining profit, remove duplicates and plot them again!</b>"
	}

	return answer
}
