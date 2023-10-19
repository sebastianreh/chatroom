# Chatroom Project

## Overview

The Chatroom Project is a browser-based chat application designed primarily with Go. The focus is on building a robust
back-end system with capabilities for real-time chatting and implementing decoupled bots, using a stock quote bot as an
example. The architecture is based on the Domain Driven Design pattern.

## Features

- **User Authentication**: Allows registered users to securely log in.

- **Real-time Chat**: Users can converse in a chatroom with real-time messaging capabilities.

- **Stock Bot & Command**: When active for a channel, users can fetch real-time stock quotes by typing messages in the
  format `/stock=stock_code`, like `/stock=aapl.us` for Apple Inc. The bot retrieves this information from an external
  API, parses the CSV data, and sends a message to the chatroom such as "APPL.US quote is $93.42 per share".

- **Message Ordering and Limit**: Chat messages are displayed in order of their timestamps, showcasing only the most
  recent 50 messages for clarity and performance.

- **Multiple Chatrooms**: Users have the flexibility to join multiple chatrooms.

- **Bot Exception Handling**: The bot is equipped to manage unrecognized messages and exceptions gracefully, ensuring
  seamless user experience.

---

## Technologies Used

- **MongoDB**
  MongoDB is employed as our primary database to store persistent users and rooms data.

- **Redis**
  We utilize Redis for real-time functionalities. It's responsible for managing chat sessions and user actions.

- **Kafka**
  Kafka handles our bot messaging system between the bot and the server. 

---

## Chatroom Server Setup Guide

This README provides a step-by-step guide to set up and run the chatroom project.

### Prerequisites

Before starting, ensure you have the following installed:

- Docker and Docker Compose
- Go (v1.19) (for running the server and bot)
- npm (for frontend dependencies)
- jq (for JSON processing)

### Step-by-Step Guide

1. **Start Docker Services**:
   `make start-compose`

   This step will use Docker Compose to start necessary services defined in the `docker-compose.yml` file in the
   server-application folder.

2. **Start the Server**:
   `make start-server`

   This command starts the Go server.

3. **Create a Chat Room**:
   `make create-room`

   This command sends a POST request to the server to create a new chatroom. It captures the ID of the newly created
   room and stores it in `room_id.env` to use it when starting the bot in next step.

4. **Start the Bot**:
   `make start-bot`

   This command will start the chat bot. Before starting, it fetches the room ID from the `room_id.env` file.

5. **Install Frontend Dependencies**:
   `make install-frontend-dependenciest`

   If you haven't installed frontend dependencies or if there are new dependencies added, run this command.

6. **Start the Frontend**:
   `make start-frontend`

   This command will start the frontend server. It uses npm to run the frontend.

7. **Kill the Server and Frontend (Optional)**:
   `make kill-project`

   If for any reason you wish to force stop the Go server and the frontend, use the following command. This will
   terminate any process running on port 8000 and 3000.

8. **Shut Down the Services (Optional)**:
   `make down-compose`

   If you wish to shut down all services started using Docker Compose, use the following command:

---

### Test

In order to run the server tests, use this **command**: `make test-server`