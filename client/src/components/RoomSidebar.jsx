import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { authFetch } from "../util";

function RoomSidebar() {
    const navigate = useNavigate();
    const [rooms, setRooms] = useState([]);
    const [collapsed, setCollapsed] = useState(false);

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
        <div
            className={`h-full p-4 border-r transition-all duration-300 ease-in-out
            ${collapsed ? "w-16" : "w-64"}`}
        >
            <div className="flex justify-between items-center mb-4">
                {!collapsed && (
                    <h2 className="text-md font-semibold">Your Rooms</h2>
                )}
                <button
                    onClick={() => setCollapsed(!collapsed)}
                    className="text-gray-600 text-sm focus:outline-none"
                    title={collapsed ? "Expand" : "Collapse"}
                >
                    {collapsed ? "»" : "«"}
                </button>
            </div>

            {!collapsed && (
                <>
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
                </>
            )}
        </div>
    );
}

export default RoomSidebar;
