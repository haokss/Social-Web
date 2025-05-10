package serializer

type TransactionRank struct {
	TransactionTime string  `json:"transactionTime"`
	Counterparty    string  `json:"counterparty"`
	Product         string  `json:"product"`
	IncomeExpense   string  `json:"incomeExpense"`
	Amount          float64 `json:"amount"`
	PaymentMethod   string  `json:"paymentMethod"`
}
