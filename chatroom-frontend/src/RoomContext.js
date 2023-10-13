import React, {createContext, useContext, useState} from 'react';

const RoomContext = createContext();

export const RoomProvider = ({children}) => {
    const [room, setRoom] = useState(null);
    const [users, setUsers] = useState([]);

    return (
        <RoomContext.Provider value={{room, setRoom, users, setUsers}}>
            {children}
        </RoomContext.Provider>
    );
};

export const useRoom = () => {
    return useContext(RoomContext);
};