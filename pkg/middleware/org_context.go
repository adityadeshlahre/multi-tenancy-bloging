package middleware

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/adityadeshlahre/multi-tenant-backend-app/repository"
)

const OrganizationKey = "organizationId"

func OrganizationContext(orgRepo repository.OrgRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgIdHeader := c.GetHeader("X-Organization-ID")
		if orgIdHeader == "" {
			log.Println("Organization ID not found in header")
			c.AbortWithStatusJSON(400, gin.H{"error": "Organization ID header required"})
			return
		}

		orgId, err := strconv.ParseUint(orgIdHeader, 10, 32)
		if err != nil {
			log.Printf("Invalid organization ID format: %v", err)
			c.AbortWithStatusJSON(400, gin.H{"error": "Invalid organization ID format"})
			return
		}

		org, err := orgRepo.GetOrganizationByID(c.Request.Context(), uint(orgId))
		if err != nil {
			log.Printf("Error fetching organization: %v", err)
			c.AbortWithStatusJSON(404, gin.H{"error": "Organization not found"})
			return
		}

		c.Set("organization", org)
		c.Set(OrganizationKey, uint(orgId))
		c.Next()
	}
}
