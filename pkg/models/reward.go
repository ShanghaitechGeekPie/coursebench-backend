package models

type RanklistResponse struct {
	NickName    string `json:"nick_name"`
	Reward      int    `json:"reward"`
	IsAnonymous bool   `json:"is_anonymous"`
}
