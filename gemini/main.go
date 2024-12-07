package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/transport"
)

const apiEndpoint = "https://discoveryengine.googleapis.com/v1/projects/%s/locations/%s/collections/default_collection/engines/%s/servingConfigs/default_config:search"

type SearchRequest struct {
	Query               string `json:"query"`
	PageSize            int    `json:"pageSize"`
	ContentSearchSpec   `json:"contentSearchSpec"`
	QueryExpansionSpec  `json:"queryExpansionSpec"`
	SpellCorrectionSpec `json:"spellCorrectionSpec"`
}

type ContentSearchSpec struct {
	SnippetSpec SnippetSpec `json:"snippetSpec"`
	SummarySpec SummarySpec `json:"summarySpec"`
}

type SnippetSpec struct {
	ReturnSnippet bool `json:"returnSnippet"`
}

type SummarySpec struct {
	SummaryResultCount           int  `json:"summaryResultCount"`
	IncludeCitations             bool `json:"includeCitations"`
	IgnoreAdversarialQuery       bool `json:"ignoreAdversarialQuery"`
	IgnoreNonSummarySeekingQuery bool `json:"ignoreNonSummarySeekingQuery"`
	ModelPromptSpec              `json:"modelPromptSpec"`
	ModelSpec                    `json:"modelSpec"`
}

type ModelPromptSpec struct {
	Preamble string `json:"preamble"`
}

type ModelSpec struct {
	Version string `json:"version"`
}

type QueryExpansionSpec struct {
	Condition string `json:"condition"`
}

type SpellCorrectionSpec struct {
	Mode string `json:"mode"`
}

func searchSample(projectID, location, engineID, searchQuery string) (string, error) {
	url := fmt.Sprintf(apiEndpoint, projectID, location, engineID)

	requestBody := SearchRequest{
		Query:    searchQuery,
		PageSize: 10,
		ContentSearchSpec: ContentSearchSpec{
			SnippetSpec: SnippetSpec{
				ReturnSnippet: true,
			},
			SummarySpec: SummarySpec{
				SummaryResultCount:           5,
				IncludeCitations:             true,
				IgnoreAdversarialQuery:       true,
				IgnoreNonSummarySeekingQuery: true,
				ModelPromptSpec: ModelPromptSpec{
					Preamble: "YOUR_CUSTOM_PROMPT",
				},
				ModelSpec: ModelSpec{
					Version: "stable",
				},
			},
		},
		QueryExpansionSpec: QueryExpansionSpec{
			Condition: "AUTO",
		},
		SpellCorrectionSpec: SpellCorrectionSpec{
			Mode: "AUTO",
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}
	ctx := context.Background()
	client, _, err := transport.NewHTTPClient(ctx, option.WithScopes("https://www.googleapis.com/auth/cloud-platform"))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP client: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 response: %v\n%s", resp.StatusCode, body)
	}

	return string(body), nil
}

func main() {
	projectID := "term6-joji-wada"                   // 例
	location := "global"                             // 例
	engineID := "hackathon-curriculum_1733546607357" // 例
	searchQuery := "入試日程"                            // 例

	response, err := searchSample(projectID, location, engineID, searchQuery)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during search: %v\n", err)
		return
	}

	fmt.Printf("Search Response: %s\n", response)
}
