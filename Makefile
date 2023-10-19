.PHONY: start-compose start-server kill-server start-bot down-compose
start-compose:
	@cd server-application && docker-compose up -d

start-server:
	@cd server-application/cmd/main && go run main.go &

kill-server:
	@kill -9 $$(lsof -t -i:8000) || echo "No process running on port 8000"

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

down-compose:
	@cd server-application && docker-compose down