import React, {useState} from 'react';
import {useUser} from "./UserContext";
import {useNavigate} from "react-router-dom";

function Home({onLogin}) {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [message, setMessage] = useState('');
    const [loggedIn, setLoggedIn] = useState(false);
    const {setUser} = useUser();
    const navigate = useNavigate();

    const handleLogin = async () => {
        if (username === '' || password === '') {
            setMessage('Empty username or password.');
        } else {
            try {
                const response = await fetch('http://localhost:8000/chatroom/user/login', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({username, password})
                });
                const res = await response.json()
                if (response.status !== 200) {
                    setMessage('Error: ' + res.message)
                } else {
                    setUser({username: username, user_id: res.user_id})
                    setLoggedIn(true)
                    navigate('/rooms');
                }
            } catch (error) {
                setMessage('Error logging in.');
            }
        }
    };

    const handleSignUp = async () => {
        if (username === '' || password === '') {
            setMessage('Empty username or password.');
        } else {
            try {
                const response = await fetch('http://localhost:8000/chatroom/user', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({username, password})
                });

                if (response.status !== 201) {
                    const res = await response.json()
                    setMessage('Error: ' + res.message)
                } else {
                    setMessage('Account created successfully. Please log in');
                }
            } catch (error) {
                setMessage('Error creating account.');
            }
        }
    };

    return (
        <div>
            {message && <div>{message}</div>}

            <input
                type="text"
                value={username}
                onChange={e => setUsername(e.target.value)}
                placeholder="Username"
            />
            <input
                type="password"
                value={password}
                onChange={e => setPassword(e.target.value)}
                placeholder="Password"
            />
            <button className="login-button" onClick={handleLogin}>Login</button>
            <button className="signup-button" onClick={handleSignUp}>Sign Up</button>
        </div>
    );

}

export default Home;