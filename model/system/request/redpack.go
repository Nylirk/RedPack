package request

type CreateRedPackRequest struct {
	Amount float64 `json:"amount"`
	Total  int     `json:"total"`
}

type GetRedPackRequest struct {
	RedPackID int64 `json:"redPackUID"`
}
