package middleware

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/adityadeshlahre/multi-tenant-backend-app/repository"
)

const OrganizationKey = "organizationId"

func OrganizationContext(orgRepo repository.OrgRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgId, exists := c.Get(OrganizationKey)
		if !exists {
			log.Println("Organization ID not found in context")
			c.AbortWithStatusJSON(400, gin.H{"error": "Organization ID not found"})
			return
		}

		org, err := orgRepo.GetOrganizationByID(c.Request.Context(), orgId.(uint))
		if err != nil {
			log.Printf("Error fetching organization: %v", err)
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
			return
		}

		c.Set("organization", org)
		c.Next()
	}
}
