package swashengine

import (
	"tsundere/packages/customization/gradient"
	"tsundere/source/master/sessions"
)

func (s *SwashEngine) Elements(extras map[string]any) map[string]any {
	var mainMap = map[string]any{
		"user":     s.session.UserProfile,
		"gradient": gradient.New,

		"sessions": map[string]any{
			"length": sessions.Count(),
		},
	}

	if extras != nil {
		for key, value := range extras {
			mainMap[key] = value
		}
	}

	return mainMap
}

func (s *SwashEngine) Packages() map[string]map[string]any {
	return map[string]map[string]any{}
}
