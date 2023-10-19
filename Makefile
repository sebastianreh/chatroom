.PHONY: start-compose start-server kill-server start-bot chatroom-frontend start-frontend down-compose
start-compose:
	@cd server-application && docker-compose up -d

start-server:
	@cd server-application/cmd/main && go run main.go &

create-room:
	@ID=$$(curl --location --request POST 'http://localhost:8000/chatroom/room' \
	--header 'Content-Type: application/json' \
	--data-raw '{"name": "chat-room","is_active": true}' | jq -r '.id'); \
	if [ -n "$$ID" ] && [ "$$ID" != "null" ]; then \
		echo "ID captured: $$ID"; \
		echo "export ROOM_ID=$$ID" > room_id.env; \
	else \
		echo "error"; \
	fi

start-bot:
	@. ./room_id.env && cd bots/stocks && go run main.go -room_id=$$ROOM_ID &

install-frontend-dependencies:
	@(cd chatroom-frontend && npm install)

start-frontend:
	@(cd chatroom-frontend && npm run start) &

down-compose:
	@cd server-application && docker-compose down

kill-project:
	@kill -9 $$(lsof -t -i:8000) || echo "No process running on port 8000"
	@kill -9 $$(lsof -t -i:3000) || echo "No process running on port 3000"

test-server:
	@(cd server-application && go test ./...)