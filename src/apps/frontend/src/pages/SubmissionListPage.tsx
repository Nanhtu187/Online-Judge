import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import Navbar from '../components/Navbar';
import { listSubmissions } from '../api/api';
import type { Submission } from '../types';
import { SubmissionStatus } from '../types';
import { Clock, Code2 } from 'lucide-react';

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

const SubmissionListPage: React.FC = () => {
  const [submissions, setSubmissions] = useState<Submission[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    listSubmissions({})
      .then((data) => {
        setSubmissions(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error('Failed to fetch submissions:', err);
        setLoading(false);
      });
  }, []);

  return (
    <div className="min-h-screen bg-[#1a1a1a] text-gray-200">
      <Navbar />

      <div className="max-w-6xl mx-auto py-10 px-6">
        <h1 className="text-2xl font-semibold mb-8">Global Submission History</h1>

        {loading ? (
          <div className="flex justify-center py-20">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500"></div>
          </div>
        ) : (
          <div className="bg-[#1e1e1e] rounded-lg border border-gray-800 overflow-hidden shadow-xl">
            <table className="w-full text-left border-collapse text-sm">
              <thead>
                <tr className="bg-[#2d2d2d] text-gray-400 font-medium">
                  <th className="py-4 px-6">ID</th>
                  <th className="py-4 px-6">Problem</th>
                  <th className="py-4 px-6">Status</th>
                  <th className="py-4 px-6">Language</th>
                  <th className="py-4 px-6">Date</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-800">
                {submissions.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="py-10 text-center text-gray-500 italic">No submissions found.</td>
                  </tr>
                ) : (
                  submissions.map((s) => (
                    <tr key={s.id} className="hover:bg-[#252526] transition-colors group">
                      <td className="py-4 px-6 text-gray-500 font-mono text-xs">
                        {s.id.split('-')[0]}...
                      </td>
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

export default SubmissionListPage;
