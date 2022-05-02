package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil

	if userType != role {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	return err
}

func MatchUserTypeToUid(c *gin.Context, user_id string) (err error) {
	// Getting current logged in user's credentials
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil
	// Checking with the request variables and userType to be Admin not normal user
	if userType == "USER" && uid != user_id {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	err = CheckUserType(c, userType)
	return err
}
