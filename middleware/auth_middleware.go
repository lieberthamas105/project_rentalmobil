package middleware

import (
	"net/http"
	"project_rentalmobil/model"
	"project_rentalmobil/utils/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	RequireToken(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService service.JwtService
}

type authHeader struct {
	AuthorizationHeader string `header:"Authorization" binding:"required"`
}

func (a *authMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var aH authHeader
		err := ctx.ShouldBindHeader(&aH)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		// Extract token from Authorization header
		token := strings.Replace(aH.AuthorizationHeader, "Bearer ", "", 1)
		tokenClaim, err := a.jwtService.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		// Convert UserId from string to uint32
		userId, err := strconv.ParseUint(tokenClaim.UserId, 10, 32)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			return
		}

		// Validate roles (only if roles are provided)
		validRole := len(roles) == 0 // Allow access if no roles are required
		for _, role := range roles {
			if role == tokenClaim.Role {
				validRole = true
				break
			}
		}
		if !validRole {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden Resource"})
			return
		}

		// Set user information to context (used in other handlers)
		ctx.Set("user", model.UserCredential{
			Id:   uint32(userId), // Convert userId to uint32
			Role: tokenClaim.Role,
		})

		ctx.Next()
	}
}

func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}
