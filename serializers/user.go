package serializers

import "github.com/elfgzp/plumber/models"

type UserSerializer struct {
	ModelSerializer
	Nickname          string `json:"nickname"`
	Email             string `json:"email"`
	MobileCountryCode string `json:"mobile_country_code"`
	Mobile            string `json:"mobile"`
}

func SerializeUser(u *models.User) UserSerializer {
	us := UserSerializer{
		Nickname:          u.Nickname,
		Email:             u.Email,
		MobileCountryCode: u.MobileCountryCode,
		Mobile:            u.Mobile,
	}

	us.serializeBaseField(u.Slug, u.CreatedAt, u.UpdatedAt)
	return us
}
