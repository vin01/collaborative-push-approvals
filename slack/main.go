package slack

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func SendApprovalRequest(request_id string, info string) (err error) {
	message := []byte(fmt.Sprintf(`{
        "text": "Login approval request",
        "attachments": [
            {
                "text": "*%[1]s*",
                "color": "#3AA3E3",
                "attachment_type": "default",
                "actions": [
                    {
                        "url": "http://localhost:8080/?action=allow&request_id=%[2]s",
                        "type": "button",
                        "text" : ":white_check_mark: Allow"
                    },
                    {
                        "url": "http://localhost:8080/?action=deny&request_id=%[2]s",
                        "type": "button",
                        "text" : ":x: Deny"
                    },
                ]
            }
        ]
}`, info, request_id))

	url, ok := os.LookupEnv("SLACK_URL")
	if !ok {
		return errors.New("SLACK_URL not found, skipping slack notification")
	}
	resp, err := http.Post(url, "application/json",
		bytes.NewBuffer(message))

	if resp.StatusCode != http.StatusOK {
		errbody, _ := ioutil.ReadAll(resp.Body)
		return errors.New("Received erroneous response: " + string(errbody))
	}

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func SendRequestStatus(request_id string, info string) (err error) {
	message := []byte(fmt.Sprintf(`{
        "text": "Request for %[1]s %[2]s",}`, request_id, info))

	url, ok := os.LookupEnv("SLACK_URL")
	if !ok {
		return errors.New("SLACK_URL not found, skipping slack notification")
	}
	resp, err := http.Post(url, "application/json",
		bytes.NewBuffer(message))

	if resp.StatusCode != http.StatusOK {
		errbody, _ := ioutil.ReadAll(resp.Body)
		return errors.New("Received erroneous response: " + string(errbody))
	}

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}
