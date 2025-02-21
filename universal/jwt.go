package universal


import (
  "time"

	"github.com/golang-jwt/jwt/v5"
)


var jwtKey = []byte(secretKey)

type Claims struct {
	Hash string `json:"hash"`
	jwt.RegisteredClaims
}

type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"` // Unix timestamp
}

func GenerateJWT(hash string, duration time.Duration) *TokenResponse {
	expirationTime := time.Now().Add(duration).Unix()

	claims := &Claims{
		Hash: hash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
      Issuer: "termo-chat",
      Subject: "user-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  signedToken, _ := token.SignedString(jwtKey)

  response := TokenResponse {
    Token: signedToken,
    ExpiresAt: expirationTime,
  }

  return &response
}
