import { useEffect } from "react";
import useChatStore from "./store";
import useWebSocket, { ReadyState } from "react-use-websocket"

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
    const token = localStorage.getItem("token");
    const WS_URL = "ws://127.0.0.1:8080/api/ws";
    const { sendJsonMessage, lastJsonMessage, readyState } = useWebSocket(
        WS_URL,
        {
            share: false,
            shouldReconnect: (e) => {
                console.log(e);
                return true
            },
            protocols: [token]
        },
    );

    const connectionStatus = {
        [ReadyState.CONNECTING]: 'Connecting',
        [ReadyState.OPEN]: 'Open',
        [ReadyState.CLOSING]: 'Closing',
        [ReadyState.CLOSED]: 'Closed',
        [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
    }[readyState];

    useEffect(() => {
        console.log("Connection state changed to", connectionStatus);
    }, [readyState]);

    useEffect(() => {
        console.log("Got a new message:");
        console.log(lastJsonMessage);
        if (lastJsonMessage != null) {
            handleMsg(lastJsonMessage.type, lastJsonMessage.data);
        }
    }, [lastJsonMessage]);
}

export function handleMsg(type, data) {
    try {
        switch (type) {
            case MsgType.Chat:
                handleChatMsg(data);
                break;
            case MsgType.Event:
                handleEventMsg(data);
                break;
            case MsgType.Error:
                handleErrorMsg(data);
                break;

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

    const { addMessageToRoom } = useChatStore.getState();

    const msg = {
        senderId: data.sender.id,
        senderName: data.sender.name,
        content: data.content,
    };

    addMessageToRoom(data.room_id, msg);
}

function handleEventMsg(data) {
    validateEventMsg(data);

    const { addMember, removeMember, addAvailableRoom, deleteRoom } = useChatStore.getState();

    switch (data.event) {
        case EventType.UserJoined:
            addMember(data.room_id, data.member)
            break;

        case EventType.UserLeft:
            removeMember(data.room_id, data.user_id);
            break;

        case EventType.RoomCreated:
            addAvailableRoom(data.room_id);
            break;

        case EventType.RoomDeleted:
            deleteRoom(data.room_id);
            break;
        default:
            throw new Error(`Invalid event type ${data.event}`);
    }
}

function handleErrorMsg(data) {
    validateErrorMsg(data);
    console.log("error message:", data.error);
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
    if (!data.text) {
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
