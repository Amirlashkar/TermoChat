package universal

import (
	"fmt"
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

func parseToken(tokenString string) (*jwt.Token, error) {
    return jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtKey, nil
    })
}

// Checks a JWT token validity due to its expiration date & signature ;
// Also gives back the related user hash
func IsTokenValid(tokenString string) (bool, error, *string) {
    token, err := parseToken(tokenString)
    if err != nil {
        return false, fmt.Errorf("token parsing error: %v", err), nil
    }

    hash := token.Claims.(*Claims).Hash

    if token.Valid {
        return true, nil, &hash // we used pointer here to handle nil values returning
    }

    return false, fmt.Errorf("token is invalid"), nil
}

func GetUHash(tokenString string) string {
    // No error cause this function  is called when the validations
    // are already checked by top function
    token, _ := parseToken(tokenString)
    return token.Claims.(*Claims).Hash
}

