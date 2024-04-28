package gorunner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	formatURL  = "https://go.dev/_/fmt?backend="
	compileURL = "https://go.dev/_/compile?backend="
)

func RunCode(code string, imports string) (string, error) {
	formData := url.Values{}
	formData.Set("body", code)
	formData.Set("imports", imports)

	resp, err := http.PostForm(formatURL, formData)
	if err != nil {
		return "", fmt.Errorf("error sending POST request to format URL: %w", err)
	}
	defer resp.Body.Close()

	var formatResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&formatResponse)
	if err != nil {
		return "", fmt.Errorf("error decoding format response: %w", err)
	}
	if formatResponse["Error"] != "" {
		return "", fmt.Errorf("formatting error: %s", formatResponse["Error"])
	} else {
		compiledCode := formatResponse["Body"].(string)
		compileFormData := url.Values{}
		compileFormData.Set("body", compiledCode)

		compileResp, err := http.PostForm(compileURL, compileFormData)
		if err != nil {
			return "", fmt.Errorf("error sending POST request to compile URL: %w", err)
		}
		defer compileResp.Body.Close()

		body, err := io.ReadAll(compileResp.Body)
		if err != nil {
			return "", fmt.Errorf("error reading response body: %w", err)
		}

		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, body, "", "  ")
		if err != nil {
			return "", fmt.Errorf("error formatting JSON: %w", err)
		}

		return prettyJSON.String(), nil
	}
}
