package models

type Transaction struct {
	ID                string `json:"internalTransactionId"`
	BookingDate       Date   `json:"bookingDate"`
	TransactionAmount struct {
		Amount   FloatString `json:"amount"`
		Currency string      `json:"currency"`
	} `json:"transactionAmount"`
	CreditorName string `json:"creditorName"`
	PurposeCode  string `json:"purposeCode"`
	Description  string `json:"description"`
	BalanceAfter struct {
		BalanceAmount struct {
			Amount   FloatString `json:"amount"`
			Currency string      `json:"currency"`
		} `json:"balanceAmount"`
	} `json:"balanceAfterTransaction"`
}
