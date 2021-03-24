package dispatch

const MessageTypeDispatch = "dispatch"

type Dispatch struct {
	BatchID string `json:"batchId"`
}
