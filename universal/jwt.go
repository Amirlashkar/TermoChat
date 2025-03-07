package universal


import (
    "time"
    "fmt"

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

func IsTokenValid(tokenString string) (bool, error) {
    // Parser deals with expired, unsigned jwt tokens itself
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtKey, nil
    })

    if err != nil {
        return false, fmt.Errorf("token parsing error: %v", err)
    }

    if token.Valid {
        return true, nil
    }

    return false, fmt.Errorf("token is invalid")
}
