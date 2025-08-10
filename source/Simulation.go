package referral

import (
	"errors"
	"math"
)

func Simulate(p float64, days int) []float64 {
    const initialReferrers = 100
    const referralCapacity = 10

    if days > 500 && p > 0.5 {
        days = 500
    }

    type referrer struct {
        remaining float64
    }

    active := []referrer{}
    for i := 0; i < initialReferrers; i++ {
        active = append(active, referrer{remaining: referralCapacity})
    }

    result := make([]float64, days)
    var total float64

    for day := 0; day < days; day++ {
        newReferrers := []referrer{}
        for i := range active {
            if active[i].remaining <= 0 {
                continue
            }
            referralsToday := p
            if referralsToday > active[i].remaining {
                referralsToday = active[i].remaining
            }
            active[i].remaining -= referralsToday
            total += referralsToday
            for r := 0; r < int(referralsToday+0.5); r++ {
                newReferrers = append(newReferrers, referrer{remaining: referralCapacity})
            }
        }
        active = append(active, newReferrers...)
        result[day] = total
    }
    return result
}


func DaysToTarget(p float64, targetTotal int) int {
	const maxDays = 10000
	sim := Simulate(p, maxDays)
	for day, total := range sim {
		if int(total+0.5) >= targetTotal {
			return day + 1 
		}
	}
	return -1 
}

func MinBonusForTarget(days int, target int, adoptionProb func(bonus float64) float64, eps float64) (*int, error) {
	if days <= 0 {
		return nil, errors.New("days must be positive")
	}
	if target <= 0 {
		zero := 0
		return &zero, nil
	}
	if adoptionProb == nil {
		return nil, errors.New("adoptionProb cannot be nil")
	}
	if eps <= 0 {
		eps = 1e-3
	}

	isAchievable := func(bonus float64) bool {
		p := adoptionProb(bonus)
		if p < 0 {
			p = 0
		}
		if p > 1 {
			p = 1
		}
		res := Simulate(p, days)
		if len(res) == 0 {
			return false
		}
		last := res[len(res)-1]
		return int(math.Floor(last+0.5)) >= target
	}

	if isAchievable(0.0) {
		z := 0
		return &z, nil
	}

	low := 0.0
	high := 10.0 
	const maxBonus = 1_000_000.0 
	for !isAchievable(high) {
		high *= 2
		if high > maxBonus {
			if !isAchievable(maxBonus) {
				return nil, nil 
			}
			high = maxBonus
			break
		}
	}

	// binary search on [low, high] for minimal bonus (real-valued)
	for high-low > eps {
		mid := (low + high) / 2.0
		if isAchievable(mid) {
			high = mid
		} else {
			low = mid
		}
	}

	// round up to nearest $10
	bonusRounded := math.Ceil(high/10.0) * 10.0
	// verify rounded value is achievable; if not, try next increments of $10
	for ; bonusRounded <= maxBonus; bonusRounded += 10.0 {
		if isAchievable(bonusRounded) {
			ans := int(bonusRounded + 0.5)
			return &ans, nil
		}
	}
	// if we exit loop, target not achievable
	return nil, nil
}
