# Mercor Challenge – Referral Network & Simulation

## Overview
This project implements:
1. **Referral Network** – A directed acyclic graph (DAG) of users, with rules for adding referrals, detecting cycles, tracking reach, and computing centrality metrics.
2. **Simulation Engine** – A referral adoption simulation to model user growth over time, calculate days to reach a target, and determine the minimal bonus needed to achieve adoption goals.

## Structure

mercor-challenge/
├── source/
│ ├── ReferralNetwork.go # Graph logic & reach calculations
│ └── simulation.go # Growth simulation & bonus calculation
└── tests/
├── referral_network_test.go # Unit tests for referral graph
└── simulation_test.go # Unit tests for simulation logic


## Approach
- **Referral Graph**:  
  - Used adjacency maps for `children` and a `parent` map for quick validations.  
  - Rules enforced:
    - No self-referrals
    - Each user can only have one referrer
    - No cycles allowed
    - Duplicate referral prevention
  - Implemented BFS-based reach calculations and sorting for top-K queries.
  - Centrality computed via shortest path counts.

- **Simulation**:  
  - Modeled active referrers with limited referral capacity.
  - At each day, new referrers are added based on adoption probability `p`.
  - `DaysToTarget` simulates until target is met or max days reached.
  - `MinBonusForTarget` uses binary search to find the minimal bonus that achieves the target within given days.

- **Testing**:  
  - Used `testify/assert` for concise assertions.
  - Covered both happy paths and failure conditions (rule violations, unreachable targets).
  - Simulation tests use controlled parameters to avoid long runtimes.

## How to Run
```bash
# Run all tests
go test ./tests -v



Time Spent
Approx. 4.5 hours:

1.5h – Referral graph implementation & validations

1.5h – Simulation logic & bonus search

1h – Writing and refining tests

0.5h – Project structuring & cleanup

