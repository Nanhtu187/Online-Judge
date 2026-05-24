import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Plus, Edit2 } from 'lucide-react';
import { listProblems } from '../api/api';
import type { ProblemSummary } from '../types';

const ProblemListPage: React.FC = () => {
  const [problems, setProblems] = useState<ProblemSummary[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    listProblems()
      .then((data) => {
        setProblems(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error('Failed to fetch problems:', err);
        setLoading(false);
      });
  }, []);

  return (
    <div className="min-h-screen bg-[#1a1a1a] text-gray-200">
      {/* Header */}
      <div className="flex items-center px-8 h-14 bg-[#252526] border-b border-gray-700">
        <div className="text-white font-bold text-lg">Online Judge</div>
        <nav className="ml-10 space-x-6">
          <Link to="/" className="text-sm font-medium text-white border-b-2 border-orange-500 pb-4 mt-4 inline-block">Problems</Link>
        </nav>
      </div>

      <div className="max-w-5xl mx-auto py-10 px-6">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-2xl font-semibold">Problem List</h1>
          <Link
            to="/problems/new"
            className="flex items-center bg-orange-600 hover:bg-orange-500 text-white px-4 py-2 rounded transition-colors text-sm font-medium"
          >
            <Plus size={16} className="mr-2" /> Create Problem
          </Link>
        </div>

        {loading ? (
          <div className="flex justify-center py-10">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500"></div>
          </div>
        ) : (
          <div className="bg-[#1e1e1e] rounded-lg border border-gray-800 overflow-hidden">
            <table className="w-full text-left border-collapse text-sm">
              <thead>
                <tr className="bg-[#2d2d2d] text-gray-400 font-medium">
                  <th className="py-3 px-6 w-20">ID</th>
                  <th className="py-3 px-6">Title</th>
                  <th className="py-3 px-6 w-32">Status</th>
                  <th className="py-3 px-6 w-24 text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-800">
                {problems.length === 0 ? (
                  <tr>
                    <td colSpan={4} className="py-10 text-center text-gray-500">No problems found.</td>
                  </tr>
                ) : (
                  problems.map((p, idx) => (
                    <tr key={p.id} className="hover:bg-[#252526] transition-colors">
                      <td className="py-4 px-6 text-gray-500">{idx + 1}</td>
                      <td className="py-4 px-6">
                        <Link 
                          to={`/problems/${p.id}`} 
                          className="text-blue-400 hover:text-blue-300 font-medium"
                        >
                          {p.title}
                        </Link>
                      </td>
                      <td className="py-4 px-6">
                        <span className="text-xs px-2 py-1 rounded bg-gray-800 text-gray-400 uppercase">Todo</span>
                      </td>
                      <td className="py-4 px-6 text-right">
                        <Link
                          to={`/problems/${p.id}/edit`}
                          className="text-gray-400 hover:text-orange-400 transition-colors inline-flex items-center"
                          title="Edit Problem"
                        >
                          <Edit2 size={16} />
                        </Link>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
};

export default ProblemListPage;
