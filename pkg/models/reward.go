package models

type RanklistResponse struct {
	NickName    string `json:"nick_name"`
	Rewards     int    `json:"reward"`
	IsAnonymous bool   `json:"is_anonymous"`
}
