import { Navigate } from "react-router-dom";
import Navbar from "./components/Navbar";
import RoomSidebar from "./components/RoomSidebar";

function AuthenticatedLayout({ children }) {
    const token = localStorage.getItem("token");

    if (!token) {
        return <Navigate to="/login" replace />;
    }

    return (
        <>
            <Navbar />
            <div className="flex h-screen">
                <RoomSidebar />
                <main className="flex-1 p-4 flex justify-center items-start overflow-auto">
                    <div className="w-full max-w-xl">{children}</div>
                </main>
            </div>
        </>
    );
}

export default AuthenticatedLayout;
