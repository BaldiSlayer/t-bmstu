package handler

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Baldislayer/t-bmstu/pkg/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	clientID       = ""
	clientSecret   = ""
	redirectURL    = "http://localhost:8080/auth/callback"
	sessionName    = "mysession"
	accessTokenKey = "access_token"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegistrationForm struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	LastName  string `json:"lastname"`
	FirstName string `json:"firstname"`
	Group     string `json:"group"`
	Email     string `json:"email"`
}

var (
	oauthConf = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
	// TODO to config
	jwtSecret = "your-secret-key"
)

func isAuth(c *gin.Context) bool {
	cookie, err := c.Cookie("token")
	if err != nil {
		return false
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return false
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	return true
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		role, err := database.GetUserRole(username)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/auth/login")
			c.Abort()
			return
		}

		c.Set("username", username)
		c.Set("role", role)

		c.Next()
	}
}

func generateToken(username string, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Генерация подписи токена
	secret := jwtSecret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *Handler) signIn(c *gin.Context) {
	// проверить авторизован ли уже пользователь
	if isAuth(c) {
		c.Redirect(http.StatusSeeOther, "/view/home")
		c.Abort()
	}

	requestMethod := c.Request.Method
	switch requestMethod {
	case "GET":
		{
			c.HTML(http.StatusOK, "login.tmpl", gin.H{})
		}
	case "POST":
		{
			var form LoginForm
			if err := c.BindJSON(&form); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
				return
			}

			valid, role := checkLoginAndPassword(form.Username, form.Password)
			if valid {
				token, err := generateToken(form.Username, role)

				if err != nil {
					c.JSON(http.StatusUnauthorized, "Ошибка генерации токена jwt")
				}

				cookie := http.Cookie{
					Name:     "token",
					Path:     "/",
					Value:    token,
					Expires:  time.Now().Add(time.Hour * 24),
					HttpOnly: true,
				}
				http.SetCookie(c.Writer, &cookie)

				c.JSON(http.StatusOK, gin.H{"message": "Успешный вход"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные логин или пароль"})
			}
		}
	default:
		{
			c.JSON(http.StatusBadRequest, "No such router for this method")
		}
	}
}

func checkLoginAndPassword(username, password string) (bool, string) {
	str := password
	hash := md5.Sum([]byte(str)) // Хэширование в MD5

	hashString := hex.EncodeToString(hash[:])

	exist, err := database.AuthenticateUser(username, hashString)

	if err != nil {
		fmt.Println("error")
		return false, ""
	}

	if !exist {
		return false, ""
	}

	role, err := database.GetUserRole(username)

	if err != nil {
		fmt.Println("error")
		return false, ""
	}

	return true, role
}

func (h *Handler) signUp(c *gin.Context) {
	requestMethod := c.Request.Method
	switch requestMethod {
	case "GET":
		{
			c.HTML(http.StatusOK, "registration.tmpl", gin.H{})
		}
	case "POST":
		{
			var form RegistrationForm
			if err := c.BindJSON(&form); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
				return
			}

			// fmt.Println(123)
			exist, err := database.CheckIfUserExists(form.Username)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
			}

			if exist {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с таким никнеймом уже есть"})
				return
			}

			hash := md5.Sum([]byte(form.Password))

			hashPassword := hex.EncodeToString(hash[:])

			database.CreateUser(database.User{
				Username:     form.Username,
				PasswordHash: hashPassword,
				LastName:     form.LastName,
				FirstName:    form.FirstName,
				Email:        form.Email,
				Group:        form.Group,
				Role:         "student",
				SolvedTasks:  []string{},
				Groups:       []json.RawMessage{},
			})

			c.Redirect(302, "/auth/login")
		}
	default:
		{
			c.JSON(http.StatusBadRequest, "No such router for this method")
		}
	}
}

func (h *Handler) githubSignUp(c *gin.Context) {
	url := oauthConf.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) githubCallback(c *gin.Context) {
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

func oldAuthMiddleware() gin.HandlerFunc {
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
