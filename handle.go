package function

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/go-redis/redis"
)


type Answers struct {
	SessionId string `json:"sessionId"`
	OptionA bool `json:"optionA"`
	OptionB bool `json:"optionB"`
	OptionC bool `json:"optionC"`
	OptionD bool `json:"optionD"`
	RemainingTime int `json:"remainingTime"`
}

type GameScore struct {
	SessionId string
	Time      time.Time
	Level     string
	LevelScore int
}

type GameTime struct{
	GameTimeId string
	SessionId string
	Level string
	Type string
	Time      time.Time
}

var redisHost = os.Getenv("REDIS_HOST") // This should include the port which is most of the time 6379
var redisPassword = os.Getenv("REDIS_PASSWORD")

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0,
	})

	points := 0
	var answers Answers

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(req.Body).Decode(&answers)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if answers.OptionA == true {
		points =  0
	}
	if answers.OptionB == true {
		points= 5
	}
	if answers.OptionC == true {
		points = 0
	}
	if answers.OptionD == true {
		points = 3 // KubeCon/KnativeCon special bonus!
	}

	points += answers.RemainingTime

	score := GameScore {
		SessionId: answers.SessionId,
		Level: "kubeconeu-question-1",
		LevelScore: points,
		Time: time.Now(),
	}
	scoreJson, err := json.Marshal(score)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err = client.RPush("score-" + answers.SessionId, string(scoreJson)).Err()
	// if there has been an error setting the value
	// handle the error
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	gt := GameTime{
		GameTimeId: "time-" + score.SessionId,
		SessionId:  score.SessionId,
		Level:      score.Level,
		Type:       "end",
		Time:       score.Time,
	}

	gameTimeJson, err := json.Marshal(gt)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = client.RPush(gt.GameTimeId, string(gameTimeJson)).Err()
	// if there has been an error setting the value
	// handle the error
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(res, string(scoreJson))
}
