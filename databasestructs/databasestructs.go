package databasestructs

type PlayerStats struct {
	PlayerID          string  `json:"playerid,omitempty"`
	Games             int64   `json:"g,omitempty"`
	GamesStarted      int64   `json:"gs,omitempty"`
	Minutes           float64 `json:"mpg,omitempty"`
	Points            float64 `json:"ppg,omitempty"`
	Rebounds          float64 `json:"rpg,omitempty"`
	Assists           float64 `json:"apg,omitempty"`
	Steals            float64 `json:"spg,omitempty"`
	Blocks            float64 `json:"bpg,omitempty"`
	Turnovers         float64 `json:"topg,omitempty"`
	FGPercentage      float64 `json:"fgpct,omitempty"`
	ThreeFGPercentage float64 `json:"threefgpct,omitempty"`
	FTPercentage      float64 `json:"ftpct,omitempty"`
	Season            string  `json:"season,omitempty"`
	Position          string  `json:"position,omitempty"`
	TeamAbbr          string  `json:"team,omitempty"`
	IsRookie          bool    `json:"rookie,omitempty"`
}

type AdvancedStats struct {
	PlayerID string  `json:"stats,omitempty"`
	TeamAbbr string  `json:"team,omitempty"`
	Season   string  `json:"season,omitempty"`
	PER      float64 `json:"per,omitempty"`
	TSPct    float64 `json:"ts,omitempty"`
	USGPCt   float64 `json:"usg,omitempty"`
	OffWS    float64 `json:"ows,omitempty"`
	DefWS    float64 `json:"dws,omitempty"`
	WS       float64 `json:"ws,omitempty"`
	OffBPM   float64 `json:"obpm,omitempty"`
	DefBPM   float64 `json:"dbpm,omitempty"`
	BPM      float64 `json:"bpm,omitempty"`
	VORP     float64 `json:"vorp,omitempty"`
	DefRtg   float64 `json:"defrtg,omitempty"`
	OffRtg   float64 `json:"offrtg,omitempty"`
}

type User struct {
	ID           int64  `json:"id,omitempty"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
	ProfilePic   string `json:"profile_pic"`
	IsAdmin      bool   `json:"is_admin"`
}

type Role struct {
	UserID int64
	Role   string
}

type TeamInfo struct {
	TeamAbbr         string
	Name             string
	Logo             string
	WinLossPct       float64
	Playoffs         int64
	DivisionTitles   int64
	ConferenceTitles int64
	Championships    int64
}

type PlayerInfo struct {
	Name          string `json:"name,omitempty"`
	ID            string `json:"playerid,omitempty"`
	College       string `json:"college,omitempty"`
	TeamAbbr      string `json:"team,omitempty"`
	Height        string `json:"height,omitempty"`
	Weight        string `json:"weight,omitempty"`
	Age           int64  `json:"age,omitempty"`
	PlayerStats   `json:"stats,omitempty"`
	AdvancedStats `json:"advstats,omitempty"`
}

type Poll struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Image         string `json:"image"`
	SelectedStats string `json:"selected_stats"`
	Season        string `json:"season"`
	UserID        int64  `json:"user_id,omitempty"`
}

type Image struct {
	ID       int64
	PollID   int64
	ImageURL string
}
