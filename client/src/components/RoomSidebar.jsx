import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { authFetch } from "../util";

function RoomSidebar() {
    const [rooms, setRooms] = useState([]);
    const navigate = useNavigate();

    useEffect(() => {
        async function fetchRooms() {
            try {
                const res = await authFetch("/api/user/rooms");
                if (!res.ok) {
                    throw new Error("Failed to fetch rooms");
                }

                const data = await res.json();
                setRooms(data.rooms);
            } catch (err) {
                console.error(err);
            }
        }

        fetchRooms();
    }, []);

    return (
        <div className="w-64 h-full p-4 space-y-2 border-r">
            <h2 className="text-lg font-semibold mb-2">Your Rooms</h2>
            {rooms.length === 0 ? (
                <p className="text-sm text-gray-500">You are not in any rooms.</p>
            ) : (
                <ul className="space-y-1">
                    {rooms.map((room) => (
                        <li
                            key={room.id}
                            className="p-2 rounded cursor-pointer hover:bg-gray-200"
                            onClick={() => navigate(`/room/${room.id}`, { replace: true })}
                        >
                            {room.name}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}

export default RoomSidebar;
