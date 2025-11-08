package jwt

import "context"

func GetContextToken(ctx context.Context) string {
	token := ""
	if v, ok := ctx.Value("bearer_token").(string); ok {
		token = v
	} else {
		if v, ok = ctx.Value("auth_token").(string); ok {
			token = v
		} else {
			if v, ok = ctx.Value("access_key").(string); ok {
				token = v
			}
		}
	}
	return token
}

func GetContextSubject(ctx context.Context) string {
	if v, ok := ctx.Value("jwt:subject").(string); ok {
		return v
	}
	return ""
}
