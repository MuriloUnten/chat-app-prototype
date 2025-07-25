import { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, useNavigate, useParams } from "react-router-dom";
import { authFetch, fetchJSON, RequireAuth } from "../util.js"


function Room() {
    const { roomId } = useParams();
    const [message, setMessage] = useState("");

    const handleSendMessage = () => {
        // this is where you'll send messages later
        console.log("Send:", message);
        setMessage("");
    };

    return (
        <div className="max-w-2xl mx-auto mt-8 flex flex-col h-[90vh]">
            <h1 className="text-2xl font-bold mb-4">Room #{roomId}</h1>

            {/* Message display area */}
            <div className="flex-1 overflow-y-auto border rounded p-4 bg-white shadow mb-4">
                {/* Messages will go here */}
                <p className="text-gray-500 italic">No messages yet</p>
            </div>

            {/* Message input */}
            <div className="flex gap-2">
                <input
                    type="text"
                    placeholder="Type a message..."
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                    className="flex-1 border p-2 rounded"
                />
                <button
                    className="bg-blue-500 text-white px-4 py-2 rounded"
                    onClick={handleSendMessage}
                >
                    Send
                </button>
            </div>
        </div>
    );
}

export default Room;
