package middlewares

import (
	V1Domains "Shorty.Server.Go.Management/internal/business/domains/v1"
	"Shorty.Server.Go.Management/internal/config"
	"Shorty.Server.Go.Management/internal/constants"
	"Shorty.Server.Go.Management/pkg/helpers"
	"Shorty.Server.Go.Management/pkg/logger"
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/mitchellh/mapstructure"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func KeycloakAuthMiddleware(rdc *redis.Client) func(roles ...string) fiber.Handler {
	return func(roles ...string) fiber.Handler {
		return func(ctx *fiber.Ctx) error {
			authHeader := ctx.Get(constants.AuthHeader)
			if authHeader == "" {

				if err := ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "No token provided"}); err != nil {
					return err
				}
				return nil
			}

			if headerArr := strings.Split(authHeader, " "); len(headerArr) != 2 || headerArr[0] != "Bearer" {
				logger.Error("Invalid token format", nil)

				if err := ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"}); err != nil {
					return err
				}
				return nil
			}

			token, err := verifyToken(ctx.Context(), rdc, strings.Split(authHeader, " ")[1])

			if err != nil {
				logger.Error("cannot verify token", logrus.Fields{"err": err.Error()})

				if err := ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()}); err != nil {
					return err
				}
				return nil
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !(ok && token.Valid) {
				logger.Error("cannot get claims", logrus.Fields{"err": err.Error()})
				if err := ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"}); err != nil {
					return err
				}
				return nil
			}

			if claims["sub"] == "" || claims["sub"] == nil {
				if err := ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"}); err != nil {
					return err
				}
				return nil
			}

			var userRoles []string
			if resourceAccess, ok := claims["resource_access"].(map[string]interface{}); ok {
				if authClient, ok := resourceAccess[config.AppConfig.AUTHClient].(map[string]interface{}); ok {
					if err := mapstructure.Decode(authClient["roles"], &userRoles); err != nil {
						logger.Error("cannot get user roles", logrus.Fields{"err": err.Error()})
						userRoles = []string{}
					}
				}
			}

			userEmail, ok := claims["email"].(string)
			if !ok {
				userEmail = ""
			}
			ctx.Locals(constants.UserDetails, V1Domains.UserDetails{
				Roles:      userRoles,
				UserId:     claims["sub"].(string),
				Email:      userEmail,
				Username:   claims["preferred_username"].(string),
				Name:       claims["name"].(string),
				FamilyName: claims["family_name"].(string),
			})

			if len(roles) == 0 {
				if err := ctx.Next(); err != nil {
					return err
				}
				return nil
			}

			if !isUserHaveRoles(roles, userRoles) {
				ctx.Status(http.StatusForbidden)
			}

			return ctx.Next()
		}
	}
}

func isUserHaveRoles(roles []string, userRoles []string) bool {
	for _, role := range roles {
		if helpers.IsArrayContains(userRoles, role) {
			return true
		}
	}
	return false
}

func verifyToken(ctx context.Context, r *redis.Client, tokenString string) (token *jwt.Token, err error) {

	//if err := verifyTokenSession(tokenString); err != nil {
	//	return nil, err
	//}

	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (rawKey interface{}, err error) {

		result, err := r.Get(ctx, constants.JwkKey).Result()
		if err == nil {
			logger.Info("Jwk get from cache", nil)
			resultKey := deserializePublicKey(result)
			return &resultKey, nil
		}

		set, err := jwk.Fetch(ctx, config.AppConfig.AUTHJwkPublicUri)
		if err != nil {
			return nil, err
		}
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			err = fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, err
		}

		keyID, ok := token.Header["kid"].(string)
		if !ok {
			err = fmt.Errorf("expecting JWT header to have string 'kid'")
			return nil, err
		}

		key, found := set.LookupKeyID(keyID)
		if !found {
			err = fmt.Errorf("unable to find key")
			return nil, err
		}

		err = key.Raw(&rawKey)
		if err != nil {
			return nil, fmt.Errorf("invalid token")
		}

		serialRawKey := serializePublicKey(*rawKey.(*rsa.PublicKey))
		r.Set(ctx, constants.JwkKey, serialRawKey, time.Duration(config.AppConfig.AUTHRefreshJwkTimeout)*time.Second)
		logger.Info("Jwk get from sso and saved to cache", nil)
		return rawKey, err
	})

	if err != nil {
		return nil, err
	}

	return token, nil

}

func serializePublicKey(key rsa.PublicKey) string {
	serialized := map[string]string{
		"N": key.N.String(),
		"E": fmt.Sprintf("%d", key.E),
	}
	serializedKey, err := json.Marshal(serialized)
	if err != nil {
		panic(err)
	}
	return string(serializedKey)
}

// DeserializePublicKey deserializes a string into an RSA public key
func deserializePublicKey(serializedKey string) rsa.PublicKey {
	var serialized map[string]string
	if err := json.Unmarshal([]byte(serializedKey), &serialized); err != nil {
		panic(err)
	}

	N := new(big.Int)
	N.SetString(serialized["N"], 10)

	E, err := strconv.ParseInt(serialized["E"], 10, 64)
	if err != nil {
		panic(err)
	}

	return rsa.PublicKey{
		N: N,
		E: int(E),
	}
}

// uncomment this if you want to verify user session in every request
//func verifyTokenSession(tokenString string) (err error) {
//	req, err := http.NewRequest("GET", config.AppConfig.AUTHUserInfoEndpoint, nil)
//	if err != nil {
//		return err
//	}
//
//	req.Header.Set("Authorization", "Bearer "+tokenString)
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		return err
//	}
//
//	if resp.StatusCode != 200 {
//		return fmt.Errorf("user session is die")
//	}
//
//	logger.Info("token session is ok", nil)
//
//	return nil
//}
