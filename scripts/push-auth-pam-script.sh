#!/bin/bash

req_id=$(uuidgen)
info="$PAM_USER@$(hostname -f)"

# Enforce the approval only for users belonging to group `ci`
if ! groups "$PAM_USER" | grep -wq "ci"; then
  exit 0
fi

# Basic auth credentials are used in curl requests
if [ "$PAM_TYPE" != "auth" ]; then
  if curl -sS  -u 'client:password' 'http://localhost:9090/request/submit'  -d "{\"request_id\":\"$req_id\", \"info\":\"$info\"}"; then
    for (( i = 0; i < 30; i++ )); do
      sleep 1
      res=$(curl -sS  -u 'client:password' 'http://localhost:9090/request/status'  -d "{\"request_id\":\"$req_id\"}")
      if echo "$res" | grep '\"allowed\"'; then
	echo "allowed!"
        exit 0
      fi
      if echo "$res" | grep '\"denied\"'; then
        echo "denied!"
        exit 2
      fi
    done
  fi
  echo "error!"
  exit 1
fi
