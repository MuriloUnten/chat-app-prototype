import { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, useNavigate, useParams } from "react-router-dom";
import { authFetch, fetchJSON, RequireAuth } from "./util.js"

import LoginForm from "./pages/LoginForm";
import RegisterForm from "./pages/RegisterForm";
import Dashboard from "./pages/Dashboard";
import CreateRoom from "./pages/CreateRoom";
import Room from "./pages/Room";

export default function App() {
    return (
        <Router>
            <Routes>
                <Route path="/login" element={<LoginForm />} />
                <Route path="/register" element={<RegisterForm />} />
                <Route
                    path="/"
                    element={
                        <RequireAuth>
                            <Dashboard />
                        </RequireAuth>
                    }
                />
                <Route
                    path="/create-room"
                    element={
                        <RequireAuth>
                            <CreateRoom />
                        </RequireAuth>
                    }
                />
                <Route
                    path="/room/:roomId"
                    element={
                        <RequireAuth>
                            <Room />
                        </RequireAuth>
                    }
                />
            </Routes>
        </Router>
    );
}

