import { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, useNavigate, useParams } from "react-router-dom";
import { fetchJSON } from "../util.js"

function LoginForm() {
    const navigate = useNavigate();
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    async function handleLogin() {
        let user = {
            name: username,
            password: password
        }
        console.log(JSON.stringify({ user }))

        const res = await fetchJSON("/api/login", {
            method: "POST",
            body: JSON.stringify({ user }),
        });

        if (res.ok) {
            const data = await res.json();
            console.log("jwt:", data.token)
            localStorage.setItem("token", data.token)
            navigate("/", { replace: true })
        } else {
            alert("Login failed");
        }
    }

    return (
        <div className="max-w-sm mx-auto mt-10 space-y-4">
            <h1 className="text-2xl font-bold">Login</h1>
            <input
                className="border w-full p-2 rounded"
                placeholder="Username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
            />
            <input
                className="border w-full p-2 rounded"
                type="password"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
            />
            <button className="bg-blue-500 text-white px-4 py-2 rounded" onClick={handleLogin}>
                Login
            </button>
            <p className="text-sm">
                Don't have an account?{" "}
                <button className="text-blue-500 underline" onClick={ () => navigate("/register", { replace: true }) }>
                    Register
                </button>
            </p>
        </div>
    );
}

export default LoginForm;
