package application

import (
	"sort"
	"strconv"
	"strings"

	"github.com/rubenbuelvas/game-organizer-be/src/domain"
)

type OrganizerService struct {
}

func NewOrganizerService() *OrganizerService {
	return &OrganizerService{}
}

func (s *OrganizerService) Organize(dto OrganizerDTO) ([]domain.Team, error) {
	teams := GenerateEmptyTeams(dto.NumberOfTeams)

	if !dto.BreakFamilies {
		// Keep families together
		organizeWithFamiliesTogether(dto, teams)
	} else {
		// Separate families
		organizeWithFamiliesSeparated(dto, teams)
	}

	return CollectTeams(teams), nil
}

// organizeWithFamiliesTogether groups players by family and assigns entire families to teams
func organizeWithFamiliesTogether(dto OrganizerDTO, teams map[int]domain.Team) {
	// Group players by family
	familyGroups := make(map[string][]domain.Player)
	for _, player := range dto.Players {
		familyKey := player.Family.Name
		familyGroups[familyKey] = append(familyGroups[familyKey], player)
	}

	// Convert map to slice for deterministic ordering
	var families [][]domain.Player
	for _, group := range familyGroups {
		families = append(families, group)
	}

	// Sort families by size (larger families first) for better distribution
	sort.Slice(families, func(i, j int) bool {
		return len(families[i]) > len(families[j])
	})

	// Assign families to teams trying to balance based on priority
	for _, family := range families {
		assignFamilyToTeam(dto, teams, family)
	}
}

// assignFamilyToTeam assigns a family group to the team with the lowest target metric
func assignFamilyToTeam(dto OrganizerDTO, teams map[int]domain.Team, family []domain.Player) {
	var targetTeamID int
	var lowestValue float64 = float64(^uint64(0) >> 1) // Max float

	// Find team with lowest target value (or smallest size if priority is balanced)
	for teamID, team := range teams {
		var value float64

		if dto.Priority == "skills" {
			value = float64(getTeamSkillSum(team, dto.SkillBalanceType)) + float64(len(team.Members))*0.1
		} else if dto.Priority == "stats" {
			value = float64(getTeamStatSum(team, dto.StatBalanceType)) + float64(len(team.Members))*0.1
		} else {
			// Default: balance by team size
			value = float64(len(team.Members))
		}

		if value < lowestValue {
			lowestValue = value
			targetTeamID = teamID
		}
	}

	// Add all family members to the target team
	for _, player := range family {
		team := teams[targetTeamID]
		team.Members = append(team.Members, player)
		teams[targetTeamID] = team
	}
}

// organizeWithFamiliesSeparated assigns individual players to teams while trying to separate families
func organizeWithFamiliesSeparated(dto OrganizerDTO, teams map[int]domain.Team) {
	// Sort players by priority value (in reverse) for better distribution
	players := make([]domain.Player, len(dto.Players))
	copy(players, dto.Players)

	sort.Slice(players, func(i, j int) bool {
		var valI, valJ float64

		if dto.Priority == "skills" {
			valI = float64(getPlayerSkillValue(players[i], dto.SkillBalanceType))
			valJ = float64(getPlayerSkillValue(players[j], dto.SkillBalanceType))
		} else if dto.Priority == "stats" {
			valI = float64(getPlayerStatValue(players[i], dto.StatBalanceType))
			valJ = float64(getPlayerStatValue(players[j], dto.StatBalanceType))
		}

		return valI > valJ // Sort in descending order
	})

	// Assign players to teams
	for _, player := range players {
		assignPlayerToTeam(dto, teams, player)
	}
}

// assignPlayerToTeam assigns a single player to the team with the lowest target metric
func assignPlayerToTeam(dto OrganizerDTO, teams map[int]domain.Team, player domain.Player) {
	var targetTeamID int
	var lowestValue float64 = float64(^uint64(0) >> 1)

	// Find team with lowest target value, preferring teams without family members
	for teamID, team := range teams {
		var value float64

		if dto.Priority == "skills" {
			value = float64(getTeamSkillSum(team, dto.SkillBalanceType))
		} else if dto.Priority == "stats" {
			value = float64(getTeamStatSum(team, dto.StatBalanceType))
		} else {
			value = float64(len(team.Members))
		}

		// Add penalty for family members already in team
		familyInTeam := countFamilyInTeam(team, player.Family.Name)
		value += float64(familyInTeam) * 100 // High penalty to separate families

		if value < lowestValue {
			lowestValue = value
			targetTeamID = teamID
		}
	}

	// Add player to the target team
	team := teams[targetTeamID]
	team.Members = append(team.Members, player)
	teams[targetTeamID] = team
}

// getPlayerSkillValue extracts the value of a specific skill from a player
func getPlayerSkillValue(player domain.Player, skillType string) int {
	switch strings.ToLower(skillType) {
	case "agility":
		return player.Skills.Agility
	case "intelligence":
		return player.Skills.Intelligence
	case "artistic":
		return player.Skills.Artistic
	case "communication":
		return player.Skills.Communication
	default:
		return 0
	}
}

// getPlayerStatValue extracts the value of a specific stat from a player
func getPlayerStatValue(player domain.Player, statType string) int {
	switch strings.ToLower(statType) {
	case "social_energy", "socialenergy":
		return player.Stats.SocialEnergy
	case "age":
		return player.Stats.Age
	default:
		return 0
	}
}

// getTeamSkillSum calculates the sum of a specific skill for all team members
func getTeamSkillSum(team domain.Team, skillType string) int {
	sum := 0
	for _, player := range team.Members {
		sum += getPlayerSkillValue(player, skillType)
	}
	return sum
}

// getTeamStatSum calculates the sum of a specific stat for all team members
func getTeamStatSum(team domain.Team, statType string) int {
	sum := 0
	for _, player := range team.Members {
		sum += getPlayerStatValue(player, statType)
	}
	return sum
}

// countFamilyInTeam counts how many members of a family are in a team
func countFamilyInTeam(team domain.Team, familyName string) int {
	count := 0
	for _, player := range team.Members {
		if player.Family.Name == familyName {
			count++
		}
	}
	return count
}

// countFamilyMembers counts how many members of a family are in any team
func countFamilyMembers(teams map[int]domain.Team, familyName string) int {
	count := 0
	for _, team := range teams {
		count += countFamilyInTeam(team, familyName)
	}
	return count
}

func GenerateEmptyTeams(numberOfTeams int) map[int]domain.Team {
	teams := make(map[int]domain.Team)
	for i := 1; i <= numberOfTeams; i++ {
		teams[i] = domain.Team{
			Name:    "Team " + strconv.Itoa(i),
			Members: []domain.Player{},
		}
	}
	return teams
}

func CollectTeams(teamMap map[int]domain.Team) []domain.Team {
	teamList := make([]domain.Team, 0, len(teamMap))
	for _, team := range teamMap {
		teamList = append(teamList, team)
	}
	return teamList
}
