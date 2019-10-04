package models

import (
	"github.com/guregu/null"
)

// BestStatistic struct used for storing a value and leg where the statistic was achieved
type BestStatistic struct {
	Value    int `json:"value"`
	LegID    int `json:"leg_id"`
	PlayerID int `json:"player_id"`
}

// BestStatisticFloat struct used for storing a value and leg where the statistic was achieved
type BestStatisticFloat struct {
	Value    float32 `json:"value"`
	LegID    int     `json:"leg_id"`
	PlayerID int     `json:"player_id"`
}

// StatisticsX01 struct used for storing statistics
type StatisticsX01 struct {
	ID                    int                 `json:"id,omitempty"`
	LegID                 int                 `json:"leg_id,omitempty"`
	PlayerID              int                 `json:"player_id,omitempty"`
	WinnerID              int                 `json:"winner_id,omitempty"`
	PPD                   float32             `json:"ppd"`
	PPDScore              int                 `json:"-"`
	FirstNinePPD          float32             `json:"first_nine_ppd"`
	FirstNinePPDScore     int                 `json:"-"`
	ThreeDartAvg          float32             `json:"three_dart_avg"`
	FirstNineThreeDartAvg float32             `json:"first_nine_three_dart_avg"`
	CheckoutPercentage    null.Float          `json:"checkout_percentage"`
	CheckoutAttempts      int                 `json:"checkout_attempts,omitempty"`
	DartsThrown           int                 `json:"darts_thrown,omitempty"`
	TotalVisits           int                 `json:"total_visits,omitempty"`
	Score60sPlus          int                 `json:"scores_60s_plus"`
	Score100sPlus         int                 `json:"scores_100s_plus"`
	Score140sPlus         int                 `json:"scores_140s_plus"`
	Score180s             int                 `json:"scores_180s"`
	Accuracy20            null.Float          `json:"accuracy_20"`
	Accuracy19            null.Float          `json:"accuracy_19"`
	AccuracyOverall       null.Float          `json:"accuracy_overall"`
	AccuracyStatistics    *AccuracyStatistics `json:"accuracy,omitempty"`
	Visits                []*Visit            `json:"visits,omitempty"`
	Hits                  map[int64]*Hits     `json:"hits,omitempty"`
	MatchesPlayed         int                 `json:"matches_played,omitempty"`
	MatchesWon            int                 `json:"matches_won,omitempty"`
	LegsPlayed            int                 `json:"legs_played,omitempty"`
	LegsWon               int                 `json:"legs_won,omitempty"`
	BestPPD               *BestStatisticFloat `json:"best_ppd,omitempty"`
	BestFirstNinePPD      *BestStatisticFloat `json:"best_first_nine_ppd,omitempty"`
	Best301               *BestStatistic      `json:"best_301,omitempty"`
	Best501               *BestStatistic      `json:"best_501,omitempty"`
	Best701               *BestStatistic      `json:"best_701,omitempty"`
	HighestCheckout       *BestStatistic      `json:"highest_checkout,omitempty"`
	StartingScore         null.Int            `json:"-"`
}

// GlobalStatistics struct used for storing global statistics
type GlobalStatistics struct {
	FishNChips int `json:"fish_n_chips"`
}

// CheckoutStatistics stuct used for storing detailed checkout statistics
type CheckoutStatistics struct {
	Checkout         int         `json:"checkout"`
	Count            int         `json:"count"`
	Completed        bool        `json:"completed"`
	Visits           []*Visit    `json:"visits,omitempty"`
	CheckoutAttempts map[int]int `json:"checkout_attempts,omitempty"`
}

// Hits sturct used to store summary of hits for players/legs
type Hits struct {
	Singles int `json:"1,omitempty"`
	Doubles int `json:"2,omitempty"`
	Triples int `json:"3,omitempty"`
}

// OfficeStatistics struct used for storing statistics for a office
type OfficeStatistics struct {
	PlayerID int `json:"player_id,omitempty"`
	LegID    int `json:"leg_id,omitempty"`
	Checkout int `json:"checkout,omitempty"`
}
