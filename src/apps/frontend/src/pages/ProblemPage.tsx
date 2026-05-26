import React, { useState, useEffect, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import Split from 'react-split';
import Navbar from '../components/Navbar';
import ProblemDescription from '../components/ProblemDescription';
import CodeEditor from '../components/CodeEditor';
import ConsolePanel from '../components/ConsolePanel';
import { getProblem, submitCode, getSubmissionResultDetail, listSubmissions } from '../api/api';
import type { Problem, TestCaseResult, Submission } from '../types';
import { SubmissionStatus, SubmissionType } from '../types';
import { useAuth } from '../AuthContext';

const DEFAULT_CODE: Record<string, string> = {
  go: `package main\n\nimport "fmt"\n\nfunc main() {\n\tfmt.Println("Hello, World!")\n}`,
  cpp: `#include <iostream>\n\nint main() {\n\tstd::cout << "Hello, World!" << std::endl;\n\treturn 0;\n}`,
};

const ProblemPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const { user } = useAuth();
  const [problem, setProblem] = useState<Problem | null>(null);
  const [code, setCode] = useState(DEFAULT_CODE['go']);
  const [language, setLanguage] = useState('go');
  const [results, setResults] = useState<TestCaseResult[] | null>(null);
  const [pastSubmissions, setPastSubmissions] = useState<Submission[] | null>(null);
  const [submissionStatus, setSubmissionStatus] = useState<SubmissionStatus | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isConsoleOpen, setIsConsoleOpen] = useState(true);

  const fetchPastSubmissions = useCallback(async () => {
    if (id && user) {
      try {
        const data = await listSubmissions({ problem_id: id, user_id: user.id });
        setPastSubmissions(data);
      } catch (err) {
        console.error('Failed to fetch past submissions:', err);
      }
    }
  }, [id, user]);

  useEffect(() => {
    if (id) {
      getProblem(id).then(setProblem).catch(console.error);
      fetchPastSubmissions();
    }
  }, [id, fetchPastSubmissions]);

  const handleLanguageChange = (lang: string) => {
    setLanguage(lang);
    if (code === DEFAULT_CODE[language === 'go' ? 'go' : 'cpp']) {
        setCode(DEFAULT_CODE[lang]);
    }
  };

  const pollResults = useCallback(async (submissionId: string) => {
    const poll = async () => {
      try {
        const data = await getSubmissionResultDetail(submissionId);
        setSubmissionStatus(data.status);

        // Update results UI
        if (data.test_case_results) {
          setResults(data.test_case_results);
        }
        
        // Stop polling if we reached a terminal status
        const isTerminal = data.status !== SubmissionStatus.PENDING && data.status !== SubmissionStatus.RUNNING;
        
        if (isTerminal) {
          setIsSubmitting(false);
          fetchPastSubmissions();
          return;
        }

        // Continue polling
        setTimeout(poll, 1500);
      } catch (err) {
        console.error('Polling error:', err);
        setIsSubmitting(false);
      }
    };
    poll();
  }, []);

  const handleSubmit = async (type: SubmissionType) => {
    if (!id) return;
    setIsSubmitting(true);
    setResults(null);
    setSubmissionStatus(null);
    setIsConsoleOpen(true);

    try {
      const submissionId = await submitCode(id, code, language, type);
      pollResults(submissionId);
    } catch (err) {
      console.error('Submit error:', err);
      setIsSubmitting(false);
    }
  };

  return (
    <div className="flex flex-col h-screen bg-[#1a1a1a]">
      <Navbar />

      <div className="flex-1 overflow-hidden">
        <Split 
          className="flex h-full"
          sizes={[40, 60]}
          minSize={300}
          gutterSize={4}
          snapOffset={0}
        >
          {/* Left Side: Description */}
          <div className="h-full overflow-hidden border-r border-gray-700">
            <ProblemDescription problem={problem} />
          </div>

          {/* Right Side: Editor + Console */}
          <div className="flex flex-col h-full overflow-hidden">
            <div className="flex-1 overflow-hidden">
              <CodeEditor 
                code={code}
                language={language}
                onChange={(val) => setCode(val || '')}
                onLanguageChange={handleLanguageChange}
              />
            </div>
            
            <ConsolePanel 
              results={results}
              testCases={problem?.test_cases || null}
              pastSubmissions={pastSubmissions}
              submissionStatus={submissionStatus}
              isSubmitting={isSubmitting}
              onRun={() => handleSubmit(SubmissionType.TEST)}
              onSubmit={() => handleSubmit(SubmissionType.OFFICIAL)}
              isOpen={isConsoleOpen}
              setIsOpen={setIsConsoleOpen}
            />
          </div>
        </Split>
      </div>
    </div>
  );
};

export default ProblemPage;
