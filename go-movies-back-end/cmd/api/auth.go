package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
)

type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookiePath    string
	CookieName    string
}

type jwtUser struct {
	ID        int      `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Roles     []string `json:"roles"`
}

type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	jwt.RegisteredClaims
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

func (j *Auth) GenerateTokenPair(user *jwtUser) (TokenPairs, error) {
	// Create a token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set the claim
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = j.Audience
	claims["iss"] = j.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"
	// Set the expiry for jWT
	claims["exp"] = time.Now().UTC().Add(j.TokenExpiry).Unix()
	claims["roles"] = user.Roles

	// Create a signed Token
	signedAccessToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// Create a refresh token and set claim
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()
	// Set the expiry for refresh token
	refreshTokenClaims["exp"] = time.Now().UTC().Add(j.RefreshExpiry).Unix()

	// Create signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(j.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	// Create TokenPairs and populate with signed tokens
	var tokenPairs = TokenPairs{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	// Return token pairs
	return tokenPairs, nil
}

func (j *Auth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Value:    refreshToken,
		Expires:  time.Now().Add(j.RefreshExpiry),
		MaxAge:   int(j.RefreshExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   j.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}
}

func (j *Auth) GetExpiredRefreshCookie() *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Domain:   j.CookieDomain,
		HttpOnly: true,
		Secure:   true,
	}
}

func (j *Auth) GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	w.Header().Add("Vary", "Authorization")

	// get auth header
	// Checking if connection is for WebSocket
	var isSocket = false
	var AuthHeader = "Authorization"
	var token string
	if websocket.IsWebSocketUpgrade(r) {
		isSocket = true
		AuthHeader = "Sec-WebSocket-Protocol"
	}

	authHeader := r.Header.Get(AuthHeader)
	fmt.Println(authHeader)

	// sanity check
	if authHeader == "" {
		return "", nil, errors.New("no auth header")
	}
	//split the header on spaces
	if isSocket {
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 3 {
			return "", nil, errors.New("invalid socket auth header")
		}
		// check to see if we have the word Bearer
		if headerParts[0] != "Authorization," {
			return "", nil, errors.New("invalid auth header")
		}
		token = strings.TrimSuffix(headerParts[1], ",")

	} else {
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			return "", nil, errors.New("invalid auth header")
		}

		// check to see if we have the word Bearer
		if headerParts[0] != "Bearer" {
			return "", nil, errors.New("invalid auth header")
		}
		token = headerParts[1]
	}

	fmt.Println(token)

	//
	//declare an empty claims
	claims := &Claims{}

	//Parse the token
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.Secret), nil
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("token expired")
		}
		return "", nil, err
	}

	if claims.Issuer != j.Issuer {
		return "", nil, errors.New("invalid issuer")
	}

	return token, claims, nil
}

// Authorization Validation
func (c *Claims) HasRight(verifyRole string) bool {
	// Check if claims contain admin role
	roleExit := false
	roles := c.Roles
	// Loop through roles and check if admin
	for _, role := range roles {
		if role == verifyRole {
			roleExit = true
			break
		}
	}

	return roleExit
}
