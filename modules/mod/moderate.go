package main

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/oauth2/google"
	language "google.golang.org/api/language/v1"
)

func ModerateText(text string) ([]Category, error) {
	ctx := context.Background()

	println("Moderating text... ", text, "\n")

	httpClient, err := google.DefaultClient(ctx, language.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("failed to get HTTP client: %v", err)
	}

	nlpService, err := language.New(httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create NLP service: %v", err)
	}

	req := &language.ModerateTextRequest{
		Document: &language.Document{
			Type:    "PLAIN_TEXT",
			Content: text,
		},
	}

	resp, err := nlpService.Documents.ModerateText(req).Do()
	if err != nil {
		return nil, fmt.Errorf("moderateText API call failed: %v", err)
	}

	var results []Category
	for _, cat := range resp.ModerationCategories {
		results = append(results, Category{
			Name:       cat.Name,
			Confidence: cat.Confidence,
		})
	}

	return results, nil
}

func getCategoryWeight(name string) float64 {
	switch strings.ToLower(name) {
	case "toxic":
		return ToxicWeight
	case "insult":
		return InsultWeight
	case "profanity":
		return ProfanityWeight
	case "derogatory":
		return DerogatoryWeight
	case "sexual":
		return SexualWeight
	case "death, harm & tragedy":
		return DeathHarmTragedyWeight
	case "violent":
		return ViolentWeight
	case "firearms & weapons":
		return FirearmsWeaponsWeight
	case "public safety":
		return PublicSafetyWeight
	case "health":
		return HealthWeight
	case "religion & belief":
		return ReligionBeliefWeight
	case "illicit drugs":
		return DrugsWeight
	case "war & conflict":
		return WarConflictWeight
	case "politics":
		return PoliticsWeight
	case "finance":
		return FinanceWeight
	case "legal":
		return LegalWeight
	default:
		return 0.0
	}
}
