package utils

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type AuditLog struct {
    ID        int             `gorm:"primaryKey;autoIncrement" json:"id"`
    RequestID int             `gorm:"autoIncrement;not null" json:"request_id"`
    Table     string          `gorm:"type:varchar(50);not null" json:"table_name"`
    RecordID  string          `gorm:"type:varchar(50);not null" json:"record_id"`
    Action    string          `gorm:"type:action;not null" json:"action"`
    OldData   json.RawMessage `gorm:"type:jsonb" json:"old_data"`
    NewData   json.RawMessage `gorm:"type:jsonb" json:"new_data"`
    CreatedAt time.Time       `gorm:"not null;default:now()" json:"created_at"`
    CreatedBy string          `gorm:"type:varchar(50);not null" json:"created_by"`
    IPAddress string          `gorm:"type:varchar(45);not null" json:"ip_address"`
}

func (AuditLog) TableName() string {
    return "audit_logs"
}

func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        userid, _ := claims["userid"].(string)
        userName, _ := claims["userName"].(string)
        role, _ := claims["role"].(string)

        c.Set("userid", userid)
        c.Set("userName", userName)
        c.Set("role", role)
        c.Next()
    }
}

func AuditLogMiddleware(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := int(time.Now().UnixNano())
        c.Set("request_id", strconv.Itoa(requestID))
        c.Set("ip_address", c.ClientIP())
        c.Header("X-Request-ID", strconv.Itoa(requestID))

        c.Next()

        if c.Request.Method == "POST" && c.Request.URL.Path == "/v1/payroll_periods/create" && c.Writer.Status() == http.StatusCreated {
            userID, exists := c.Get("useri")
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
                Table:     "payroll_periods",
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