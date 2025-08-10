package referral

import (
	"errors"
	"fmt"
	"sort"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrSelfReferral       = errors.New("self-referral is not allowed")
	ErrAlreadyHasReferrer = errors.New("candidate already has a referrer")
	ErrWouldCreateCycle   = errors.New("adding this referral would create a cycle")
	ErrReferralExists     = errors.New("referral already exists")
)

type Graph struct {
	users    map[string]struct{}
	parent   map[string]string
	children map[string]map[string]struct{} 
}

func NewGraph() *Graph {
	return &Graph{
		users:    make(map[string]struct{}),
		parent:   make(map[string]string),
		children: make(map[string]map[string]struct{}),
	}
}

func (g *Graph) AddUser(id string) {
	if id == "" {
		return
	}
	g.users[id] = struct{}{}
	if _, ok := g.children[id]; !ok {
		g.children[id] = make(map[string]struct{})
	}
}

func (g *Graph) HasUser(id string) bool {
	_, ok := g.users[id]
	return ok
}

func (g *Graph) DirectReferrals(id string) ([]string, error) {
	if !g.HasUser(id) {
		return nil, ErrUserNotFound
	}
	out := make([]string, 0, len(g.children[id]))
	for c := range g.children[id] {
		out = append(out, c)
	}
	return out, nil
}

func (g *Graph) AddReferral(referrer, candidate string) error {
	if !g.HasUser(referrer) {
		return fmt.Errorf("referrer %s: %w", referrer, ErrUserNotFound)
	}
	if !g.HasUser(candidate) {
		return fmt.Errorf("candidate %s: %w", candidate, ErrUserNotFound)
	}
	if referrer == candidate {
		return ErrSelfReferral
	}
	if _, has := g.parent[candidate]; has {
		if g.parent[candidate] == referrer {
			return ErrReferralExists
		}
		return ErrAlreadyHasReferrer
	}
	if g.wouldCreateCycle(referrer, candidate) {
		return ErrWouldCreateCycle
	}
	g.children[referrer][candidate] = struct{}{}
	g.parent[candidate] = referrer
	return nil
}

func (g *Graph) wouldCreateCycle(referrer, candidate string) bool {
	visited := make(map[string]struct{})
	stack := []string{candidate}
	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if curr == referrer {
			return true
		}
		if _, seen := visited[curr]; seen {
			continue
		}
		visited[curr] = struct{}{}
		for child := range g.children[curr] {
			if _, ok := visited[child]; !ok {
				stack = append(stack, child)
			}
		}
	}
	return false
}

type UserReach struct {
	User  string
	Reach int
}

func (g *Graph) TotalReach(user string) (int, error) {
	if !g.HasUser(user) {
		return 0, ErrUserNotFound
	}
	visited := make(map[string]struct{})
	queue := []string{user}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		for child := range g.children[curr] {
			if _, seen := visited[child]; !seen {
				visited[child] = struct{}{}
				queue = append(queue, child)
			}
		}
	}

	return len(visited), nil
}

func (g *Graph) TopKByReach(k int) ([]UserReach, error) {
	if k <= 0 {
		return nil, errors.New("k must be positive")
	}

	reaches := make([]UserReach, 0, len(g.users))
	for u := range g.users {
		cnt, _ := g.TotalReach(u)
		reaches = append(reaches, UserReach{User: u, Reach: cnt})
	}

	sort.Slice(reaches, func(i, j int) bool {
		if reaches[i].Reach == reaches[j].Reach {
			return reaches[i].User < reaches[j].User
		}
		return reaches[i].Reach > reaches[j].Reach
	})

	if k > len(reaches) {
		k = len(reaches)
	}
	return reaches[:k], nil
}

func (g *Graph) FullDownstream(user string) (map[string]struct{}, error) {
	if !g.HasUser(user) {
		return nil, ErrUserNotFound
	}
	visited := make(map[string]struct{})
	queue := []string{user}
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		for child := range g.children[curr] {
			if _, seen := visited[child]; !seen {
				visited[child] = struct{}{}
				queue = append(queue, child)
			}
		}
	}
	return visited, nil
}

func (g *Graph) UniqueReachExpansion() ([]UserReach, error) {
	downstreams := make(map[string]map[string]struct{})
	for u := range g.users {
		ds, _ := g.FullDownstream(u)
		downstreams[u] = ds
	}

	covered := make(map[string]struct{})
	result := []UserReach{}

	for len(covered) < len(g.users) {
		var bestUser string
		bestNewCount := -1
		for u, ds := range downstreams {
			newCount := 0
			for target := range ds {
				if _, ok := covered[target]; !ok {
					newCount++
				}
			}
			if newCount > bestNewCount {
				bestNewCount = newCount
				bestUser = u
			}
		}
		if bestNewCount <= 0 {
			break
		}
		for target := range downstreams[bestUser] {
			covered[target] = struct{}{}
		}
		result = append(result, UserReach{User: bestUser, Reach: bestNewCount})
		delete(downstreams, bestUser)
	}
	return result, nil
}

func (g *Graph) FlowCentrality() ([]UserReach, error) {
	users := make([]string, 0, len(g.users))
	for u := range g.users {
		users = append(users, u)
	}

	dist := make(map[string]map[string]int)
	for _, u := range users {
		dist[u] = bfsDistances(g, u)
	}

	score := make(map[string]int)
	for _, s := range users {
		for _, t := range users {
			if s == t || dist[s][t] == -1 {
				continue
			}
			for _, v := range users {
				if v == s || v == t {
					continue
				}
				if dist[s][v] != -1 && dist[v][t] != -1 &&
					dist[s][v]+dist[v][t] == dist[s][t] {
					score[v]++
				}
			}
		}
	}

	reaches := make([]UserReach, 0, len(score))
	for u, sc := range score {
		reaches = append(reaches, UserReach{User: u, Reach: sc})
	}

	sort.Slice(reaches, func(i, j int) bool {
		if reaches[i].Reach == reaches[j].Reach {
			return reaches[i].User < reaches[j].User
		}
		return reaches[i].Reach > reaches[j].Reach
	})

	return reaches, nil
}

func bfsDistances(g *Graph, start string) map[string]int {
	dist := make(map[string]int)
	for u := range g.users {
		dist[u] = -1
	}
	dist[start] = 0
	queue := []string{start}
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		for child := range g.children[curr] {
			if dist[child] == -1 {
				dist[child] = dist[curr] + 1
				queue = append(queue, child)
			}
		}
	}
	return dist
}

