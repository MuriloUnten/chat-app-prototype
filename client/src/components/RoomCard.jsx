import { useState } from "react";

import JoinRoomDialog from "./JoinRoomDialog"

function RoomCard({ room }) {
    const [showDialog, setShowDialog] = useState(false);
    const [password, setPassword] = useState("");

    const handleJoinClick = () => setShowDialog(true);

    const handleConfirmJoin = () => {
        console.log("Joining room", room.id, "with password", password);
        setShowDialog(false);
        setPassword("");
    };

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
                    onConfirm={handleConfirmJoin}
                    setPassword={setPassword} // pass setter for password input
                />
            )}
        </>
    );
}

export default RoomCard;
