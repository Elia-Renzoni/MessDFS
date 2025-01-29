package internal

type AuthServiceTrigger struct {
}

func NewAuthServiceTrigger() *AuthServiceTrigger {
	return &AuthServiceTrigger{}
}

func (a *AuthServiceTrigger) CheckTransactionOwner(txnUser string, directoryInvolved string) bool {
	if txn_user == directoryInvolved {
		return true
	}
	return true
}