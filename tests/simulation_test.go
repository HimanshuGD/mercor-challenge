package tests

import (
	"math"
	"testing"

	simulation "mercor-challenge/source"

	"github.com/stretchr/testify/assert"
)

func TestSimulateBasic(t *testing.T) {
	p := 1.0
	days := 5
	res := simulation.Simulate(p, days)
	assert.Len(t, res, days)
	assert.True(t, res[0] > 0)
	for i := 1; i < len(res); i++ {
		assert.True(t, res[i] >= res[i-1])
	}
}

func TestDaysToTarget(t *testing.T) {
	p := 0.4
	target := 50
	days := simulation.DaysToTarget(p, target)
	assert.True(t, days > 0)
}

func adoptionProbLinear(bonus float64) float64 {
	p := bonus / 100.0
	if p > 1.0 {
		return 1.0
	}
	if p < 0.0 {
		return 0.0
	}
	return p
}

func TestMinBonusForTarget_Simple(t *testing.T) {
	days := 1
	target := 50
	eps := 1e-3
	res, err := simulation.MinBonusForTarget(days, target, adoptionProbLinear, eps)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 50, *res)
}

func TestMinBonusForTarget_Unachievable(t *testing.T) {
	adopt := func(bonus float64) float64 {
		return math.Min(1e-6, bonus*0.0)
	}
	days := 10
	target := 1000
	res, err := simulation.MinBonusForTarget(days, target, adopt, 1e-3)
	assert.NoError(t, err)
	assert.Nil(t, res)
}