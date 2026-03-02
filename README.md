# game-organizer-be

This service organizes players into balanced teams using weighted player statistics and simple heuristics to respect or separate families.

**Algorithm (organizer service)**

The core implementation is in [src/application/organizer_service.go](src/application/organizer_service.go#L1-L200).

- **Inputs:** number of teams, list of players (each with stats and family ID), and a `Stats` weight object that assigns importance to each stat.
- **Weighted value:** each player's score is computed as a linear combination of their stats multiplied by the provided weights.

- **Mode A — Families Together:**
	- Group players by their family ID.
	- Compute the family's total weighted value (sum of its members' weighted values).
	- Sort families descending by that total so larger/stronger families are placed first.
	- For each family, pick the team with the lowest current metric (team's weighted sum plus a small size factor) and assign the whole family to that team. This helps distribute family strength across teams while keeping families intact.

- **Mode B — Families Separated:**
	- Treat each player individually and sort players descending by their weighted value.
	- For each player, evaluate teams by their current weighted sum and add a large penalty proportional to how many members of the player's family are already on that team. This strongly discourages placing same-family members together and spreads them across teams.

- **Helpers:**
	- `GenerateEmptyTeams(n)` creates `n` teams with empty member lists.
	- `getPlayerWeightedStatValue`, `getTeamWeightedStatSum`, and `familyWeightedValue` compute the numeric metrics used by the heuristics.
	- `CollectTeams` converts the internal team map back to a slice for output.

These heuristics prioritize balance by weighted stats while offering a simple, deterministic strategy for family handling. They are efficient (sorting + linear assignment) and easy to adjust by changing the stat weights or the family penalty/size factors.