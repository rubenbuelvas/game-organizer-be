package application

import (
	"math"
	"sort"
	"strconv"

	"github.com/rubenbuelvas/game-organizer-be/src/domain"
)

type OrganizerService struct {
}

func NewOrganizerService() *OrganizerService {
	return &OrganizerService{}
}

func (s *OrganizerService) Organize(dto OrganizerDTO) ([]domain.Team, error) {
	teams := generateEmptyTeams(dto.NumberOfTeams)

	if !dto.BreakFamilies {
		organizeWithFamiliesTogether(dto, teams)
	} else {
		organizeWithFamiliesSeparated(dto, teams)
	}

	for teamId, team := range teams {
		team.Score = getTeamWeightedStatSum(team, dto.Stats)
		teams[teamId] = team
	}

	return collectTeams(teams), nil
}

// organizeWithFamiliesTogether groups players by family and assigns entire families to teams
func organizeWithFamiliesTogether(dto OrganizerDTO, teams map[int]domain.Team) {
	// Group players by family
	familyGroups := make(map[int64][]domain.Player)
	for _, player := range dto.Players {
		familyKey := player.Family.ID
		familyGroups[familyKey] = append(familyGroups[familyKey], player)
	}

	// Convert map to slice for deterministic ordering
	var families [][]domain.Player
	for _, group := range familyGroups {
		families = append(families, group)
	}

	// Sort families by total weighted stat (larger first) for better distribution
	sort.Slice(families, func(i, j int) bool {
		return familyWeightedValue(families[i], dto.Stats) > familyWeightedValue(families[j], dto.Stats)
	})

	// Assign families to teams trying to balance based on weighted stats
	for _, family := range families {
		assignFamilyToTeam(dto, teams, family)
	}
}

// assignFamilyToTeam assigns a family group to the team with the lowest target metric
func assignFamilyToTeam(dto OrganizerDTO, teams map[int]domain.Team, family []domain.Player) {
	var targetTeamID int
	var lowestValue float64 = math.MaxFloat64

	// Find team with lowest target value (or smallest size if balanced)
	for teamID, team := range teams {
		value := float64(getTeamWeightedStatSum(team, dto.Stats)) + float64(len(team.Members))*0.1
		if value < lowestValue {
			lowestValue = value
			targetTeamID = teamID
		}
	}

	for _, player := range family {
		team := teams[targetTeamID]
		team.Members = append(team.Members, player)
		teams[targetTeamID] = team
	}
}

// organizeWithFamiliesSeparated assigns individual players to teams while trying to separate families
func organizeWithFamiliesSeparated(dto OrganizerDTO, teams map[int]domain.Team) {
	players := make([]domain.Player, len(dto.Players))
	copy(players, dto.Players)

	sort.Slice(players, func(i, j int) bool {
		return getPlayerWeightedStatValue(players[i], dto.Stats) > getPlayerWeightedStatValue(players[j], dto.Stats)
	})

	for _, player := range players {
		assignPlayerToTeam(dto, teams, player)
	}
}

// assignPlayerToTeam assigns a single player to the team with the lowest target metric
func assignPlayerToTeam(dto OrganizerDTO, teams map[int]domain.Team, player domain.Player) {
	var targetTeamID int
	var lowestValue float64 = math.MaxFloat64

	for teamID, team := range teams {
		value := float64(getTeamWeightedStatSum(team, dto.Stats))
		familyInTeam := countFamilyInTeam(team, player.Family.ID)
		value += float64(familyInTeam) * 100
		if value < lowestValue {
			lowestValue = value
			targetTeamID = teamID
		}
	}

	team := teams[targetTeamID]
	team.Members = append(team.Members, player)
	teams[targetTeamID] = team
}

// getPlayerWeightedStatValue computes a player's total stats weighted by importance
func getPlayerWeightedStatValue(player domain.Player, weights domain.Stats) int {
	return player.Stats.Agility*weights.Agility +
		player.Stats.Intelligence*weights.Intelligence +
		player.Stats.ArtisticSkill*weights.ArtisticSkill +
		player.Stats.Communication*weights.Communication +
		player.Stats.SocialEnergy*weights.SocialEnergy
}

// getTeamWeightedStatSum calculates sum of weighted stat values for all team members
func getTeamWeightedStatSum(team domain.Team, weights domain.Stats) int {
	sum := 0
	for _, player := range team.Members {
		sum += getPlayerWeightedStatValue(player, weights)
	}
	return sum
}

// familyWeightedValue returns the combined weighted stat value for all players in a family
func familyWeightedValue(family []domain.Player, weights domain.Stats) int {
	total := 0
	for _, player := range family {
		total += getPlayerWeightedStatValue(player, weights)
	}
	return total
}

// countFamilyInTeam counts how many members of a family are in a team
func countFamilyInTeam(team domain.Team, familyId int64) int {
	count := 0
	for _, player := range team.Members {
		if player.Family.ID == familyId {
			count++
		}
	}
	return count
}

func generateEmptyTeams(numberOfTeams int) map[int]domain.Team {
	teams := make(map[int]domain.Team)
	for i := 1; i <= numberOfTeams; i++ {
		teams[i] = domain.Team{
			Name:    "Team " + strconv.Itoa(i),
			Members: []domain.Player{},
		}
	}
	return teams
}

func collectTeams(teamMap map[int]domain.Team) []domain.Team {
	teamList := make([]domain.Team, 0, len(teamMap))
	for _, team := range teamMap {
		teamList = append(teamList, team)
	}
	return teamList
}
