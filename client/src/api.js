import { authFetch, fetchJSON } from "./util";

async function getUserById(userId) {
    const res = await fetchJSON(`/api/user/${userId}`);
    if (res.ok) {
        return res.json();
    }
}

async function getRooms() {
    const res = await fetchJSON("/api/room");
    if (res.ok) {
        return res.json();
    }
}

async function getJoinedRooms() {
    const res = await authFetch("/api/user/rooms");
    if (res.ok) {
        return res.json();
    }
}

async function getRoomMembers(roomId) {
    const res = await authFetch(`/api/room/${roomId}/users`);
    if (res.ok) {
        return res.json();
    }
}

export const api = {
    getUserById,
    getRooms,
    getJoinedRooms,
    getRoomMembers,
}
