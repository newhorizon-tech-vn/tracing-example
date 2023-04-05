package authorize

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/newhorizon-tech-vn/tracing-example/pkg/util"
	"github.com/newhorizon-tech-vn/tracing-example/services"
)

const KeyUserId = "UserId"
const KeyRoleId = "RoleId"

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := ParseUserIdFromToken(c)
		if err != nil {
			c.Abort()
			return
		}
		user, err := (&services.UserService{UserId: userId}).GetUser(c.Request.Context())
		if err != nil {
			c.Abort()
			return
		}

		c.Set(KeyUserId, userId)
		c.Set(KeyRoleId, user.GetRoleId())
		c.Next()
	}
}

func ParseUserIdFromToken(c *gin.Context) (int, error) {
	splitToken := strings.Split(c.GetHeader("Authorization"), "Bearer")
	if len(splitToken) < 2 {
		return 0, fmt.Errorf("Token invalid")
	}

	return util.StringToInt(strings.TrimSpace(splitToken[1]))
}
