package repository

import (
	"strings"
	"testing"
)

func TestScoreboardQueriesKeepRegistrationsWithMissingInstitution(t *testing.T) {
	queries := []string{
		scoreboardByEventLevelQuery(),
		customLeaderboardQuery(),
	}

	for _, query := range queries {
		if strings.Contains(query, "\n\t\tJOIN institutions i ON t.insti_id = i.id") {
			t.Fatal("scoreboard queries must not inner join institutions because it drops registrations with incomplete institution data")
		}

		if !strings.Contains(query, "LEFT JOIN institutions i ON t.insti_id = i.id") {
			t.Fatal("scoreboard queries must left join institutions")
		}

		if !strings.Contains(query, "r.is_kick = false") {
			t.Fatal("scoreboard queries must exclude kicked registrations")
		}
	}
}
