package middleware

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iamsainadh/neovault/api-gateway/internal/repository"
)

func OrganizationRoleMiddleware(repo repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")

		orgID := c.Param("orgId")

		role, err := repo.GetUserRole(c.Request.Context(), userID, orgID)

		log.Println("Role" + role)

		if err != nil {
			c.AbortWithStatusJSON(403, gin.H{"error": "not a member"})
			return
		}

		c.Set("role", role)

		fmt.Println(role)

		c.Next()
	}
}
