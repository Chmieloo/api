package data

import (
	"github.com/guregu/null"
	"github.com/kcapp/api/models"
)

// GetDartsAtXStatistics will return statistics for all players active during the given period
func GetDartsAtXStatistics(from string, to string, startingScores ...int) ([]*models.StatisticsDartsAtX, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT m.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
			SUM(s.singles) as 'singles',
			SUM(s.doubles) as 'doubles',
			SUM(s.triples) as 'triples',
			SUM(s.singles + s.doubles + s.triples) / (99 * COUNT(DISTINCT l.id)) as 'hit_rate'
		FROM statistics_darts_at_x s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE m.updated_at >= ? AND m.updated_at < ?
			AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_practice = 0
			AND m.match_type_id = 5
		GROUP BY p.id
		ORDER BY(COUNT(DISTINCT m2.id) / COUNT(DISTINCT m.id)) DESC, matches_played DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsDartsAtX, 0)
	for rows.Next() {
		s := new(models.StatisticsDartsAtX)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.AvgScore, &s.Singles, &s.Doubles, &s.Triples, &s.HitRate)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetDartsAtXStatisticsForLeg will return statistics for all players in the given leg
func GetDartsAtXStatisticsForLeg(id int) ([]*models.StatisticsDartsAtX, error) {
	rows, err := models.DB.Query(`
		SELECT
			l.id,
			p.id,
			s.score,
			s.singles,
			s.doubles,
			s.triples,
			s.hit_rate
		FROM statistics_darts_at_x s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
		WHERE l.id = ? GROUP BY p.id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsDartsAtX, 0)
	for rows.Next() {
		s := new(models.StatisticsDartsAtX)
		err := rows.Scan(&s.LegID, &s.PlayerID, &s.Score, &s.Singles, &s.Doubles, &s.Triples, &s.HitRate)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetDartsAtXStatisticsForMatch will return statistics for all players in the given match
func GetDartsAtXStatisticsForMatch(id int) ([]*models.StatisticsDartsAtX, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
			SUM(s.singles) as 'singles',
			SUM(s.doubles) as 'doubles',
			SUM(s.triples) as 'triples',
			SUM(s.singles + s.doubles + s.triples) / 99 * COUNT(DISTINCT l.id) as 'hit_rate'
		FROM statistics_darts_at_x s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			JOIN player2leg p2l ON p2l.leg_id = l.id AND p2l.player_id = s.player_id
		WHERE m.id = ?
		GROUP BY p.id
		ORDER BY p2l.order`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsDartsAtX, 0)
	for rows.Next() {
		s := new(models.StatisticsDartsAtX)
		err := rows.Scan(&s.PlayerID, &s.AvgScore, &s.Singles, &s.Doubles, &s.Triples, &s.HitRate)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

// GetPlayerDartsAtXStatistics will return statistics for the given player
func GetPlayerDartsAtXStatistics(id int) (*models.StatisticsDartsAtX, error) {
	rows, err := models.DB.Query(`
		SELECT
			p.id AS 'player_id',
			COUNT(DISTINCT m.id) AS 'matches_played',
			COUNT(DISTINCT m2.id) AS 'matches_won',
			COUNT(DISTINCT m.id) AS 'legs_played',
			COUNT(DISTINCT l2.id) AS 'legs_won',
			CAST(SUM(s.score) / COUNT(DISTINCT l.id) AS SIGNED) as 'avg_score',
			SUM(s.singles) as 'singles',
			SUM(s.doubles) as 'doubles',
			SUM(s.triples) as 'triples',
			SUM(s.singles + s.doubles + s.triples) / (99 * COUNT(DISTINCT l.id)) as 'hit_rate'
		FROM statistics_darts_at_x s
			JOIN player p ON p.id = s.player_id
			JOIN leg l ON l.id = s.leg_id
			JOIN matches m ON m.id = l.match_id
			LEFT JOIN leg l2 ON l2.id = s.leg_id AND l2.winner_id = p.id
			LEFT JOIN matches m2 ON m2.id = l.match_id AND m2.winner_id = p.id
		WHERE s.player_id = ?
			AND l.is_finished = 1 AND m.is_abandoned = 0 AND m.is_practice = 0
			AND m.match_type_id = 5
		GROUP BY s.player_id`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*models.StatisticsDartsAtX, 0)
	for rows.Next() {
		s := new(models.StatisticsDartsAtX)
		err := rows.Scan(&s.PlayerID, &s.MatchesPlayed, &s.MatchesWon, &s.LegsPlayed, &s.LegsWon, &s.AvgScore, &s.Singles, &s.Doubles, &s.Triples, &s.HitRate)
		if err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	if len(stats) > 0 {
		return stats[0], nil
	}
	return new(models.StatisticsDartsAtX), nil
}

// CalculateDartsAtXStatistics will generate cricket statistics for the given leg
func CalculateDartsAtXStatistics(legID int) (map[int]*models.StatisticsDartsAtX, error) {
	visits, err := GetLegVisits(legID)
	if err != nil {
		return nil, err
	}

	leg, err := GetLeg(legID)
	if err != nil {
		return nil, err
	}

	players, err := GetPlayersScore(legID)
	if err != nil {
		return nil, err
	}
	statisticsMap := make(map[int]*models.StatisticsDartsAtX)
	playerHitsMap := make(map[int]map[int]int64)
	for _, player := range players {
		stats := new(models.StatisticsDartsAtX)
		stats.PlayerID = player.PlayerID
		stats.Score = null.IntFrom(int64(player.CurrentScore))
		statisticsMap[player.PlayerID] = stats
		playerHitsMap[player.PlayerID] = make(map[int]int64)
	}

	number := leg.StartingScore
	for i := 0; i < len(visits); i++ {
		visit := visits[i]
		stats := statisticsMap[visit.PlayerID]

		addDart(number, visit.FirstDart, stats)
		addDart(number, visit.SecondDart, stats)
		addDart(number, visit.ThirdDart, stats)
	}
	for _, stat := range statisticsMap {
		stat.HitRate = float32(stat.Singles+stat.Doubles+stat.Triples) / 99
	}
	return statisticsMap, nil
}

func addDart(number int, dart *models.Dart, stats *models.StatisticsDartsAtX) {
	if dart.ValueRaw() == number {
		if dart.IsTriple() {
			stats.Triples++
		} else if dart.IsDouble() {
			stats.Doubles++
		} else {
			stats.Singles++
		}
	}
}
