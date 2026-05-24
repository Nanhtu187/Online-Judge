import React, { useState, useEffect } from 'react';
import { Play, Send, ChevronUp, ChevronDown, Beaker, Terminal } from 'lucide-react';
import type { TestCaseResult, TestCase } from '../types';
import { SubmissionStatus } from '../types';

interface Props {
  results: TestCaseResult[] | null;
  testCases: TestCase[] | null;
  submissionStatus: SubmissionStatus | null;
  isSubmitting: boolean;
  onRun: () => void;
  onSubmit: () => void;
  isOpen: boolean;
  setIsOpen: (open: boolean) => void;
}

const getStatusColor = (status: SubmissionStatus | null) => {
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

const ConsolePanel: React.FC<Props> = ({ results, testCases, submissionStatus, isSubmitting, onRun, onSubmit, isOpen, setIsOpen }) => {
  const [activeTab, setActiveTab] = useState<'testcase' | 'result'>('testcase');

  // Switch to result tab automatically when submission starts
  useEffect(() => {
    if (isSubmitting) {
      setActiveTab('result');
      setIsOpen(true);
    }
  }, [isSubmitting, setIsOpen]);

  return (
    <div className={`flex flex-col bg-[#1e1e1e] border-t border-gray-700 transition-all duration-300 ${isOpen ? 'h-80' : 'h-12'}`}>
      <div className="flex items-center justify-between px-4 h-12 bg-[#252526] shrink-0 border-b border-gray-800">
        <div className="flex items-center space-x-6 h-full">
          <button 
            onClick={() => setIsOpen(!isOpen)}
            className="flex items-center space-x-2 text-sm text-gray-400 hover:text-white transition-colors mr-2"
          >
            {isOpen ? <ChevronDown size={18} /> : <ChevronUp size={18} />}
            <span>Console</span>
          </button>

          {isOpen && (
            <div className="flex space-x-4 h-full items-center">
              <button
                onClick={() => setActiveTab('testcase')}
                className={`flex items-center space-x-2 text-xs font-bold uppercase tracking-wider h-full border-b-2 transition-all ${
                  activeTab === 'testcase' ? 'text-white border-orange-500' : 'text-gray-500 border-transparent hover:text-gray-300'
                }`}
              >
                <Beaker size={14} />
                <span>Testcase</span>
              </button>
              <button
                onClick={() => setActiveTab('result')}
                className={`flex items-center space-x-2 text-xs font-bold uppercase tracking-wider h-full border-b-2 transition-all ${
                  activeTab === 'result' ? 'text-white border-orange-500' : 'text-gray-500 border-transparent hover:text-gray-300'
                }`}
              >
                <Terminal size={14} />
                <span>Result</span>
              </button>
            </div>
          )}
        </div>
        
        <div className="flex items-center space-x-2">
          <button
            onClick={onRun}
            disabled={isSubmitting}
            className="flex items-center space-x-1 px-3 py-1 bg-[#3c3c3c] hover:bg-[#4c4c4c] text-gray-200 text-sm rounded transition-colors disabled:opacity-50"
          >
            <Play size={14} fill="currentColor" />
            <span>Run</span>
          </button>
          <button
            onClick={onSubmit}
            disabled={isSubmitting}
            className="flex items-center space-x-1 px-3 py-1 bg-green-600 hover:bg-green-700 text-white text-sm rounded transition-colors disabled:opacity-50"
          >
            <Send size={14} fill="currentColor" />
            <span>Submit</span>
          </button>
        </div>
      </div>

      {isOpen && (
        <div className="flex-1 overflow-y-auto p-4 font-mono text-sm text-gray-300">
          {activeTab === 'testcase' && (
            <div className="space-y-4">
              {!testCases || testCases.length === 0 ? (
                <div className="text-gray-500 italic">No sample test cases available.</div>
              ) : (
                testCases.map((tc, idx) => (
                  <div key={idx} className="space-y-3 bg-[#2d2d2d]/30 p-3 rounded border border-gray-800">
                    <div className="text-[10px] text-gray-500 font-bold uppercase tracking-widest">Case {idx + 1}</div>
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <div className="text-[9px] text-gray-600 uppercase font-bold mb-1">Input</div>
                        <pre className="bg-[#1e1e1e] p-2 rounded text-xs border border-gray-800">{tc.input || '(empty)'}</pre>
                      </div>
                      <div>
                        <div className="text-[9px] text-gray-600 uppercase font-bold mb-1">Expected</div>
                        <pre className="bg-[#1e1e1e] p-2 rounded text-xs border border-gray-800">{tc.expected_output || '(empty)'}</pre>
                      </div>
                    </div>
                  </div>
                ))
              )}
            </div>
          )}

          {activeTab === 'result' && (
            <div>
              {!submissionStatus && !isSubmitting && (
                <div className="text-gray-500 italic">No results yet. Run or Submit your code.</div>
              )}
              
              {isSubmitting && (
                <div className="flex items-center space-x-2 text-blue-400 animate-pulse mb-4">
                  <div className="w-2 h-2 bg-blue-400 rounded-full"></div>
                  <span>Processing submission...</span>
                </div>
              )}

              {submissionStatus && submissionStatus !== SubmissionStatus.PENDING && submissionStatus !== SubmissionStatus.RUNNING && (
                <div className="mb-6 pb-4 border-b border-gray-800">
                  <div className="text-[10px] text-gray-500 uppercase font-bold mb-1 tracking-wider">Final Status</div>
                  <div className={`text-xl font-bold ${getStatusColor(submissionStatus)}`}>
                    {submissionStatus.replace('SUBMISSION_STATUS_', '').replace('_', ' ')}
                  </div>
                </div>
              )}

              {submissionStatus === SubmissionStatus.COMPILE_ERROR && (
                <div className="bg-red-900/20 border border-red-500/50 p-4 rounded text-red-400">
                  <div className="text-sm font-bold uppercase mb-2 text-red-500">Compilation Error</div>
                  <pre className="text-xs whitespace-pre-wrap font-mono">{results?.[0]?.actual_output || 'Unknown compilation error'}</pre>
                </div>
              )}

              {results && submissionStatus !== SubmissionStatus.COMPILE_ERROR && (
                <div className="space-y-4">
                  {results.map((res, idx) => (
                    <div key={idx} className="border-l-2 border-gray-700 pl-3">
                      <div className="flex items-center space-x-2 mb-1">
                        <span className="text-gray-500">Testcase {idx + 1}:</span>
                        <span className={`font-bold ${getStatusColor(res.status)}`}>
                          {res.status.replace('SUBMISSION_STATUS_', '')}
                        </span>
                      </div>
                      {res.actual_output && (
                        <div className="bg-[#2d2d2d] p-2 rounded text-xs">
                          <div className="text-gray-500 mb-1">Output:</div>
                          <pre className="whitespace-pre-wrap">{res.actual_output}</pre>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default ConsolePanel;
