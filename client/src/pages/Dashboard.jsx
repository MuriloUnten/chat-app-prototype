import { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, useNavigate, useParams } from "react-router-dom";

import RoomCard  from "../components/RoomCard"

function Dashboard() {
    const navigate = useNavigate();
    const [rooms, setRooms] = useState([]);

    useEffect(() => {
        fetchRooms();
    }, []);

    async function fetchRooms() {
        const res = await fetch("/api/room");
        if (res.ok) {
            const data = await res.json();
            setRooms(data);
        }
    }

    return (
        <div className="max-w-xl mx-auto mt-10 space-y-4">
            <h1 className="text-2xl font-bold">Dashboard</h1>
            <div className="flex gap-2">
                <button
                    className="bg-blue-500 text-white px-4 py-2 rounded"
                    onClick={() => navigate("/create-room")}
                >
                    Create Room
                </button>
                <button className="bg-gray-300 px-4 py-2 rounded" onClick={fetchRooms}>
                    Refresh
                </button>
            </div>
            <ul className="space-y-2">
                {rooms.map((room) => (
                    <RoomCard key={room.id} room={room} />
                ))}
            </ul>
        </div>
    );
}

export default Dashboard;
