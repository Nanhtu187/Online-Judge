import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Plus, Edit2 } from 'lucide-react';
import Navbar from '../components/Navbar';
import { useAuth } from '../AuthContext';
import { listProblems } from '../api/api';
import type { ProblemSummary } from '../types';

const getDifficultyColor = (diff: string) => {
  switch (diff) {
    case 'DIFFICULTY_EASY': return 'text-green-500 bg-green-500/10 border-green-500/20';
    case 'DIFFICULTY_MEDIUM': return 'text-orange-500 bg-orange-500/10 border-orange-500/20';
    case 'DIFFICULTY_HARD': return 'text-red-500 bg-red-500/10 border-red-500/20';
    default: return 'text-gray-500 bg-gray-500/10 border-gray-500/20';
  }
};

const ProblemListPage: React.FC = () => {
  const { hasPermission } = useAuth();
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
      <Navbar />

      <div className="max-w-6xl mx-auto py-10 px-6">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-2xl font-semibold">Problem List</h1>
          {hasPermission('PROBLEM_CREATE') && (
            <Link
              to="/problems/new"
              className="flex items-center bg-orange-600 hover:bg-orange-500 text-white px-4 py-2 rounded transition-colors text-sm font-medium"
            >
              <Plus size={16} className="mr-2" /> Create Problem
            </Link>
          )}
        </div>

        {loading ? (
          <div className="flex justify-center py-10">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500"></div>
          </div>
        ) : (
          <div className="bg-[#1e1e1e] rounded-lg border border-gray-800 overflow-hidden shadow-xl">
            <table className="w-full text-left border-collapse text-sm">
              <thead>
                <tr className="bg-[#2d2d2d] text-gray-400 font-medium">
                  <th className="py-4 px-6 w-20">ID</th>
                  <th className="py-4 px-6">Title</th>
                  <th className="py-4 px-6 w-32">Difficulty</th>
                  <th className="py-4 px-6">Tags</th>
                  <th className="py-4 px-6 w-32 text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-800">
                {problems.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="py-10 text-center text-gray-500">No problems found.</td>
                  </tr>
                ) : (
                  problems.map((p, idx) => (
                    <tr key={p.id} className="hover:bg-[#252526] transition-colors group">
                      <td className="py-4 px-6 text-gray-500 font-mono text-xs">{idx + 1}</td>
                      <td className="py-4 px-6">
                        <Link 
                          to={`/problems/${p.id}`} 
                          className="text-blue-400 hover:text-blue-300 font-medium transition-colors"
                        >
                          {p.title}
                        </Link>
                      </td>
                      <td className="py-4 px-6">
                        <span className={`text-[10px] px-2 py-0.5 rounded border ${getDifficultyColor(p.difficulty)} font-bold uppercase tracking-wider`}>
                          {p.difficulty.replace('DIFFICULTY_', '')}
                        </span>
                      </td>
                      <td className="py-4 px-6">
                        <div className="flex flex-wrap gap-1.5">
                          {p.tags && p.tags.map(tag => (
                            <span key={tag} className="text-[9px] px-1.5 py-0.5 rounded bg-gray-800 text-gray-400 border border-gray-700">
                              {tag}
                            </span>
                          ))}
                        </div>
                      </td>
                      <td className="py-4 px-6 text-right">
                        {hasPermission('PROBLEM_UPDATE') && (
                          <Link
                            to={`/problems/${p.id}/edit`}
                            className="text-gray-400 hover:text-orange-400 transition-colors inline-flex items-center p-1.5 hover:bg-[#3d3d3d] rounded-md"
                            title="Edit Problem"
                          >
                            <Edit2 size={16} />
                          </Link>
                        )}
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
