import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";

function Navbar() {
    const navigate = useNavigate();
    const [isMenuOpen, setIsMenuOpen] = useState(false);

    function handleLogout() {
        setIsMenuOpen(false);
        console.log("logging out")

        localStorage.removeItem("token");
        navigate("/login", { replace: true });
    }

    return (
        <nav className="bg-white shadow px-6 py-4 flex justify-between items-center">
            <div className="text-lg font-bold">
                <Link to="/dashboard">Dashboard</Link>
            </div>

            <div className="relative w-20 h-8 rounded-full bg-gray-200 flex items-center justify-center">
                <button
                    onClick={ handleLogout}
                    className=""
                >
                    Logout
                </button>
            </div>
        </nav>
    );
}

export default Navbar;

