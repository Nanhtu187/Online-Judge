import React from 'react';
import Editor from '@monaco-editor/react';

interface Props {
  code: string;
  language: string;
  onChange: (value: string | undefined) => void;
  onLanguageChange: (lang: string) => void;
}

const CodeEditor: React.FC<Props> = ({ code, language, onChange, onLanguageChange }) => {
  return (
    <div className="flex flex-col h-full bg-[#1e1e1e]">
      <div className="flex items-center justify-between px-4 py-2 border-b border-gray-700 bg-[#252526]">
        <div className="flex items-center space-x-4">
          <select
            value={language}
            onChange={(e) => onLanguageChange(e.target.value)}
            className="bg-[#3c3c3c] text-white text-sm rounded px-2 py-1 outline-none border border-transparent focus:border-blue-500"
          >
            <option value="go">Go</option>
            <option value="cpp">C++</option>
          </select>
        </div>
        <div className="text-xs text-gray-400 font-mono">
          main.{language === 'cpp' ? 'cpp' : 'go'}
        </div>
      </div>
      <div className="flex-1">
        <Editor
          height="100%"
          language={language === 'cpp' ? 'cpp' : 'go'}
          theme="vs-dark"
          value={code}
          onChange={onChange}
          options={{
            minimap: { enabled: false },
            fontSize: 14,
            lineNumbers: 'on',
            scrollBeyondLastLine: false,
            automaticLayout: true,
            tabSize: 2,
          }}
        />
      </div>
    </div>
  );
};

export default CodeEditor;
