package oauth2

import "time"

const noExpiryYears = 100

// ResolveExpiresAt returns a safe expiration timestamp for tokens.
// If expiresIn is <= 0, we treat it as non-expiring and set a far-future time.
func ResolveExpiresAt(expiresIn int) time.Time {
	if expiresIn <= 0 {
		return time.Now().AddDate(noExpiryYears, 0, 0)
	}
	return time.Now().Add(time.Duration(expiresIn) * time.Second)
}
