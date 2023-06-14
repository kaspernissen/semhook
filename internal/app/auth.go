package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

func authMiddleware(issuer string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		jwksURL := "https://github.com/.well-known/jwks.json"
		set, err := jwk.Fetch(context.Background(), jwksURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch JWKS keys"})
			c.Abort()
			return
		}

		token, err := jwt.Parse([]byte(tokenString),
			jwt.WithKeySet(set),
			jwt.WithValidate(true),
			jwt.WithIssuer(issuer),
			jwt.WithAudience("")) //TODO figure out audience
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims := token.PrivateClaims()
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		/*
			requiredScopes := getRequiredScopes(c.Request.URL.Path)
			scopes, ok := claims.Get("scope").(jwt.Set)
			if !ok {
				c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient scope"})
				c.Abort()
				return
			}

			if !hasRequiredScopes(requiredScopes, scopes) {
				c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient scope"})
				c.Abort()
				return
			}
		*/
		c.Set("user", token)
		c.Next()
	}
}
