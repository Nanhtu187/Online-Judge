import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { Plus, Trash2, Save, ArrowLeft } from 'lucide-react';
import { getProblem, listTestCases, upsertProblem, upsertTestCases } from '../api/api';
import type { TestCase } from '../types';

const UpsertProblemPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const isEditing = Boolean(id);

  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [inputFormat, setInputFormat] = useState('');
  const [outputFormat, setOutputFormat] = useState('');
  const [timeLimit, setTimeLimit] = useState(1000);
  const [memoryLimit, setMemoryLimit] = useState(256);
  const [testCases, setTestCases] = useState<TestCase[]>([]);
  
  const [loading, setLoading] = useState(isEditing);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isEditing && id) {
      Promise.all([
        getProblem(id),
        listTestCases(id)
      ])
        .then(([problemData, testCasesData]) => {
          setTitle(problemData.title);
          setContent(problemData.content);
          setInputFormat(problemData.input_format);
          setOutputFormat(problemData.output_format);
          setTimeLimit(problemData.time_limit);
          setMemoryLimit(problemData.memory_limit);
          setTestCases(testCasesData);
        })
        .catch(err => {
          console.error("Failed to load problem details", err);
          setError("Failed to load problem details.");
        })
        .finally(() => {
          setLoading(false);
        });
    }
  }, [id, isEditing]);

  const handleAddTestCase = () => {
    setTestCases([...testCases, { input: '', expected_output: '' }]);
  };

  const handleRemoveTestCase = (index: number) => {
    const newTestCases = [...testCases];
    newTestCases.splice(index, 1);
    setTestCases(newTestCases);
  };

  const handleTestCaseChange = (index: number, field: keyof TestCase, value: string) => {
    const newTestCases = [...testCases];
    newTestCases[index] = { ...newTestCases[index], [field]: value };
    setTestCases(newTestCases);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setError(null);

    try {
      // 1. Upsert Problem
      const problemId = await upsertProblem({
        id: isEditing ? id : undefined,
        title,
        content,
        input_format: inputFormat,
        output_format: outputFormat,
        time_limit: timeLimit,
        memory_limit: memoryLimit
      });

      // 2. Upsert Test Cases
      await upsertTestCases(problemId, testCases);

      // Redirect back to problem list
      navigate('/');
    } catch (err: any) {
      console.error("Failed to save problem", err);
      setError(err?.response?.data?.message || "An error occurred while saving.");
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-[#1a1a1a] flex justify-center py-20">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[#1a1a1a] text-gray-200">
      {/* Header */}
      <div className="flex items-center px-8 h-14 bg-[#252526] border-b border-gray-700">
        <Link to="/" className="text-gray-400 hover:text-white mr-4 flex items-center">
          <ArrowLeft size={18} className="mr-1" /> Back
        </Link>
        <div className="text-white font-bold text-lg">
          {isEditing ? 'Edit Problem' : 'Create New Problem'}
        </div>
      </div>

      <div className="max-w-4xl mx-auto py-8 px-6">
        {error && (
          <div className="bg-red-500/10 border border-red-500/50 text-red-500 px-4 py-3 rounded mb-6">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-8">
          {/* Problem Details Section */}
          <div className="bg-[#1e1e1e] p-6 rounded-lg border border-gray-800 space-y-6">
            <h2 className="text-xl font-semibold text-white mb-4">Problem Details</h2>
            
            <div>
              <label className="block text-sm font-medium text-gray-400 mb-2">Title</label>
              <input
                type="text"
                required
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="w-full bg-[#2d2d2d] border border-gray-700 rounded px-4 py-2 text-white focus:outline-none focus:border-orange-500"
                placeholder="e.g., Two Sum"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-400 mb-2">Description</label>
              <textarea
                required
                rows={6}
                value={content}
                onChange={(e) => setContent(e.target.value)}
                className="w-full bg-[#2d2d2d] border border-gray-700 rounded px-4 py-2 text-white focus:outline-none focus:border-orange-500 font-mono text-sm"
                placeholder="Problem description in Markdown..."
              />
            </div>

            <div className="grid grid-cols-2 gap-6">
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-2">Input Format</label>
                <textarea
                  rows={4}
                  value={inputFormat}
                  onChange={(e) => setInputFormat(e.target.value)}
                  className="w-full bg-[#2d2d2d] border border-gray-700 rounded px-4 py-2 text-white focus:outline-none focus:border-orange-500 font-mono text-sm"
                  placeholder="Describe the input format..."
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-2">Output Format</label>
                <textarea
                  rows={4}
                  value={outputFormat}
                  onChange={(e) => setOutputFormat(e.target.value)}
                  className="w-full bg-[#2d2d2d] border border-gray-700 rounded px-4 py-2 text-white focus:outline-none focus:border-orange-500 font-mono text-sm"
                  placeholder="Describe the output format..."
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-6">
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-2">Time Limit (ms)</label>
                <input
                  type="number"
                  required
                  value={timeLimit}
                  onChange={(e) => setTimeLimit(parseInt(e.target.value) || 0)}
                  className="w-full bg-[#2d2d2d] border border-gray-700 rounded px-4 py-2 text-white focus:outline-none focus:border-orange-500"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-2">Memory Limit (MB)</label>
                <input
                  type="number"
                  required
                  value={memoryLimit}
                  onChange={(e) => setMemoryLimit(parseInt(e.target.value) || 0)}
                  className="w-full bg-[#2d2d2d] border border-gray-700 rounded px-4 py-2 text-white focus:outline-none focus:border-orange-500"
                />
              </div>
            </div>
          </div>

          {/* Test Cases Section */}
          <div className="bg-[#1e1e1e] p-6 rounded-lg border border-gray-800">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-xl font-semibold text-white">Test Cases</h2>
              <button
                type="button"
                onClick={handleAddTestCase}
                className="flex items-center text-sm bg-[#2d2d2d] hover:bg-[#3d3d3d] text-white px-3 py-1.5 rounded border border-gray-700 transition-colors"
              >
                <Plus size={16} className="mr-1" /> Add Test Case
              </button>
            </div>

            {testCases.length === 0 ? (
              <div className="text-center py-8 text-gray-500 border border-dashed border-gray-700 rounded">
                No test cases added. Click "Add Test Case" to create one.
              </div>
            ) : (
              <div className="space-y-6">
                {testCases.map((tc, idx) => (
                  <div key={idx} className="bg-[#252526] p-4 rounded border border-gray-700 relative">
                    <button
                      type="button"
                      onClick={() => handleRemoveTestCase(idx)}
                      className="absolute top-4 right-4 text-gray-500 hover:text-red-400 transition-colors"
                      title="Remove Test Case"
                    >
                      <Trash2 size={18} />
                    </button>
                    <h3 className="text-sm font-medium text-gray-400 mb-4">Test Case #{idx + 1}</h3>
                    
                    <div className="flex items-center space-x-2 mb-4">
                      <input
                        type="checkbox"
                        id={`is-sample-${idx}`}
                        checked={tc.is_sample}
                        onChange={(e) => handleTestCaseChange(idx, 'is_sample', e.target.checked as any)}
                        className="w-4 h-4 rounded border-gray-700 bg-[#1e1e1e] text-orange-600 focus:ring-orange-500"
                      />
                      <label htmlFor={`is-sample-${idx}`} className="text-sm text-gray-400 cursor-pointer">
                        Sample Test Case
                      </label>
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <label className="block text-xs text-gray-500 mb-1">Input</label>
                        <textarea
                          required
                          rows={3}
                          value={tc.input}
                          onChange={(e) => handleTestCaseChange(idx, 'input', e.target.value)}
                          className="w-full bg-[#1e1e1e] border border-gray-800 rounded px-3 py-2 text-white focus:outline-none focus:border-orange-500 font-mono text-sm"
                        />
                      </div>
                      <div>
                        <label className="block text-xs text-gray-500 mb-1">Expected Output</label>
                        <textarea
                          required
                          rows={3}
                          value={tc.expected_output}
                          onChange={(e) => handleTestCaseChange(idx, 'expected_output', e.target.value)}
                          className="w-full bg-[#1e1e1e] border border-gray-800 rounded px-3 py-2 text-white focus:outline-none focus:border-orange-500 font-mono text-sm"
                        />
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Submit Action */}
          <div className="flex justify-end pt-4">
            <button
              type="submit"
              disabled={saving}
              className="flex items-center bg-orange-600 hover:bg-orange-500 text-white font-medium py-2 px-6 rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {saving ? (
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
              ) : (
                <Save size={18} className="mr-2" />
              )}
              {saving ? 'Saving...' : 'Save Problem'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default UpsertProblemPage;
