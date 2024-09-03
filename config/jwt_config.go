package config

import "time"

var JWTSecretkey = []byte("your-secret-key")

const JTWExpiredTime = 5 * time.Minute
