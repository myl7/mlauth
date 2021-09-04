package conf

var SiteHost = env("SITE_HOST", nil)

var DbDriver = env("DB_DRIVER", strPtr("postgres"))
var DbSource = env("DB_SOURCE", nil)

var KvAddr = env("KV_ADDR", strPtr("localhost:6379"))
var KvPassword = env("KV_PASSWORD", strPtr(""))
var KvDb = envInt("KV_DB", intPtr(0))

var SecretKey = []byte(env("SECRET_KEY", strPtr("debug-secret-key")))

var AccessTokenAge = envInt("ACCESS_TOKEN_AGE", intPtr(1*60*60))
var UpdateTokenAge = envInt("UPDATE_TOKEN_AGE", intPtr(7*24*60*60))

var SmtpHost = env("SMTP_HOST", nil)
var SmtpPort = envInt("SMTP_PORT", intPtr(587))
var SmtpSender = env("SMTP_SENDER", nil)
var SmtpUsername = env("SMTP_USERNAME", nil)
var SmtpPassword = env("SMTP_PASSWORD", nil)

var UserActiveEmailAge = envInt("USER_ACTIVE_EMAIL_AGE", intPtr(1*24*60*60))
