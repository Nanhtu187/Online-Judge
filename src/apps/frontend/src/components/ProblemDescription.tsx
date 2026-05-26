import React from 'react';
import type { Problem } from '../types';

interface Props {
  problem: Problem | null;
}

const getDifficultyColor = (diff: string) => {
  switch (diff) {
    case 'DIFFICULTY_EASY': return 'text-green-500 bg-green-500/10 border-green-500/20';
    case 'DIFFICULTY_MEDIUM': return 'text-orange-500 bg-orange-500/10 border-orange-500/20';
    case 'DIFFICULTY_HARD': return 'text-red-500 bg-red-500/10 border-red-500/20';
    default: return 'text-gray-500 bg-gray-500/10 border-gray-500/20';
  }
};

const ProblemDescription: React.FC<Props> = ({ problem }) => {
  if (!problem) return <div className="p-4">Loading problem...</div>;

  return (
    <div className="flex flex-col h-full overflow-y-auto p-6 bg-[#1e1e1e] text-gray-200">
      <div className="flex flex-wrap items-center gap-3 mb-2">
        <h1 className="text-2xl font-bold">{problem.title}</h1>
        <span className={`text-[10px] px-2 py-0.5 rounded border ${getDifficultyColor(problem.difficulty)} font-bold uppercase tracking-wider`}>
          {problem.difficulty.replace('DIFFICULTY_', '')}
        </span>
      </div>

      <div className="flex flex-wrap gap-2 mb-6">
        {problem.tags && problem.tags.map(tag => (
          <span key={tag} className="text-[10px] px-2 py-0.5 rounded bg-gray-800 text-gray-400 border border-gray-700">
            {tag}
          </span>
        ))}
      </div>
      
      <div className="flex space-x-4 mb-6 text-xs text-gray-400">
        <div className="bg-[#2d2d2d] px-2 py-1 rounded">Time Limit: <span className="text-white font-medium">{problem.time_limit}ms</span></div>
        <div className="bg-[#2d2d2d] px-2 py-1 rounded">Memory Limit: <span className="text-white font-medium">{problem.memory_limit}MB</span></div>
      </div>
      
      <div className="prose prose-invert max-w-none">
        <section className="mb-6">
          <h2 className="text-lg font-semibold border-b border-gray-700 pb-2 mb-3">Description</h2>
          <div className="whitespace-pre-wrap">{problem.content}</div>
        </section>

        <section className="mb-6">
          <h2 className="text-lg font-semibold border-b border-gray-700 pb-2 mb-3">Input Format</h2>
          <div className="bg-[#2d2d2d] p-3 rounded">{problem.input_format}</div>
        </section>

        <section className="mb-6">
          <h2 className="text-lg font-semibold border-b border-gray-700 pb-2 mb-3">Output Format</h2>
          <div className="bg-[#2d2d2d] p-3 rounded">{problem.output_format}</div>
        </section>

        {problem.test_cases && problem.test_cases.filter(tc => tc.is_sample).length > 0 && (
          <section className="mb-6">
            <h2 className="text-lg font-semibold border-b border-gray-700 pb-2 mb-3">Sample Test Cases</h2>
            <div className="space-y-4">
              {problem.test_cases.filter(tc => tc.is_sample).map((tc, idx) => (
                <div key={tc.id || idx} className="bg-[#2d2d2d] rounded overflow-hidden border border-gray-700">
                  <div className="bg-[#37373d] px-3 py-1.5 text-xs font-medium text-gray-400 border-b border-gray-700">
                    Sample #{idx + 1}
                  </div>
                  <div className="grid grid-cols-2 divide-x divide-gray-700">
                    <div className="p-3">
                      <div className="text-[10px] uppercase text-gray-500 mb-1 font-bold">Input</div>
                      <pre className="text-sm font-mono whitespace-pre-wrap">{tc.input}</pre>
                    </div>
                    <div className="p-3">
                      <div className="text-[10px] uppercase text-gray-500 mb-1 font-bold">Expected Output</div>
                      <pre className="text-sm font-mono whitespace-pre-wrap">{tc.expected_output}</pre>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </section>
        )}
      </div>
    </div>
  );
};

export default ProblemDescription;
