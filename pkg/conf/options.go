package conf

var DbDriver = env("DB_DRIVER", strPtr("postgres"))
var DbSource = env("DB_SOURCE", nil)

var SecretKey = []byte(env("SECRET_KEY", strPtr("debug-secret-key")))

var AccessTokenAge = envInt("ACCESS_TOKEN_AGE", intPtr(1*60*60))
var UpdateTokenAge = envInt("UPDATE_TOKEN_AGE", intPtr(7*24*60*60))
