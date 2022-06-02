package httpRequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	log "github.com/projectkeas/sdks-service/logger"
	"go.uber.org/zap"
)

func PostJson(uri string, payload interface{}, headers map[string]string) (int, string, error) {

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonData))

	headers["Content-Type"] = "application/json"
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	start := time.Now()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Logger.Error(err.Error(), zap.Error(err))
		return resp.StatusCode, "", err
	}

	bytesReceived := 0

	defer func() {
		resp.Body.Close()
		log.Logger.Debug(fmt.Sprintf("Sent request to '%s'...", uri), zap.Any("http", map[string]interface{}{
			"statusCode":    resp.StatusCode,
			"url":           uri,
			"method":        "POST",
			"latency":       time.Since(start).Milliseconds(),
			"bytesSent":     len(jsonData),
			"bytesReceived": bytesReceived,
		}))
	}()

	responseBody, err := ioutil.ReadAll(resp.Body)
	bytesReceived = len(responseBody)
	if err != nil {
		log.Logger.Error(err.Error(), zap.Error(err))
		if resp.StatusCode == 400 {
			return resp.StatusCode, string(responseBody), err
		} else {
			return resp.StatusCode, "", err
		}
	}

	if resp.StatusCode == 400 {
		responseData := map[string]interface{}{}
		err := json.Unmarshal(responseBody, &responseData)
		if err == nil {
			// We need to check for a suitable ingestion rejection reason
			val, found := responseData["reason"]
			if found && val.(string) == "ingestion-service-rejected" {
				return fiber.StatusAccepted, "", nil
			}
		}

		return resp.StatusCode, string(responseBody), nil
	}

	return resp.StatusCode, "", nil
}
