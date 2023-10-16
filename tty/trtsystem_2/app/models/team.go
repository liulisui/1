package models

type Team struct {
	ID         uint   `json:"team_id"`
	TeamName   string `json:"team_name"`
	TeamLeader string `json:"leader_name"`
	Number     int    `json:"number"`
	Password   []byte `json:"password"`
	Total      int    `json:"total"`
	State      int    `json:"state"`
}
