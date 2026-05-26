import React from 'react';
import { Link } from 'react-router-dom';
import { ShieldAlert } from 'lucide-react';
import Navbar from '../components/Navbar';

const ForbiddenPage: React.FC = () => {
  return (
    <div className="min-h-screen bg-[#1a1a1a] text-gray-200 flex flex-col">
      <Navbar />
      <div className="flex-1 flex flex-col items-center justify-center p-6">
        <div className="w-20 h-20 bg-red-600/10 rounded-full flex items-center justify-center mb-6 border border-red-600/20">
          <ShieldAlert className="text-red-500" size={40} />
        </div>
        <h1 className="text-4xl font-bold text-white mb-2">403 - Forbidden</h1>
        <p className="text-gray-400 text-lg mb-8 text-center max-w-md">
          You don't have permission to access this page. Please contact an administrator if you believe this is an error.
        </p>
        <Link
          to="/"
          className="px-6 py-2.5 bg-[#3c3c3c] hover:bg-[#4c4c4c] text-white font-medium rounded-lg transition-colors"
        >
          Back to Problems
        </Link>
      </div>
    </div>
  );
};

export default ForbiddenPage;
