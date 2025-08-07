import { useEffect } from "react";
import useChatStore from "./store";

const MsgType = Object.freeze({
    Chat: "chat",
    Event: "event",
    Error: "error",
});

const EventType = Object.freeze({
    UserJoined:  "user_joined",
	UserLeft:    "user_left",
	RoomCreated: "room_created",
	RoomDeleted: "room_deleted",
});

export function useWebSocketConnection() {
    useEffect(() => {
        const ws = WebSocket("/api/ws");

        ws.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            const { type, data } = msg;

            handleMsg(type, data);
        };

        
        return () => {
            ws.close();
        };
    }, []);
}

function handleMsg(type, data) {
    try {
        switch (type) {
            case MsgType.Chat:
                handleChatMsg(data);
            case MsgType.Event:
                handleEventMsg(data);
            case MsgType.Error:
                handleErrorMsg(data);

            default:
                console.log("received fucked message");
                break;
        }
    }
    catch (e) {
        console.log("received malformed message");
        console.log(e);
    }
}

function handleChatMsg(data) {
    validateChatMsg(data);

    const { addMessage } = useChatStore();

    const msg = {
        senderId: data.sender.id,
        senderName: data.sender.name,
        content: data.content,
    };

    addMessage(data.room_id, msg);
}

function handleEventMsg(data) {
    validateEventMsg(data);

    const { addRoom, deleteRoom } = useChatStore();

    switch (data.event) {
        case EventType.UserJoined:

        case EventType.UserLeft:

        case EventType.RoomCreated:
            addRoom(data.room_id);
        case EventType.RoomDeleted:

        default:
            throw new Error(`Invalid event type ${data.event}`);
    }
}

function handleErrorMsg(data) {
    validateErrorMsg(data);
}

function validateChatMsg(data) {
    if (!data.room_id) {
        throw new Error("Missing room id");
    }
    if (!data.sender) {
        throw new Error("Missing sender data");
    }
    if (!data.sender.id) {
        throw new Error("Missing user id");
    }
    if (!data.sender.name) {
        throw new Error("Missing user name");
    }
    if (!data.content) {
        throw new Error("Missing content");
    }
}

function validateEventMsg(data) {
    if (!data.event) {
        throw new Error("Missing event type");
    }
    if (!data.user_id) {
        throw new Error("Missing user id");
    }
    if (!data.room_id) {
        throw new Error("Missing room id");
    }
}

function validateErrorMsg(data) {
    if (!data.error) {
        throw new Error("Missing error");
    }
}
