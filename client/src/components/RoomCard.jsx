import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { authFetch } from "../util";

import JoinRoomDialog from "./JoinRoomDialog"

function RoomCard({ room }) {
    const navigate = useNavigate()
    const [showDialog, setShowDialog] = useState(false);
    const [password, setPassword] = useState("");

    async function handleJoinClick() {
        if (room.private) {
            setShowDialog(true);
        } else {
            await joinRoom();
        }
    }

    async function joinRoom() {
        const res = await authFetch("/api/join-room", {
            method: "POST",
            body: JSON.stringify({
                "room_id": room.id,
                "password": password,
            })
        });
        if (res.ok) {
            setShowDialog(false);
            setPassword("");
            navigate(`/room/${room.id}`, { replace: true })
        }
        if (res.status == 401) {
            alert("wrong room password")
        }
    }

    return (
        <>
            <li className="border p-4 rounded bg-white shadow flex justify-between items-center">
                <div>
                    <h2 className="text-lg font-semibold">{room.name}</h2>
                    <p className="text-sm text-gray-500">
                        {room.private ? "🔒 Private " : "🌐 Public"}
                    </p>
                </div>
                <button
                    className="bg-green-500 text-white px-3 py-1 rounded hover:bg-green-600"
                    onClick={ handleJoinClick }
                >
                    Join
                </button>
            </li>

            {showDialog && room.private && (
                <JoinRoomDialog
                    room={room}
                    onClose={() => setShowDialog(false)}
                    onConfirm={{joinRoom, setPassword}}
                    setPassword={setPassword}
                />
            )}
        </>
    );
}

export default RoomCard;
