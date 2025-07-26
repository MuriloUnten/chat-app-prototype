import { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, useNavigate, useParams } from "react-router-dom";

import LoginForm from "./pages/LoginForm";
import RegisterForm from "./pages/RegisterForm";
import Dashboard from "./pages/Dashboard";
import CreateRoom from "./pages/CreateRoom";
import Room from "./pages/Room";
import AuthenticatedLayout from "./AuthenticadedLayout";

export default function App() {
    const token = localStorage.getItem("token")

    return (
        <>
            <Router>
                <Routes>
                    <Route path="/login" element={<LoginForm />} />
                    <Route path="/register" element={<RegisterForm />} />
                    <Route
                        path="/"
                        element={
                            <AuthenticatedLayout>
                                <Dashboard />
                            </AuthenticatedLayout>
                        }
                    />
                    <Route
                        path="/create-room"
                        element={
                            <AuthenticatedLayout>
                                <CreateRoom />
                            </AuthenticatedLayout>
                        }
                    />
                    <Route
                        path="/room/:roomId"
                        element={
                            <AuthenticatedLayout>
                                <Room />
                            </AuthenticatedLayout>
                        }
                    />
                </Routes>
            </Router>
        </>
    );
}
