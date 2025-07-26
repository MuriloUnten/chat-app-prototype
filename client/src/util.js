import { Navigate } from "react-router-dom";

export function authFetch(url, options = {}) {
    const token = localStorage.getItem("token");
    return fetchJSON(url, {
        ...options,
        headers: {
            ...options.headers,
            "Authorization": `Bearer ${token}`,
        }
    });
}

export function fetchJSON(url, options = {}) {
    return fetch(url, {
        ...options,
        headers: {
            ...options.headers,
            "Content-Type": "application/json"
        }
    });
}
