// Package telemetry provides metrics collection and telemetry filtering capabilities.
package telemetry

import (
	"net"
	"net/http"
	"strings"

	"digital-innovation/gostrategy/internal/logging"
	"digital-innovation/gostrategy/internal/utils"

	"github.com/gin-gonic/gin"
)

var (
	ginUnauthorizedIPaddress  = gin.H{"error": "unauthorized IP address"}
	ginInvalidClientIPaddress = gin.H{"error": "invalid client ip"}
)

// IPFilterMiddleware restricts endpoint access to whitelisted IPs/CIDRs from env.
func IPFilterMiddleware() gin.HandlerFunc {
	allowedIPsStr := utils.GetEnv("TELEMETRY_ALLOWED_IPS", "")

	if allowedIPsStr == "" {
		logging.Error("No allowed IPs configured. No one can access telemetry", nil)
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusForbidden, ginUnauthorizedIPaddress)
		}

	}

	allowedIPs := strings.Split(allowedIPsStr, ",")
	for i, ip := range allowedIPs {
		allowedIPs[i] = strings.TrimSpace(ip)
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		parsedClientIP := net.ParseIP(clientIP)
		if parsedClientIP == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, ginInvalidClientIPaddress)
			return
		}

		allowed := false
		for _, allowedIP := range allowedIPs {
			if strings.Contains(allowedIP, "/") {
				_, ipNet, err := net.ParseCIDR(allowedIP)
				if err == nil && ipNet.Contains(parsedClientIP) {
					allowed = true
					break
				}
			} else if allowedIP == clientIP {
				allowed = true
				break
			}
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, ginUnauthorizedIPaddress)
			return
		}

		c.Next()
	}
}
