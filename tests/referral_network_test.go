package tests

import (
	"testing"

	referral "mercor-challenge/source"

	"github.com/stretchr/testify/assert"
)

func TestGraph_AddReferral_Rules(t *testing.T) {
    tests := []struct {
        name       string
        setup      func(*referral.Graph)
        referrer   string
        candidate  string
        expectErr  error
    }{
        {
            name: "Valid referral",
            setup: func(g *referral.Graph) {
                // no pre-existing edge
                g.AddUser("A")
                g.AddUser("B")
            },
            referrer: "A",
            candidate: "B",
            expectErr: nil,
        },
        {
            name: "Self referral",
            setup: func(g *referral.Graph) {
                g.AddUser("A")
            },
            referrer: "A",
            candidate: "A",
            expectErr: referral.ErrSelfReferral,
        },
        {
            name: "Unique referrer violation",
            setup: func(g *referral.Graph) {
                g.AddUser("A")
                g.AddUser("B")
                g.AddUser("C")
                _ = g.AddReferral("A", "B")
            },
            referrer: "C",
            candidate: "B",
            expectErr: referral.ErrAlreadyHasReferrer,
        },
        {
            name: "Cycle detection",
            setup: func(g *referral.Graph) {
                g.AddUser("A")
                g.AddUser("B")
                g.AddUser("C")
                _ = g.AddReferral("A", "B")
                _ = g.AddReferral("B", "C")
            },
            referrer: "C",
            candidate: "A",
            expectErr: referral.ErrWouldCreateCycle,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            g := referral.NewGraph()
            tc.setup(g)
            err := g.AddReferral(tc.referrer, tc.candidate)
            if tc.expectErr == nil {
                assert.NoError(t, err)
            } else {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tc.expectErr.Error())
            }
        })
    }
}


func TestGraph_DirectReferrals(t *testing.T) {
	g := referral.NewGraph()
	g.AddUser("R")
	g.AddUser("C1")
	g.AddUser("C2")
	assert.NoError(t, g.AddReferral("R", "C1"))
	assert.NoError(t, g.AddReferral("R", "C2"))

	children, err := g.DirectReferrals("R")
	assert.NoError(t, err)
	assert.Len(t, children, 2)
}

func TestTotalReach(t *testing.T) {
	g := referral.NewGraph()
	for _, u := range []string{"A", "B", "C", "D"} {
		g.AddUser(u)
	}
	_ = g.AddReferral("A", "B")
	_ = g.AddReferral("B", "C")
	_ = g.AddReferral("C", "D")

	count, err := g.TotalReach("A")
	assert.NoError(t, err)
	assert.Equal(t, 3, count)

	count, err = g.TotalReach("B")
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	_, err = g.TotalReach("X")
	assert.Error(t, err)
}

func TestTopKByReach(t *testing.T) {
	g := referral.NewGraph()
	for _, u := range []string{"A", "B", "C", "D", "E"} {
		g.AddUser(u)
	}
	_ = g.AddReferral("A", "B")
	_ = g.AddReferral("B", "C")
	_ = g.AddReferral("A", "D")
	_ = g.AddReferral("B", "E")

	top, err := g.TopKByReach(3)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(top))
	assert.Equal(t, "A", top[0].User)
	assert.Equal(t, 4, top[0].Reach)
}

func TestUniqueReachExpansion(t *testing.T) {
	g := referral.NewGraph()
	for _, u := range []string{"A", "B", "C", "D", "E"} {
		g.AddUser(u)
	}
	_ = g.AddReferral("A", "B")
	_ = g.AddReferral("A", "C")
	_ = g.AddReferral("B", "D")
	_ = g.AddReferral("C", "E")

	result, err := g.UniqueReachExpansion()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(result), 1)
	assert.Equal(t, "A", result[0].User)
}

func TestFlowCentrality(t *testing.T) {
	g := referral.NewGraph()
	for _, u := range []string{"A", "B", "C", "D"} {
		g.AddUser(u)
	}
	_ = g.AddReferral("A", "B")
	_ = g.AddReferral("B", "C")
	_ = g.AddReferral("A", "D")
	_ = g.AddReferral("D", "C")

	result, err := g.FlowCentrality()
	assert.NoError(t, err)
	assert.True(t, len(result) > 0)
	assert.Contains(t, []string{"B", "D"}, result[0].User)
}
