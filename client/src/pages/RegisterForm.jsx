import { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, useNavigate, useParams } from "react-router-dom";
import { fetchJSON } from "../util.js"

function RegisterForm() {
    const navigate = useNavigate();
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    async function handleRegister() {
        let user = {
            name: username,
            password: password
        }

        const res = await fetchJSON("/api/user", {
            method: "POST",
            body: JSON.stringify({ user }),
        });

        if (res.ok) {
            const data = await res.json()
            localStorage.setItem("token", data.token)
        } else {
            alert("Registration failed");
        }
    }

    return (
        <div className="max-w-sm mx-auto mt-10 space-y-4">
            <h1 className="text-2xl font-bold">Register</h1>
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
            <button className="bg-green-500 text-white px-4 py-2 rounded" onClick={handleRegister}>
                Register
            </button>
            <p className="text-sm">
                Already have an account?{" "}
                <button className="text-blue-500 underline" onClick={ () => navigate("/login", { replace: true }) }>
                    Login
                </button>
            </p>
        </div>
    );
}

export default RegisterForm;
