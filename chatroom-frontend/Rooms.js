import React, {useState, useEffect, memo} from 'react';
import {useNavigate} from 'react-router-dom';
import {useUser} from "./UserContext";
import {useRoom} from "./RoomContext";

const ChatUsersContext = React.createContext();

function Rooms() {
    const [rooms, setRooms] = useState([]);
    const navigate = useNavigate();
    const {user} = useUser();
    const {setRoom, setUsers} = useRoom();
    const [isLoading, setIsLoading] = useState(true);

    const fetchRooms = async () => {
        try {
            const response = await fetch('http://localhost:8000/chatroom/room');
            const data = await response.json();
            const activeRooms = data.rooms.filter(room => room.is_active);
            setRooms(activeRooms);
        } catch (error) {
            console.error('Error fetching rooms:', error);
        } finally {
            setIsLoading(false);
        }
    };


    useEffect(() => {
        fetchRooms()
    }, []);

    const joinRoom = async (room_id) => {
        try {

            const response = await fetch('http://localhost:8000/chatroom/session/join', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({
                    room_id: room_id,
                    user_id: user.user_id,
                    username: user.username
                })
            });

            if (response.status === 200) {
                const jsonRes = await response.json()
                setRoom({room_id: room_id})
                setUsers(jsonRes.users)
                navigate('/chat');
            } else {
                const res = await response.json();
                console.error('Failed to join room:', res.message);
            }
        } catch (error) {
            console.error('Error joining room:', error);
        }
    };

    if (isLoading) {
        return <div>Loading...</div>;  // Display a loading message or spinner
    } else if (user) {
        return (
            <div>
                <h2>Available Rooms</h2>
                <ul className="rooms-list">
                    {rooms.map(room => (
                        <li key={room.id} className="room-item">
                            <button className="room-button" onClick={() => joinRoom(room.id)}>
                                {room.name}
                            </button>
                        </li>
                    ))}
                </ul>
            </div>
        );
    } else {
        return (
            <div>
                <h2>Not logged in</h2>
            </div>
        );
    }
}

export default Rooms;