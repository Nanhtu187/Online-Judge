import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import Navbar from '../components/Navbar';
import { useAuth } from '../AuthContext';
import { listSubmissions } from '../api/api';
import type { Submission } from '../types';
import { SubmissionStatus } from '../types';
import { Clock, Code2, User as UserIcon, Mail, Calendar } from 'lucide-react';

const getStatusColor = (status: SubmissionStatus) => {
  switch (status) {
    case SubmissionStatus.ACCEPTED: return 'text-green-500';
    case SubmissionStatus.COMPILE_ERROR:
    case SubmissionStatus.RUNTIME_ERROR:
    case SubmissionStatus.LIMIT_EXCEEDED:
    case SubmissionStatus.WRONG_ANSWER:
    case SubmissionStatus.FAILED:
      return 'text-red-500';
    case SubmissionStatus.PENDING:
    case SubmissionStatus.RUNNING:
      return 'text-blue-400';
    default: return 'text-gray-400';
  }
};

const PortfolioPage: React.FC = () => {
  const { user, loading: authLoading } = useAuth();
  const [submissions, setSubmissions] = useState<Submission[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (user) {
      listSubmissions({ user_id: user.id })
        .then((data) => {
          setSubmissions(data);
          setLoading(false);
        })
        .catch((err) => {
          console.error('Failed to fetch personal submissions:', err);
          setLoading(false);
        });
    }
  }, [user]);

  if (authLoading) {
    return (
      <div className="min-h-screen bg-[#1a1a1a] flex justify-center items-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[#1a1a1a] text-gray-200">
      <Navbar />

      <div className="max-w-6xl mx-auto py-10 px-6">
        {/* Profile Header */}
        <div className="bg-[#1e1e1e] rounded-xl border border-gray-800 p-8 mb-10 shadow-2xl flex flex-col md:flex-row items-center md:items-start space-y-6 md:space-y-0 md:space-x-8">
          <div className="w-24 h-24 bg-orange-600/10 rounded-full flex items-center justify-center border border-orange-600/20 shrink-0">
            <UserIcon className="text-orange-500" size={40} />
          </div>
          <div className="flex-1 text-center md:text-left">
            <h1 className="text-3xl font-bold text-white mb-2">{user?.name}</h1>
            <div className="flex flex-wrap justify-center md:justify-start gap-4 text-sm text-gray-400">
              <div className="flex items-center space-x-1.5">
                <Mail size={14} className="text-gray-500" />
                <span>{user?.email}</span>
              </div>
              <div className="flex items-center space-x-1.5">
                <Calendar size={14} className="text-gray-500" />
                <span>Joined May 2026</span>
              </div>
            </div>
            
            <div className="mt-6 flex space-x-8 border-t border-gray-800 pt-6">
              <div>
                <div className="text-2xl font-bold text-white">{submissions.length}</div>
                <div className="text-[10px] uppercase text-gray-500 font-bold tracking-wider">Submissions</div>
              </div>
              <div>
                <div className="text-2xl font-bold text-green-500">
                  {new Set(submissions.filter(s => s.status === SubmissionStatus.ACCEPTED).map(s => s.problem_id)).size}
                </div>
                <div className="text-[10px] uppercase text-gray-500 font-bold tracking-wider">Solved</div>
              </div>
            </div>
          </div>
        </div>

        {/* Submissions Section */}
        <h2 className="text-xl font-semibold mb-6 flex items-center">
          <Clock size={20} className="mr-2 text-orange-500" />
          Recent Submissions
        </h2>

        {loading ? (
          <div className="flex justify-center py-20">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500"></div>
          </div>
        ) : (
          <div className="bg-[#1e1e1e] rounded-lg border border-gray-800 overflow-hidden shadow-xl">
            <table className="w-full text-left border-collapse text-sm">
              <thead>
                <tr className="bg-[#2d2d2d] text-gray-400 font-medium">
                  <th className="py-4 px-6">Problem</th>
                  <th className="py-4 px-6">Status</th>
                  <th className="py-4 px-6">Language</th>
                  <th className="py-4 px-6">Date</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-800">
                {submissions.length === 0 ? (
                  <tr>
                    <td colSpan={4} className="py-10 text-center text-gray-500 italic">No submissions found.</td>
                  </tr>
                ) : (
                  submissions.map((s) => (
                    <tr key={s.id} className="hover:bg-[#252526] transition-colors">
                      <td className="py-4 px-6">
                        <Link 
                          to={`/problems/${s.problem_id}`} 
                          className="text-blue-400 hover:text-blue-300 font-medium transition-colors"
                        >
                          {s.problem_title || 'Unknown Problem'}
                        </Link>
                      </td>
                      <td className="py-4 px-6">
                        <span className={`font-bold ${getStatusColor(s.status)}`}>
                          {s.status.replace('SUBMISSION_STATUS_', '').replace('_', ' ')}
                        </span>
                      </td>
                      <td className="py-4 px-6">
                        <div className="flex items-center space-x-1.5">
                          <Code2 size={14} className="text-gray-500" />
                          <span className="uppercase text-gray-400 font-medium">{s.language}</span>
                        </div>
                      </td>
                      <td className="py-4 px-6 text-gray-500">
                        <div className="flex items-center space-x-1.5 text-xs">
                          <Clock size={12} />
                          <span>{s.created_at}</span>
                        </div>
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

export default PortfolioPage;
