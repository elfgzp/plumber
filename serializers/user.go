package serializers

import (
	"github.com/elfgzp/plumber/models"
	"time"
)

type UserSerializer struct {
	Nickname          string    `json:"nickname"`
	Email             string    `json:"email"`
	MobileCountryCode string    `json:"mobile_country_code"`
	Mobile            string    `json:"mobile"`
	Slug              string    `json:"slug"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func SerializeUser(u *models.User) UserSerializer {
	us := UserSerializer{
		Nickname:          u.Nickname,
		Email:             u.Email,
		MobileCountryCode: u.MobileCountryCode,
		Mobile:            u.Mobile,
		Slug:              u.Slug,
		CreatedAt:         u.CreatedAt,
		UpdatedAt:         u.UpdatedAt,
	}

	return us
}
