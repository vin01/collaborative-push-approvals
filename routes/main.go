package routes

import (
	"database/sql"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

const userkey = "user_id"

var secret = []byte(os.Getenv("SESSION_SECRET"))

var db *sql.DB
var err error

func Run() {
	r := SetupRouter()
	r.Run(":8080")
}

func SetupRouter() *gin.Engine {
	db, err = sql.Open("sqlite3", "./collaborative-push-approvals.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
        create table if not exists requests (request_id TEXT primary key not null, request_info TEXT not null, created_at DATETIME not null, status TEXT not null);
        `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil
	}

	r := gin.Default()
	r.Static("/static", "./static/")
	r.LoadHTMLGlob("templates/*")

	clients, err := db.Query("SELECT username, password FROM clients")
	if err != nil {
		return nil
	}
	defer clients.Close()

	client_accounts := gin.Accounts{}

	var (
		username string
		password string
	)

	for clients.Next() {
		if err := clients.Scan(&username, &password); err != nil {
			log.Println(err)
		}
		client_accounts[username] = password
	}

	if len(secret) == 0 {
		log.Fatal("Found empty session secret, exiting!\nPlease set a secure secret using `SESSION_SECRET` environment variable.")
	}
	r.Use(sessions.Sessions("login", cookie.NewStore(secret)))
	r.POST("/login", login)
	//r.GET("/login", func(c *gin.Context) {
	//      c.Redirect(http.StatusFound, "/")
	//})
	r.GET("/logout", logout)

	authorized_approvers := r.Group("/")
	authorized_approvers.Use(AuthRequired)
	authorized_approvers.GET("/", approverHandler)
	authorized_clients := r.Group("/request", gin.BasicAuth(client_accounts))
	authorized_clients.POST("/submit", approvalRequestHandler)
	authorized_clients.POST("/status", approvalRequestStatusHandler)
	return r
}
