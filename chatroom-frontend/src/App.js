import React, {useState} from 'react';
import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';  // Notice the change from Switch to Routes
import './App.css';
import Home from './Home';
import Rooms from './Rooms';
import Chat from './Chat';
import {UserProvider} from "./UserContext";
import {RoomProvider} from "./RoomContext";

function App() {
    const [showRooms, setShowRooms] = useState(false);

    return (
        <UserProvider>
            <RoomProvider>
                <Router>
                    <div className="App">
                        <header className="App-header">
                            <Routes>
                                <Route path="/chat" element={<Chat/>}/>
                                <Route path="/rooms" element={<Rooms/>}/>
                                <Route path="/" element=<Home/>/>
                            </Routes>
                        </header>
                    </div>
                </Router>
            </RoomProvider>
        </UserProvider>
    );
}

export default App;