import { Navigate } from "react-router-dom";
import Navbar from "./components/Navbar";

function AuthenticatedLayout({ children }) {
    const token = localStorage.getItem("token");

    if (!token) {
        return <Navigate to="/login" replace />;
    }

    return (
        <>
            <Navbar />
            <div className="p-4">{children}</div>
        </>
    );
}

export default AuthenticatedLayout;
