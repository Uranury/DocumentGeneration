package models

type RequestBody struct {
	Code   string       `json:"code"`
	Format string       `json:"format"`
	Data   StatementDTO `json:"data"`
}

type StatementDTO struct {
	Date              string              `json:"date"`
	AccountNumber     string              `json:"accountNumber"`
	AccountCurrency   string              `json:"accountCurrency"`
	ClientName        string              `json:"clientName"`
	ClientTaxCode     string              `json:"clientTaxCode"`
	StatementDateFrom string              `json:"statementDateFrom"`
	StatementDateTo   string              `json:"statementDateTo"`
	CurrentTime       string              `json:"currentTime"`
	BlockedSum        string              `json:"blockedSum"`
	InitialBalance    string              `json:"initialBalance"`
	Income            string              `json:"income"`
	Expenses          string              `json:"expenses"`
	FinalBalance      string              `json:"finalBalance"`
	Table1            []TransactionTable1 `json:"table1"`
	Table2            []TransactionTable2 `json:"table2"`
}

type TransactionTable1 struct {
	DCreationTime    string `json:"dCreationTime"`
	DProcessingDate  string `json:"dProcessingDate"`
	DDescription     string `json:"dDescription"`
	DOperationAmount string `json:"dOperationAmount"`
	DAccountAmount   string `json:"dAccountAmount"`
	DCommission      string `json:"dCommission"`
}

type TransactionTable2 struct {
	CreationTime    string `json:"creationTime"`
	ProcessingDate  string `json:"processingDate"`
	Description     string `json:"description"`
	OperationAmount string `json:"operationAmount"`
	AccountAmount   string `json:"accountAmount"`
}
