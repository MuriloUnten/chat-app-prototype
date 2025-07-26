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

            <div className="relative">
                <button
                    onClick={() => setIsMenuOpen(!isMenuOpen)}
                    className="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center"
                >
                    <span className="text-sm">👤</span>
                </button>

                {isMenuOpen && (
                    <div className="absolute right-0 mt-2 w-48 bg-white shadow-lg rounded p-2 z-10">
                        <button className="w-full text-left px-2 py-1 hover:bg-gray-100">
                            Profile
                        </button>
                        <button className="w-full text-left px-2 py-1 hover:bg-gray-100">
                            Settings
                        </button>
                        <button
                            className="w-full text-left px-2 py-1 hover:bg-gray-100"
                            onClick={ handleLogout }
                        >
                            Logout
                        </button>
                    </div>
                )}
            </div>
        </nav>
    );
}

export default Navbar;

