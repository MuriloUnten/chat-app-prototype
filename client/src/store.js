import { create } from "zustand";

const useChatStore = create((set, get) => ({
    availableRooms: [],
    joinedRooms: {},

    setAvailableRooms: (rooms) => set({ availableRooms: rooms }),

    addAvailableRoom: (roomId) => {
        const current = get().availableRooms;
        set({ availableRooms: [...current, roomId]});
    },

    joinRoom: (roomId) => {
        const current = get().joinedRooms;
        if (current[roomId]) {
            return;
        }

        set({
            joinedRooms: {
                ...current,
                [roomId]: { messages: [], members: [] },
            },
        });
    },

    leaveRoom: (roomId) => {
        const current = get().joinedRooms;
        const { [roomId]: _, ...rest } = current;
        set({ joinedRooms: rest });
    },

    deleteRoom: (roomId) => {
        const joined = get().joinedRooms;
        const { [roomId]: _, ...restJoined } = joined;
        const available = get().availableRooms;
        const { [roomId]: __, ...restAvailable } = available;

        set({
            joinedRooms: restJoined,
            availableRooms: restAvailable,
        });
    },

    setRoomMembers: (roomId, members) => {
        const current = get().joinedRooms;
        const room = current[roomId];
        if (!room) {
            return;
        }

        set({
            joinedRooms: {
                ...current,
                [roomId]: {
                    ...room,
                    messages: room.messages,
                    members,
                },
            },
        });
    },

    addMember: (roomId, memberId) => {
        const current = get().joinedRooms;
        const room = current[roomId];
        if (!room) {
            return;
        }

        const alreadyIn = room.members.find((m) => m === memberId);
        if (alreadyIn) {
            return;
        }

        set({
            joinedRooms: {
                ...current,
                [roomId]: {
                    ...room,
                    messages: room.messages,
                    members: [...room.members, memberId],
                }
            }
        });
    },

    removeMember: (roomId, memberId) => {
        const current = get().joinedRooms;
        const room = current[roomId];
        if (!room) {
            return;
        }

        const alreadyIn = room.members.find(m => m === memberId);
        if (!alreadyIn) {
            return;
        }

        set({
            joinedRooms: {
                ...current,
                [roomId]: {
                    messages: room.messages,
                    members: room.members.filter((m) => m !== memberId),
                },
            },
        });
    },

    addMessageToRoom: (roomId, message) => {
        const current = get().joinedRooms;
        const room = current[roomId];
        if (!room) {
            return;
        }

        set({
            joinedRooms: {
                ...current,
                [roomId]: {
                    ...room,
                    messages: [...room.messages, message],
                    members: room.members,
                },
            },
        });
    }
}));

export default useChatStore;
