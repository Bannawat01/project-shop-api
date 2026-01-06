package exception

type HistoryOfPurchaseRecording struct{}

func (e *HistoryOfPurchaseRecording) Error() string {
	return "Failed to record purchase history"
}
