import { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, useNavigate, useParams } from "react-router-dom";
import { authFetch, fetchJSON, RequireAuth } from "../util.js"

function CreateRoom() {
    const navigate = useNavigate();
    const [roomName, setRoomName] = useState("");
    const [password, setPassword] = useState("");
    const [isPrivate, setPrivate] = useState(false);

    async function handleCreateRoom() {
        const room = {
            name: roomName,
            password: password,
            "private": isPrivate
        }

        const res = await authFetch("/api/room", {
            method: "POST",
            body: JSON.stringify({ room }),
        });

        if (res.ok) {
            const data = await res.json()
            const joinResponse = await authFetch("/api/join-room", {
                method: "POST",
                body: JSON.stringify({
                    "room_id": data.room.id,
                    "password": password,
                })
            });
            if (joinResponse.ok) {
                console.log("room joined. Must navigate to it")
                navigate(`/room/${data.room.id}`, { replace: true })
            }
            else {
                alert("failed to join room that was just created")
            }

        } else {
            alert("error creating room");
        }
    }

    return (
        <div className="max-w-sm mx-auto mt-10 space-y-4">
            <h1 className="text-2xl font-bold">Create Room</h1>
            <input
                className="border w-full p-2 rounded"
                placeholder="Room Name"
                value={roomName}
                onChange={(e) => setRoomName(e.target.value)}
            />
            <input
                className="border w-full p-2 rounded"
                type="password"
                placeholder="Password"
                value={password}
                disabled={ !isPrivate }
                onChange={(e) => setPassword(e.target.value)}
            />
            <div className="flex items-center gap-4">
                <button
                    onClick={() => setPrivate(!isPrivate)}
                    className={`relative w-12 h-6 rounded-full transition-colors ${
                        isPrivate ? "bg-green-500" : "bg-gray-300"
                    }`}
                >
                    <span
                        className={`absolute top-0.5 left-0.5 w-5 h-5 rounded-full bg-white transition-transform ${
                            isPrivate ? "translate-x-6" : "translate-x-0"
                        }`}
                    />
                </button>
                <span className="text-sm font-medium">Private</span>
            </div>
            <button className="bg-green-500 text-white px-4 py-2 rounded" onClick={handleCreateRoom}>
                Create Room
            </button>
        </div>
    );
}

export default CreateRoom;
