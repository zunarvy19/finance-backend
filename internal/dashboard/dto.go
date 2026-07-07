package dashboard

import "github.com/zunarvy19/finance-backend/internal/transaction"

type DashboardResponse struct {
	TotalBalance       int64                             `json:"total_balance"`
	TodayIncome        int64                             `json:"today_income"`
	TodayExpense       int64                             `json:"today_expense"`
	MonthlyIncome      int64                             `json:"monthly_income"`
	MonthlyExpense     int64                             `json:"monthly_expense"`
	LatestTransactions []transaction.TransactionResponse `json:"latest_transactions"`
}
