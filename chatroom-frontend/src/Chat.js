import React, {useEffect, useState, useRef, useContext} from 'react';
import {useUser} from "./UserContext";
import {useRoom} from "./RoomContext";
import {useNavigate, usePrompt} from 'react-router-dom';

function ChatInput({onSend}) {
    const [message, setMessage] = useState('');
    const {user} = useUser();

    const handleSend = () => {
        if (message.trim() !== '') {
            onSend(message, user.username, user.user_id, new Date().toISOString());
            setMessage('');
        }
    };

    const handleKeyDown = (e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            handleSend();
        }
    };

    return (
        <div className="chat-input">
            <input
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="Type a message..."
            />
            <button onClick={handleSend}>Send</button>
        </div>
    );
}

function ChatBox({messages}) {
    const lastMessageRef = useRef(null);

    useEffect(() => {
        if (lastMessageRef.current) {
            lastMessageRef.current.scrollIntoView({behavior: 'smooth'});
        }
    }, [messages]);

    return (
        <div className="chat-box">
            {messages.map((msg, index) => (
                <div key={msg.id} className={`message ${msg.type}`}>
                    {msg.type === "joined" ? (
                        <span>* {msg.username} has joined</span>
                    ) : (
                        <>
                            <div className="message-content">
                                <strong>{msg.username}</strong>
                                <span>{msg.text}</span>
                            </div>
                            <span className="timestamp">{new Date(msg.timestamp).toLocaleTimeString()}</span>
                        </>
                    )}
                    {index === messages.length - 1 && <div ref={lastMessageRef}></div>}
                </div>
            ))}
        </div>
    );
}

function UsersBox({}) {
    const {users} = useRoom();
    return (
        <div className="users-box">
            {users.map((username, index) => (
                <div key={index}>{username}</div>
            ))}
        </div>
    );
}

function Chat() {
    const {user} = useUser();
    const {room, users, setUsers} = useRoom();
    const [messages, setMessages] = useState([]);
    //const [users, setUsers] = useState([]);
    const wsRef = useRef(null);
    const navigate = useNavigate();

    useEffect(() => {
        const connectionString = `ws://localhost:8000/chatroom/session/chat?room_id=${room.room_id}&user_id=${user.user_id}&username=${user.username}`
        const ws = new WebSocket(connectionString);
        wsRef.current = ws;

        ws.onopen = () => {
            console.log('WebSocket connection opened');
        };

        ws.onerror = (error) => {
            console.log('WebSocket Error:', error);
        }

        fetch(`http://localhost:8000/chatroom/session/messages/${room.room_id}`)
            .then(response => response.json())
            .then(data => {
                if (data) {
                    // Transforming the fetched data to match the format needed by ChatBox
                    const transformedMessages = data.map(item => ({
                        id: item.user_id + item.created_at, // Create a unique ID
                        username: item.username,
                        text: item.content,
                        type: "message",  // Assuming 'message' type is for user messages
                        timestamp: new Date(item.created_at).getTime()
                    }));

                    setMessages(transformedMessages);
                }
            });
    }, []);

    const exitChat = () => {
        const payload = {
            room_id: room.room_id,
            user_id: user.user_id,
            username: user.username
        };

        try {
            fetch('http://localhost:8000/chatroom/session/exit', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(payload)
            }).then(response => response)
                .then((response) => {
                    if (response.ok) {
                        navigate('/rooms');
                    } else {
                        console.error('Failed to exit the chatroom.');
                    }
                });
        } catch (error) {
            console.error('There was an error exiting the room:', error);
        }
    };

    useEffect(() => {
        wsRef.current.onmessage = (event) => {
            const messageData = JSON.parse(event.data);
            console.log((messageData))
            if (isChatMessage(messageData)) {
                handleChatMessage(messageData)
            } else if (isActionMessage(messageData)) {
                handleActionMessage(messageData)
            } else if (isStockMessage(messageData)) {
                console.log((messageData))
                handleStockMessage(messageData)
            }
        };

        const isChatMessage = (data) => {
            return data &&
                data.hasOwnProperty('user_id') &&
                data.hasOwnProperty('username') &&
                data.hasOwnProperty('content') &&
                data.hasOwnProperty('created_at')
        };

        const isActionMessage = (data) => {
            return data &&
                data.hasOwnProperty('type') &&
                data.hasOwnProperty('username') &&
                data.hasOwnProperty('user_id')
        };

        const isStockMessage = (data) => {
            return data &&
                data.hasOwnProperty('room_id') &&
                data.hasOwnProperty('bot_message') &&
                data.hasOwnProperty('created_at')
        };

        const handleChatMessage = (messageData) => {
            const formattedMessage = {
                id: messageData.id || new Date().getTime(),
                type: messageData.type || "message",
                username: messageData.username,
                text: messageData.content,
                timestamp: messageData.created_at
            };
            setMessages((prevMessages) => {
                const newMessages = [...prevMessages, formattedMessage];
                while (newMessages.length > 50) {
                    newMessages.shift(); // Remove the oldest message
                }
                return newMessages;
            });
        }

        const handleActionMessage = (messageData) => {
            console.log(messageData.type)
            if (messageData.type === 'join') {
                if (!users.includes(messageData.username)) {
                    setUsers(prevUsers => [...prevUsers, messageData.username]);
                }
            } else if (messageData.type === 'exit') {
                setUsers(prevUsers => prevUsers.filter(username => username !== messageData.username));
            }
        }

        const handleStockMessage = (messageData) => {
            const formattedMessage = {
                id: messageData.id || new Date().getTime(),
                type: "message",
                username: 'stock-bot',
                text: messageData.bot_message,
                timestamp: messageData.created_at
            };
            setMessages((prevMessages) => {
                const newMessages = [...prevMessages, formattedMessage];
                while (newMessages.length > 50) {
                    newMessages.shift(); // Remove the oldest message
                }
                return newMessages;
            });
        }
    }, [room.room_id, user, users]);

    useEffect(() => {
        const handleUnload = () => {
            if (wsRef.current) {
                wsRef.current.close();
            }
        };

        window.addEventListener('beforeunload', handleUnload);

        return () => {
            window.removeEventListener('beforeunload', handleUnload);
        };
    }, []);

    const handleSend = (message, username, user_id, created_at) => {
        if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
            wsRef.current.send(JSON.stringify({
                username: username,
                user_id: user_id,
                created_at: created_at,
                content: message
            }));
        }
    };


    return (
        <>
            <button
                onClick={exitChat}
                style={{
                    padding: '8px 12px',
                    fontSize: '16px',
                    zIndex: 1000,
                    position: 'fixed',
                    top: '10px',
                    right: '10px',
                    background: '#333',
                    color: '#fff',
                    border: 'none',
                    borderRadius: '5px'
                }}>
                Exit
            </button>
            <div className="chat-container">
                <div className="content">
                    <ChatBox messages={messages}/>
                    <UsersBox users={users}/>
                </div>
                <ChatInput onSend={handleSend}/>
            </div>
        </>
    );
}


export default Chat;
