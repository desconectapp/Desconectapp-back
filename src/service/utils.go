package service

import (
	"errors"
	"fmt"
	repository "gin/db/generated"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5"

	"bytes"
	"encoding/json"
	"net/http"
)

func CreateAccessAndRefreshTokens(userID int32, is_admin bool, accesExpirationTime int64, refreshExpirationTime int64) (string, string, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	body := jwt.MapClaims{
		"sub": userID,
		"exp": accesExpirationTime,
	}
	if is_admin {
		body["is_admin"] = true
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, body)

	body = jwt.MapClaims{
		"sub": userID,
		"exp": refreshExpirationTime,
	}
	if is_admin {
		body["is_admin"] = true
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, body)

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

func ValidateSession(tokenString string) (int32, bool, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return -1, false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp := int64(claims["exp"].(float64))
		if time.Now().Unix() > exp {
			return -2, false, errors.New("invalid token")
		}

		var isAdmin bool = false
		if claims["is_admin"] != nil {
			isAdmin = claims["is_admin"].(bool)
		}

		return int32(claims["sub"].(float64)), isAdmin, nil
	}

	return -3, false, errors.New("invalid token")
}

func getLocationFromCoordinates(lat string, long string) (string, error) {
	apiKey := os.Getenv("LOCATION_IQ_API_KEY")
	url := "https://us1.locationiq.com/v1/reverse?key=" + apiKey + "&lat=" + lat + "&lon=" + long + "&format=json"
	fmt.Printf("url:", url)
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		time.Sleep(2 * time.Second)
		return getLocationFromCoordinates(lat, long)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Print("Response body:", string(body))

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	var addressMap map[string]interface{}

	if a, ok := result["address"].(map[string]interface{}); ok {
		addressMap = a
	}

	fieldsPriority := []string{
		"town",
		"neighbourhood",
		"suburb",
		"state_district",
		"state",
		"country",
	}

	for _, key := range fieldsPriority {
		if val, ok := addressMap[key]; ok {
			if str, ok := val.(string); ok && str != "" {
				return str, nil
			}
		}
	}

	// fallback: display_name or formatted string fields
	if v, ok := result["display_name"].(string); ok && strings.TrimSpace(v) != "" {
		return strings.TrimSpace(v), nil
	}

	return "Ubicacion no disponible", nil
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

	prompt := "Dada una lista de actividades y un nombre de actividad, devolvé en JSON la mejor coincidencia o creá una nueva si no existe o si la coincidencia es demasiado amplia. " +
		"Si la coincidencia es amplia (ej: 'ceramics' → 'arts & crafts'), generá una nueva actividad con nombre corto, emoji y categoría específica. " +
		"Formato del JSON: {\"name\":\"<nombre>\",\"icon\":\"<emoji>\",\"category\":\"<SPORT|CREATIVE|OUTDOOR|INDOOR|GAME|SOCIAL|WELLNESS>\"}. " +
		"Lista de actividades: " + fmt.Sprintf("%v", activitiesNames) + ". " +
		"Nombre de actividad: " + name + "."

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
	req, err := http.NewRequest("POST", "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash-lite:generateContent", bytes.NewBuffer(body))
	if err != nil {
		return repository.Activity{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-goog-api-key", os.Getenv("GEMINI_API_KEY"))

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

	if resp.StatusCode != http.StatusOK {
		
		// 
		icon := "🧉"
		newActivity, err := s.queries.CreateActivity(s.ctx, repository.CreateActivityParams{
			Name:     "Truco",
			Icon:     &icon,
			Category: repository.CategoriesGAME,
		})
		if err != nil {
			return repository.Activity{}, err
		}

		return newActivity, nil



		// return repository.Activity{}, fmt.Errorf("failed to create activity in external API: %s", resp.Status)
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
