package service

import (
	"errors"
	"fmt"
	repository "gin/db/generated"
	"io"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5"

	"bytes"
	"encoding/json"
	"net/http"
)

func CreateAccessAndRefreshTokens(userID int32, accesExpirationTime int64, refreshExpirationTime int64) (string, string, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": accesExpirationTime,
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": refreshExpirationTime,
	})

	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func Expirations() (int64, int64) {
	accesExpirationTime := time.Now().Add(7 * 24 * time.Hour).Unix()
	refreshExpirationTime := time.Now().Add(7 * 24 * time.Hour).Unix()

	return accesExpirationTime, refreshExpirationTime
}

func NewTestToken(userID int32) (string, error) {
	accesExpirationTime := time.Now().Add(7 * 24 * time.Hour).Unix()

	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": accesExpirationTime,
	})

	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

func ValidateSession(tokenString string) (int32, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return -1, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp := int64(claims["exp"].(float64))
			if time.Now().Unix() > exp {
				return -2, errors.New("invalid token")
			}

		return int32(claims["sub"].(float64)), nil
	}

	return -3, errors.New("invalid token")
}

type OpenAIPayload struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type GeminiPayload struct {
	Contents Contents `json:"contents"`
}

type Contents struct {
	Parts Parts `json:"parts"`
}

type Parts struct {
	Text string `json:"text"`
}

func GenerateActivity(s *ActivitiesRequestService, name string) (repository.Activity, error) {

	// Get activities names from DB
	activities, err := s.queries.GetActivities(s.ctx, repository.GetActivitiesParams{Limit: 100, Offset: 0})
	if err != nil {
		return repository.Activity{}, err
	}
	var activitiesNames []string
	for _, activity := range activities {
		activitiesNames = append(activitiesNames, activity.Name)
	}
	
	prompt := "Given the activities listed below and the activity name given, return " +
		"the activity that best matches the meaning of the given name, or generate a new short name for the activity: " +
		"In case the activity is new, also generate an emoji that represents the activity and a category for it. " +
		"The response must be a JSON object of the form: " +
		"{ " +
		"\"name\": \"<name>\", " +
		"\"icon\": \"<emoji>\", " +
		"\"category\": \"<category>\" " +
		"}" + " where category is one of: SPORT, CREATIVE, OUTDOOR, INDOOR, GAME, SOCIAL, WELLNESS. " +
		"Activities: " + fmt.Sprintf("%v", activitiesNames) +
		". Activity name: " + name + "."

	fmt.Print(prompt)

	payload := GeminiPayload{
		Contents: Contents{
			Parts: Parts{
				Text: prompt,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return repository.Activity{}, err
	}

	// Make HTTP request with Authorization header
	req, err := http.NewRequest("POST", "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent", bytes.NewBuffer(body))
	if err != nil {
		return repository.Activity{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-goog-api-key", os.Getenv("GEMINI_API_KEY"))
	fmt.Print("Request body: ", string(body))
	fmt.Print("GEMINI_API_KEY: ", os.Getenv("GEMINI_API_KEY"))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return repository.Activity{}, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return repository.Activity{}, err
	}

	fmt.Print("Response status: ", resp.Status)
	fmt.Print("Response body: ", string(responseBody))

	if resp.StatusCode != http.StatusOK {
		return repository.Activity{}, fmt.Errorf("failed to create activity in external API: %s", resp.Status)
	}
	
	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return repository.Activity{}, err
	}

	// Extract the text from the nested response structure
	candidates, ok := result["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return repository.Activity{}, errors.New("no candidates in response")
	}
	content, ok := candidates[0].(map[string]interface{})["content"].(map[string]interface{})
	if !ok {
		return repository.Activity{}, errors.New("invalid content structure")
	}
	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return repository.Activity{}, errors.New("no parts in content")
	}
	text, ok := parts[0].(map[string]interface{})["text"].(string)
	if !ok {
		return repository.Activity{}, errors.New("invalid text in parts")
	}

	// Remove markdown and parse the JSON
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start == -1 || end == -1 || end <= start {
		return repository.Activity{}, errors.New("invalid JSON format in text")
	}
	jsonStr := text[start : end+1]

	var activityData struct {
		Name     string `json:"name"`
		Icon     string `json:"icon"`
		Category string `json:"category"`
	}
	err = json.Unmarshal([]byte(jsonStr), &activityData)
	if err != nil {
		return repository.Activity{}, err
	}

	// First, check if the activity already exists by name
	existingActivity, err := s.queries.GetActivityByName(s.ctx, activityData.Name)
	if err == nil {
		// Activity found, return it with correct ID
		return existingActivity, nil
	}
	
	// Activity doesn't exist, create it in the database
	if err != pgx.ErrNoRows {
		// This is an unexpected error, not just "not found"
		return repository.Activity{}, err
	}

	// Insert new activity into database
	newActivity, err := s.queries.CreateActivity(s.ctx, repository.CreateActivityParams{
		Name:     activityData.Name,
		Icon:     &activityData.Icon,
		Category: repository.Categories(activityData.Category),
	})
	if err != nil {
		return repository.Activity{}, err
	}

	return newActivity, nil
}