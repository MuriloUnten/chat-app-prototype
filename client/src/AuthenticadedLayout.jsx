import { Navigate } from "react-router-dom";
import Navbar from "./components/Navbar";
import RoomSidebar from "./components/RoomSidebar";

function AuthenticatedLayout({ children }) {
    const token = localStorage.getItem("token");

    if (!token) {
        return <Navigate to="/login" replace />;
    }

    return (
        <div className="flex flex-col h-screen">
            <Navbar />
            <div className="flex flex-1 overflow-hidden">
                <RoomSidebar />
                <main className="flex-1 p-4 flex justify-center items-start overflow-auto">
                    <div className="h-full w-full max-w-xl">{children}</div>
                </main>
            </div>
        </div>
    );
}

export default AuthenticatedLayout;
