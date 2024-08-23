package routes

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"collaborative-push-approvals/slack"
	"strings"
)

type RequestBody struct {
	Request_id string `json:"request_id"`
	Info       string `json:"info"`
}

func approvalRequestHandler(c *gin.Context) {

	var request_body RequestBody
	if err := c.BindJSON(&request_body); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to parse json payload",
		})
		c.Abort()
	}
	if request_body.Request_id != "" && request_body.Info != "" {
		stmt, err := db.Prepare("insert into requests values (?, ?, CURRENT_TIMESTAMP, \"pending\")")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Error preparing query to insert requests",
			})
			c.Abort()
		}
		defer stmt.Close()
		_, err = stmt.Exec(request_body.Request_id, request_body.Info)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Error inserting into requests",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success": "Submitted " + request_body.Request_id,
			})

			err = slack.SendApprovalRequest(request_body.Request_id, request_body.Info)
			if err != nil {
				log.Println("Unable to send slack message: ", err)
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "info and request_id should not be empty in payload",
		})
	}
}

func approvalRequestStatusHandler(c *gin.Context) {
	var request_body RequestBody
	var status string
	if err := c.BindJSON(&request_body); err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error.tmpl", "Unable to parse json payload")
		c.Abort()
	}
	if strings.Trim(request_body.Request_id, " ") != "" {
		err = db.QueryRow("select status from requests where request_id = ?", request_body.Request_id).Scan(&status)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusBadRequest, "error.tmpl", "Error querying the status of request")
			c.Abort()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": status,
			})
		}
	}
}
