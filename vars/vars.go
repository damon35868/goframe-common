package vars

var (
	ClientJwtKey = "jwt.client"
	AdminJwtKey  = "jwt.admin"
)

func SetJwtKey(keys ...string) {
	if len(keys) > 0 {
		ClientJwtKey = keys[0]
	}
	if len(keys) > 1 {
		AdminJwtKey = keys[1]
	}
}
