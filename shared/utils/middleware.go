package utils

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)


type AuditLog struct {
	ID         int             `gorm:"primaryKey;autoIncrement" json:"id"`        
	RequestID  string          `gorm:"not null" json:"request_id"`  
	Table  string              `gorm:"type:varchar(50);not null" json:"table_name"`
	RecordID   string          `gorm:"not null" json:"record_id"`   
	Action     string          `gorm:"type:action;not null" json:"action"`     
	OldData    json.RawMessage `gorm:"type:jsonb" json:"old_data"`             
	NewData    json.RawMessage `gorm:"type:jsonb" json:"new_data"`             
	CreatedAt  time.Time       `gorm:"not null;default:now()" json:"created_at"` 
	CreatedBy  string          `gorm:"not null" json:"created_by"`   
	IPAddress  string          `gorm:"type:varchar(45);not null" json:"ip_address"`
}

// TableName menentukan nama tabel untuk model AuditLog.
func (AuditLog) TableName() string {
	return "audit_logs"
}


// JWTAuthMiddleware is used to check the validity of the JWT token in the Authorization header.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// Token should be in the format: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token and extract user information
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		
		// Extract the claims for user data
		userID, _ := claims["userId"].(int)	
		userName, _ := claims["userName"].(string)
		role , _ := claims["role"].(string)
		// Set userID in context for further use in the handler	
		c.Set("userId", userID)
		c.Set("userName", userName)
		c.Set("role", role)
		c.Next()
	}
	
}

func AuditLogMiddleware(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := uuid.New().String()
        c.Set("request_id", requestID)
        c.Set("ip_address", c.ClientIP())
        c.Header("X-Request-ID", requestID)

        c.Next()

        if c.Request.Method == "POST" && c.Request.URL.Path == "/v1/payroll_periods/create" && c.Writer.Status() == http.StatusCreated {
            userID, exists := c.Get("userId")
            if !exists {
                return
            }
            recordID, exists := c.Get("record_id")
            if !exists {
                return
            }
            newData, _ := json.Marshal(c.MustGet("response_data"))
            auditLog := AuditLog{
                RequestID: requestID,
                Table:    "payroll_periods",
                RecordID:  recordID.(string),
                Action:    "CREATE",
                NewData:   newData,
                CreatedAt: time.Now(),
                CreatedBy: userID.(string),
                IPAddress: c.ClientIP(),
            }
            db.Create(&auditLog)
        }
    }
}