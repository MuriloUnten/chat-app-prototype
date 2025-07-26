function JoinRoomDialog({ room, onClose, onConfirm }) {
    return (
        // backdrop
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            {/* dialog box */}
            <div className="bg-white rounded p-6 max-w-md w-full shadow-lg">
                <h2 className="text-xl font-bold mb-4">Join Room: {room.name}</h2>
                <p className="mb-4">
                    This room is private. Please enter the password to join.
                </p>

                {room.private && (
                    <input
                        type="password"
                        placeholder="Password"
                        className="border w-full p-2 rounded mb-4"
                        // You can manage this input's state via props or inside this component
                        // For simplicity, handle password state and pass it up on confirm
                        onChange={(e) => onConfirm.setPassword(e.target.value)}
                    />
                )}

                <div className="flex justify-end gap-4">
                    <button
                        className="px-4 py-2 rounded border"
                        onClick={onClose}
                    >
                        Cancel
                    </button>
                    <button
                        className="bg-blue-600 text-white px-4 py-2 rounded"
                        onClick={() => onConfirm.joinRoom()}
                    >
                        Join
                    </button>
                </div>
            </div>
        </div>
    );
}

export default JoinRoomDialog;
