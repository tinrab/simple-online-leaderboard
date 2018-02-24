package app

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type PostScoreQuery struct {
	Name     string `form:"name"`
	Password string `form:"password"`
	Score    int64  `form:"score"`
}

type GetScoresQuery struct {
	Skip int `form:"skip"`
	Take int `form:"take"`
}

type Player struct {
	Name     string
	Score    int64
	Password []byte
}

type PlayerResult struct {
	Name  string `json:"name"`
	Score string `json:"score"`
}

func postScoreEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	var query PostScoreQuery
	if c.ShouldBindQuery(&query) == nil {
		// Check parameters
		if len(query.Name) < 3 || len(query.Name) > 24 {
			errorResponse(c, http.StatusBadRequest, "Invalid name length")
			return
		}
		if len(query.Password) < 5 || len(query.Password) > 24 {
			errorResponse(c, http.StatusBadRequest, "Invalid password length")
			return
		}

		player := &Player{Name: query.Name}

		if err := readPlayer(ctx, player); err != nil {
			if err == datastore.ErrNoSuchEntity {
				// Create new player
				passwordHash, err := generatePasswordHash(query.Password)
				if err != nil {
					errorResponse(c, http.StatusInternalServerError, "Could not hash password")
					return
				}
				player.Name = query.Name
				player.Score = 0
				player.Password = passwordHash
			} else {
				errorResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			// Check password
			if !comparePassword(player.Password, []byte(query.Password)) {
				errorResponse(c, http.StatusBadRequest, "Incorrect password")
				return
			}
		}

		// Store player score
		if query.Score > player.Score {
			player.Score = query.Score
			if writePlayer(ctx, player) != nil {
				errorResponse(c, http.StatusInternalServerError, "Could not write score")
			}
		}
	} else {
		errorResponse(c, http.StatusBadRequest, "Invalid parameters")
	}
}

func getScoresEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	query := GetScoresQuery{
		Skip: 0,
		Take: 100,
	}
	c.ShouldBindQuery(&query)
	if query.Skip < 0 {
		query.Skip = 0
	}
	if query.Take > 100 {
		query.Take = 100
	}

	if players, err := readLeaderboard(ctx, query.Skip, query.Take); err != nil {
		errorResponse(c, http.StatusInternalServerError, "Could not read scores")
	} else {
		// Convert scores to strings
		result := []PlayerResult{}
		for _, player := range players {
			result = append(result, PlayerResult{
				Name:  player.Name,
				Score: strconv.FormatInt(player.Score, 10),
			})
		}
		c.JSON(http.StatusOK, gin.H{"data": result})
	}
}

func readPlayer(ctx context.Context, player *Player) error {
	key := datastore.NewKey(ctx, "Player", player.Name, 0, nil)
	return datastore.Get(ctx, key, player)
}

func writePlayer(ctx context.Context, player *Player) error {
	key := datastore.NewKey(ctx, "Player", player.Name, 0, nil)
	_, err := datastore.Put(ctx, key, player)
	return err
}

func readLeaderboard(ctx context.Context, skip int, take int) ([]Player, error) {
	var players []Player
	if _, err := datastore.
		NewQuery("Player").
		Order("-Score").
		Offset(skip).
		Limit(take).
		GetAll(ctx, &players); err != nil {
		return nil, err
	}
	return players, nil
}
