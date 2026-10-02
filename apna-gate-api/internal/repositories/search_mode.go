package repository

func normalizeSearchMode(mode string, allowed ...string) string {
	if mode == "" {
		return "all"
	}
	for _, allowedMode := range allowed {
		if mode == allowedMode {
			return mode
		}
	}
	return "all"
}
