package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		// Abort the request with the appropriate error code
		c.HTML(http.StatusOK, "login.tmpl", nil)
		c.Abort()
	}
	c.Next()
}

func login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")
	var (
		user_id       string
		password_hash []byte
	)

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.HTML(http.StatusBadRequest, "login.tmpl", gin.H{"error": "Input fields must not be empty"})
		return
	}

	err = db.QueryRow("select id, password_hash from approvers where username = ?", username).Scan(&user_id, &password_hash)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "login.tmpl", gin.H{"error": "Error validatng password"})
		return
	}

	err = bcrypt.CompareHashAndPassword(password_hash, []byte(password))
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{"error": "Authentication failed"})
		return
	}

	session.Set(userkey, user_id)
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "login.tmpl", gin.H{"error": "Failed to save session"})
		return
	}
	c.Redirect(http.StatusFound, "/")
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.HTML(http.StatusBadRequest, "login.tmpl", gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "login.tmpl", gin.H{"error": "Failed to save session"})
		return
	}
	c.HTML(http.StatusOK, "login.tmpl", nil)
}
