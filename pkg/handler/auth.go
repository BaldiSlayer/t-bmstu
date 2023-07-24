package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	clientID       = "e5f257f8b08516d4da8f"
	clientSecret   = "26a51b3240af0d236f6e2c004842b9b2a85095c8"
	redirectURL    = "http://localhost:8080/auth/callback"
	sessionName    = "mysession"
	accessTokenKey = "access_token"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
)

func (h *Handler) signUp(c *gin.Context) {
	url := oauthConf.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) callback(c *gin.Context) {
	code := c.Query("code")

	token, err := oauthConf.Exchange(c, code)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Token exchange error: %s", err.Error()))
		return
	}

	client := oauthConf.Client(c, token)
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("API call error: %s", err.Error()))
		return
	}
	defer resp.Body.Close()

	// Получение токена из сессии
	session := sessions.Default(c)
	session.Set(accessTokenKey, token.AccessToken)
	session.Save()

	redirectURL := session.Get("redirect_url")
	if redirectURL != nil {
		c.Redirect(http.StatusFound, redirectURL.(string))
	} else {
		c.Redirect(http.StatusFound, "/constests")
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлечение токена из сессии
		session := sessions.Default(c)
		token := session.Get(accessTokenKey)

		if token == nil {
			session.Set("redirect_url", c.Request.RequestURI)
			session.Save()
			c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
			return
		}

		// Использование токена для аутентификации запроса
		client := oauthConf.Client(c, &oauth2.Token{AccessToken: token.(string)})

		// Запрос информации о пользователе
		resp, err := client.Get("https://api.github.com/user")
		if err != nil {
			session.Set("redirect_url", c.Request.RequestURI)
			session.Save()
			c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
			return
		}
		defer resp.Body.Close()

		var user map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			fmt.Println("Failed to parse user JSON:", err)
			return
		}

		// Извлечение email, имени аккаунта и ссылки на аватарку из полученных данных
		//email := user["email"].(string)
		// (user["id"].(float64))
		login := user["login"].(string)
		avatarURL := user["avatar_url"].(string)

		c.Set("login", login)
		c.Set("avatar_url", avatarURL)

		client = oauthConf.Client(c, &oauth2.Token{AccessToken: token.(string)})
		resp, err = client.Get("https://api.github.com/user/emails")
		if err != nil {
			session.Set("redirect_url", c.Request.RequestURI)
			session.Save()
			c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
			return
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body.(io.Reader))
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read response body")
			return
		}

		email := string(bodyBytes)
		var emails []map[string]interface{}
		err = json.Unmarshal([]byte(email), &emails)
		if err != nil {
			fmt.Println("Failed to parse JSON:", err)
			return
		}

		c.Set("email", emails[0]["email"])

		c.Next()
	}
}
