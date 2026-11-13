package auth

import "aibuddy/internal/domain"

type Policy interface {
	Allowed(role domain.Role, tool string) bool
}

type StaticPolicy map[domain.Role]map[string]struct{}

func (p StaticPolicy) Allowed(role domain.Role, tool string) bool {
	_, ok := p[role][tool]
	return ok
}

func set(names ...string) map[string]struct{} {
	m := make(map[string]struct{}, len(names))
	for _, n := range names {
		m[n] = struct{}{}
	}
	return m
}

func DefaultPolicy() StaticPolicy {
	return StaticPolicy{
		domain.RoleIntern: set(
			"get_my_progress", "search_kb", "get_who_to_ask",
			"create_absence_request", "get_absence_status", "request_access",
		),
		domain.RoleMentor:  set("get_my_progress", "search_kb", "get_who_to_ask"),
		domain.RoleManager: set("get_my_progress", "search_kb", "get_who_to_ask"),
		domain.RoleHRBP:    set("search_kb", "get_who_to_ask"),
	}
}
