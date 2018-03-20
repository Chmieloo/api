package data

import (
	"log"

	"github.com/kcapp/api/models"
)

// GetPlayers returns a map of all players
func GetPlayers() (map[int]*models.Player, error) {
	played, err := GetGamesPlayedPerPlayer()
	if err != nil {
		return nil, err
	}

	rows, err := models.DB.Query(`SELECT p.id, p.name, p.nickname, p.games_played, p.games_won, p.created_at FROM player p`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make(map[int]*models.Player)
	for rows.Next() {
		p := new(models.Player)
		err := rows.Scan(&p.ID, &p.Name, &p.Nickname, &p.GamesPlayed, &p.GamesWon, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		players[p.ID] = p

		if val, ok := played[p.ID]; ok {
			p.GamesPlayed = val.GamesPlayed
			p.GamesWon = val.GamesWon
			p.MatchesPlayed = val.MatchesPlayed
			p.MatchesWon = val.MatchesWon
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}

// GetPlayer returns the player for the given ID
func GetPlayer(id int) (*models.Player, error) {
	p := new(models.Player)
	err := models.DB.QueryRow(`SELECT p.id, p.name, p.nickname, p.created_at FROM player p WHERE p.id = ?`, id).Scan(&p.ID, &p.Name, &p.Nickname, &p.CreatedAt)
	if err != nil {
		return nil, err
	}

	pld, err := GetGamesPlayedPerPlayer()
	if err != nil {
		return nil, err
	}
	played := pld[p.ID]
	if played != nil {
		p.GamesPlayed = played.GamesPlayed
		p.GamesWon = played.GamesWon
		p.MatchesPlayed = played.MatchesPlayed
		p.MatchesWon = played.MatchesWon
	}

	return p, nil
}

// AddPlayer will add a new player to the database
func AddPlayer(player models.Player) error {
	// Prepare statement for inserting data
	stmt, err := models.DB.Prepare("INSERT INTO player (name, nickname) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(player.Name, player.Nickname)
	log.Printf("Created new player %s", player.Name)
	return err
}

// GetPlayerScore will get the score for the given player in the given match
func GetPlayerScore(playerID int, matchID int) (int, error) {
	scores, err := GetPlayersScore(matchID)
	if err != nil {
		return 0, err
	}
	return scores[playerID], nil
}

// GetPlayersScore will get the score for all players in the given match
func GetPlayersScore(matchID int) (map[int]int, error) {
	rows, err := models.DB.Query(`
		SELECT
			p2m.player_id,
			m.starting_score - IFNULL(
				SUM(first_dart * first_dart_multiplier) +
				SUM(second_dart * second_dart_multiplier) +
				SUM(third_dart * third_dart_multiplier), 0)
				* IF(g.game_type_id = 2,  -1, 1)
				AS 'current_score'
		FROM player2match p2m
		LEFT JOIN `+"`match`"+` m ON m.id = p2m.match_id
		LEFT JOIN score s ON s.match_id = p2m.match_id AND s.player_id = p2m.player_id
		LEFT JOIN game g on g.id = m.game_id
		WHERE p2m.match_id = ? AND (s.is_bust IS NULL OR is_bust = 0)
		GROUP BY p2m.player_id`, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scores := make(map[int]int)
	for rows.Next() {
		var playerID int
		var score int
		err := rows.Scan(&playerID, &score)
		if err != nil {
			return nil, err
		}
		scores[playerID] = score
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return scores, nil
}

// GetGamesPlayedPerPlayer will get the number of games and matches played and won for each player
func GetGamesPlayedPerPlayer() (map[int]*models.Player, error) {
	rows, err := models.DB.Query(`
		SELECT
			player_id,
			MAX(games_played) AS 'games_played',
			MAX(games_won) AS 'games_won',
			MAX(matches_played) AS 'matches_played',
			MAX(matches_won) AS 'matches_won'
		FROM (
			SELECT
				p2m.player_id,
				COUNT(DISTINCT g.id) AS 'games_played',
				0 AS 'games_won',
				COUNT(DISTINCT m.id) AS 'matches_played',
				0 AS 'matches_won'
			FROM player2match p2m
			JOIN game g ON g.id = p2m.game_id
			JOIN ` + "`match`" + ` m ON m.id = p2m.match_id
			WHERE m.is_finished = 1
			GROUP BY p2m.player_id
			UNION ALL
			SELECT
				p2m.player_id,
				0 AS 'games_played',
				COUNT(DISTINCT g.id) AS 'games_won',
				0 AS 'matches_played',
				COUNT(DISTINCT m.id) AS 'matches_won'
			FROM game g
			JOIN ` + "`match`" + ` m ON m.game_id = g.id
			JOIN player2match p2m ON p2m.player_id = g.winner_id AND p2m.game_id = g.id
			WHERE m.is_finished = 1
			GROUP BY g.winner_id
		) games
		GROUP BY player_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	played := make(map[int]*models.Player)
	for rows.Next() {
		p := new(models.Player)
		err := rows.Scan(&p.ID, &p.GamesPlayed, &p.GamesWon, &p.MatchesPlayed, &p.MatchesWon)
		if err != nil {
			return nil, err
		}
		played[p.ID] = p
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return played, nil
}

// GetPlayerCheckouts will return a list containing all checkouts done by the given player
func GetPlayerCheckouts(playerID int) (map[int][]*models.Visit, error) {
	rows, err := models.DB.Query(`
		SELECT
			player_id,
			first_dart, first_dart_multiplier,
			second_dart, second_dart_multiplier,
			third_dart, third_dart_multiplier,
			(IFNULL(first_dart, 0) * first_dart_multiplier +
				IFNULL(second_dart, 0) * second_dart_multiplier +
				IFNULL(third_dart, 0) * third_dart_multiplier) AS 'checkout'
		FROM score WHERE id IN (
			SELECT MAX(id) FROM score WHERE match_id IN (
				SELECT id FROM `+"`match`"+` WHERE winner_id = ?) GROUP BY match_id)
		ORDER BY checkout`, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	playerVisits := make(map[int]map[string][]*models.Visit)
	for rows.Next() {
		var checkout int
		v := new(models.Visit)
		v.FirstDart = new(models.Dart)
		v.SecondDart = new(models.Dart)
		v.ThirdDart = new(models.Dart)
		err := rows.Scan(&v.PlayerID,
			&v.FirstDart.Value, &v.FirstDart.Multiplier,
			&v.SecondDart.Value, &v.SecondDart.Multiplier,
			&v.ThirdDart.Value, &v.ThirdDart.Multiplier,
			&checkout)
		if err != nil {
			return nil, err
		}

		s := v.GetVisitString()
		if visits, ok := playerVisits[checkout]; ok {
			if val, ok := visits[s]; ok {
				val = append(val, v)
			} else {
				arr := make([]*models.Visit, 0)
				arr = append(arr, v)
				visits[s] = arr
			}
		} else {
			visits := make(map[string][]*models.Visit)
			arr := make([]*models.Visit, 0)
			arr = append(arr, v)
			visits[s] = arr
			playerVisits[checkout] = visits
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	checkouts := make(map[int][]*models.Visit)
	for i := 2; i < 171; i++ {
		if i == 169 || i == 168 || i == 166 || i == 165 || i == 163 || i == 162 || i == 159 {
			// Skip values which cannot be checkouts
			continue
		}
		if val, ok := playerVisits[i]; ok {
			v := make([]*models.Visit, 0)
			for _, visits := range val {
				v = append(v, visits...)
			}
			checkouts[i] = v
		} else {
			checkouts[i] = nil
		}
	}

	return checkouts, nil
}
