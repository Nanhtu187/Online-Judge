import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../AuthContext';
import { LogOut, User as UserIcon, LogIn } from 'lucide-react';

const Navbar: React.FC = () => {
  const { user, isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="flex items-center justify-between px-8 h-14 bg-[#252526] border-b border-gray-700 shrink-0">
      <div className="flex items-center">
        <Link to="/" className="text-white font-bold text-lg mr-10">Online Judge</Link>
        <nav className="space-x-6 h-full flex items-center">
          <Link to="/" className="text-sm font-medium text-gray-300 hover:text-white transition-colors">Problems</Link>
          <Link to="/submissions" className="text-sm font-medium text-gray-300 hover:text-white transition-colors">All Submissions</Link>
        </nav>
      </div>

      <div className="flex items-center space-x-4">
        {isAuthenticated ? (
          <div className="flex items-center space-x-4">
            <Link 
              to="/portfolio" 
              className="flex items-center space-x-2 text-sm text-gray-300 hover:text-white transition-colors"
            >
              <UserIcon size={16} className="text-orange-500" />
              <span>{user?.name}</span>
            </Link>
            <button
              onClick={handleLogout}
              className="flex items-center space-x-1 px-3 py-1.5 bg-[#3c3c3c] hover:bg-[#4c4c4c] text-gray-200 text-xs rounded transition-colors"
            >
              <LogOut size={14} />
              <span>Logout</span>
            </button>
          </div>
        ) : (
          <div className="flex items-center space-x-3">
            <Link
              to="/login"
              className="flex items-center space-x-1 px-3 py-1.5 text-gray-300 hover:text-white text-sm transition-colors"
            >
              <LogIn size={16} />
              <span>Login</span>
            </Link>
            <Link
              to="/register"
              className="px-4 py-1.5 bg-orange-600 hover:bg-orange-500 text-white text-sm font-medium rounded transition-colors"
            >
              Register
            </Link>
          </div>
        )}
      </div>
    </div>
  );
};

export default Navbar;
