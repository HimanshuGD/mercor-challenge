```markdown
# Mercor Challenge – Referral Network

Go implementation of the Mercor Take-Home Challenge (Parts 1–5).  
Covers referral graph logic, reach calculations, expansion strategy, flow centrality, and growth simulation.

## Structure
```
mercor-challenge/
├── source/        # Core logic
│   ├── ReferralNetwork.go  # Graph & reach functions
│   └── simulation.go       # Simulation & bonus calc
└── tests/                  # Unit tests
```

## Features
- **Referral rules**: no self-referral, no cycles, one referrer per user.
- Reach metrics: total reach, top-K by reach.
- Unique reach expansion & flow centrality.
- Simulation for adoption growth.
- Days-to-target & minimum bonus calculation.
- Full test coverage with `testify`.

## Run
```bash
go mod tidy
go test ./...
```
