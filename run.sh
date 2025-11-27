#!/bin/bash
(cd backend && go run cmd/worldle/main.go) &
BACKEND_PID=$!

is_port_free() {
  local port=$1
  if lsof -i :$port >/dev/null 2>&1; then
    return 1  # Not free
  else
    return 0  # Free
  fi
}

cleanup() {
  echo "Stopping backend (PID $BACKEND_PID)..."
  kill $BACKEND_PID
  wait $BACKEND_PID 2>/dev/null

  sleep 1

  is_port_free 8080
  BACKEND_FREE=$?
  is_port_free 3000
  FRONTEND_FREE=$?

  if [[ $BACKEND_FREE -eq 0 && $FRONTEND_FREE -eq 0 ]]; then
    echo "Both ports 8080 (backend) and 3000 (frontend) are free."
  else
    [[ $BACKEND_FREE -ne 0 ]] && echo "Port 8080 is still in use."
    [[ $FRONTEND_FREE -ne 0 ]] && echo "Port 3000 is still in use."
  fi

  exit
}

# Trap SIGINT (Ctrl+C) and SIGTERM
trap cleanup SIGINT SIGTERM

(cd frontend && npm start)

cleanup