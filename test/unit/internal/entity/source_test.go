package entity_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/fyvri/fresh-proxy-list/internal/entity"
)

var (
	expectedButGotMessage      = "Expected %v = %v, but got = %v"
	expectedErrorButGotMessage = "Expected %v error = %v, but got = %v"
	testMethod                 = "LIST"
	testCategory               = "HTTP"
	testURL                    = "http://example.com"
	testIsChecked              = false
)

func TestUnmarshalJSONWithIsChecked(t *testing.T) {
	var (
		source = entity.Source{}
		data   = []byte(`{
			"method": "` + testMethod + `",
			"category": "` + testCategory + `",
			"url": "` + testURL + `",
			"is_checked": ` + strconv.FormatBool(testIsChecked) + `
		}`)
	)
	err := json.Unmarshal(data, &source)

	if err != nil {
		t.Errorf(expectedErrorButGotMessage, "unmarshal", nil, err)
	}

	if source.Method != testMethod {
		t.Errorf(expectedButGotMessage, "method", testMethod, source.Method)
	}

	if source.Category != testCategory {
		t.Errorf(expectedButGotMessage, "category", testCategory, source.Category)
	}

	if source.URL != testURL {
		t.Errorf(expectedButGotMessage, "url", testURL, source.Category)
	}

	if source.IsChecked != testIsChecked {
		t.Errorf(expectedButGotMessage, "is_checked", testIsChecked, source.IsChecked)
	}
}

func TestUnmarshalJSONWithoutIsChecked(t *testing.T) {
	var (
		source = entity.Source{}
		data   = []byte(`{
			"method": "` + testMethod + `",
			"category": "` + testCategory + `",
			"url": "` + testURL + `"
		}`)
	)
	err := json.Unmarshal(data, &source)

	if err != nil {
		t.Errorf(expectedErrorButGotMessage, "unmarshal", nil, err)
	}

	if source.Method != testMethod {
		t.Errorf(expectedButGotMessage, "method", testMethod, source.Method)
	}

	if source.Category != testCategory {
		t.Errorf(expectedButGotMessage, "category", testCategory, source.Category)
	}

	if source.URL != testURL {
		t.Errorf(expectedButGotMessage, "url", testURL, source.Category)
	}

	if source.IsChecked != true {
		t.Errorf(expectedButGotMessage, "is_checked", true, source.IsChecked)
	}
}

func TestUnmarshalJSONWithInvalidData(t *testing.T) {
	var (
		source = entity.Source{}
		data   = []byte(`{
			"method": "` + testMethod + `",
			"category": "` + testCategory + `",
			"url": "` + testURL + `",
			"is_checked": "string_instead_of_bool"
		}`)
	)
	err := json.Unmarshal(data, &source)
	if err == nil {
		t.Errorf(expectedButGotMessage, "unmarshal", "any error", err)
	}
}

func TestUnmarshalJSONWithEmptyData(t *testing.T) {
	var (
		source = entity.Source{}
		data   = []byte(`{}`)
	)
	err := json.Unmarshal(data, &source)

	if err != nil {
		t.Errorf(expectedErrorButGotMessage, "unmarshal", nil, err)
	}

	if source.Method != "" {
		t.Errorf(expectedButGotMessage, "method", "empty", source.Method)
	}

	if source.Category != "" {
		t.Errorf(expectedButGotMessage, "category", "empty", source.Category)
	}

	if source.URL != "" {
		t.Errorf(expectedButGotMessage, "url", "empty", source.Category)
	}

	if source.IsChecked != true {
		t.Errorf(expectedButGotMessage, "is_checked", true, source.IsChecked)
	}
}
