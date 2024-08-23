package routes

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"collaborative-push-approvals/slack"
)

func approverHandler(c *gin.Context) {
	var (
		created_at    string
		info          string
		request_id    string
		rows_affected int64
		status        string
	)

	data := gin.H{}
	action := c.Query("action")
	action_id := c.Query("request_id")

	if action == "allow" {
		stmt, err := db.Prepare("update requests set status=\"allowed\" where request_id=? AND CAST(strftime(\"%s\", CURRENT_TIMESTAMP) as integer) - CAST(strftime(\"%s\", created_at) as integer) < 30 AND status=\"pending\"")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		res, err := stmt.Exec(action_id)
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.tmpl", "Error allowing request_id "+action_id)
			c.Abort()
		}
		rows_affected, err = res.RowsAffected()
		log.Println(rows_affected)
		log.Println("----")
		if err != nil {
			log.Println("Unable to determine if rows were changed")
			c.Abort()
		}
	}

	if action == "deny" {
		stmt, err := db.Prepare("update requests set status=\"denied\" where request_id=? AND CAST(strftime(\"%s\", CURRENT_TIMESTAMP) as integer) - CAST(strftime(\"%s\", created_at) as integer) < 30 AND status=\"pending\"")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		res, err := stmt.Exec(action_id)
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.tmpl", "Error denying request id "+action_id)
			c.Abort()
		}
		rows_affected, err = res.RowsAffected()
		if err != nil {
			log.Println("Unable to determine if rows were changed")
			c.Abort()
		}
	}

	rows, err := db.Query("select request_id, CAST(strftime('%s', CURRENT_TIMESTAMP) as integer) - CAST(strftime('%s', created_at) as integer) as time_diff, request_info, status from requests where time_diff < 30 order by time_diff;")
	if err != nil {
		c.HTML(http.StatusOK, "error.tmpl", "Error querying database for pending requests.")
		c.Abort()
	}

	for rows.Next() {
		if err := rows.Scan(&request_id, &created_at, &info, &status); err != nil {
			c.HTML(http.StatusOK, "error.tmpl", "Error parsing database entries.")
			c.Abort()
		}
		data[string(len(data)+1)] = []interface{}{status, created_at, info, request_id}
		if action_id == request_id && rows_affected > 0 {
			if action == "allow" {
				err = slack.SendRequestStatus(info, "allowed :white_check_mark:")
			}
			if action == "deny" {
				err = slack.SendRequestStatus(info, "denied :x:")
			}
			if err != nil {
				log.Println("Unable to send slack message: ", err)
			}
		}
	}

	c.HTML(http.StatusOK, "app.tmpl", data)
}
