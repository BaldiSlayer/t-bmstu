import React, { useState } from "react";
import { useParams } from "react-router-dom";
import { Play, RotateCcw, CheckCircle, XCircle, Clock, AlertTriangle, FileText } from "lucide-react";
import { Task, Submission, Language } from "../types/task";
import CodeMirror from '@uiw/react-codemirror';
import { cpp } from '@codemirror/lang-cpp';
import { python } from '@codemirror/lang-python';
import { java } from '@codemirror/lang-java';
import { javascript } from '@codemirror/lang-javascript';
import { oneDark } from '@codemirror/theme-one-dark';

const TaskPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [code, setCode] = useState(`#include <iostream>\n#include <vector>\nusing namespace std;\n\nint main() {\n    int n;\n    cin >> n;\n    \n    long long sum = 0;\n    for (int i = 1; i <= n; i++) {\n        sum += i;\n    }\n    \n    cout << sum << endl;\n    return 0;\n}`);
  const [selectedLanguage, setSelectedLanguage] = useState("cpp");

  // Моковые данные для демонстрации
  const languages: Language[] = [
    { value: "cpp", label: "C++", icon: "⚡" },
    { value: "python", label: "Python", icon: "🐍" },
    { value: "java", label: "Java", icon: "☕" },
    { value: "javascript", label: "JavaScript", icon: "🟨" },
  ];

  const submissions: Submission[] = [
    {
      id: "1",
      language: "C++",
      verdict: "Accepted",
      test: "1-10",
      executionTime: 15,
      memoryUsed: 1024,
      submittedAt: "2024-01-15 14:30:00"
    },
    {
      id: "2", 
      language: "Python",
      verdict: "Wrong Answer",
      test: "3",
      executionTime: 45,
      memoryUsed: 2048,
      submittedAt: "2024-01-15 14:25:00"
    },
    {
      id: "3",
      language: "C++",
      verdict: "Time Limit Exceeded",
      test: "5",
      executionTime: 2000,
      memoryUsed: 512,
      submittedAt: "2024-01-15 14:20:00"
    },
    {
      id: "4",
      language: "Java",
      verdict: "Compilation Error",
      test: "-",
      executionTime: 0,
      memoryUsed: 0,
      submittedAt: "2024-01-15 14:15:00"
    }
  ];

  const getVerdictIcon = (verdict: Submission['verdict']) => {
    switch (verdict) {
      case 'Accepted':
        return <CheckCircle className="w-4 h-4 text-green-500" />;
      case 'Wrong Answer':
      case 'Time Limit Exceeded':
      case 'Memory Limit Exceeded':
      case 'Runtime Error':
        return <XCircle className="w-4 h-4 text-red-500" />;
      case 'Compilation Error':
        return <AlertTriangle className="w-4 h-4 text-orange-500" />;
      case 'Waiting':
      case 'Compiling':
        return <Clock className="w-4 h-4 text-blue-500" />;
      default:
        return <AlertTriangle className="w-4 h-4 text-gray-500" />;
    }
  };

  const getVerdictColor = (verdict: Submission['verdict']) => {
    switch (verdict) {
      case 'Accepted':
        return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300';
      case 'Wrong Answer':
      case 'Time Limit Exceeded':
      case 'Memory Limit Exceeded':
      case 'Runtime Error':
        return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300';
      case 'Compilation Error':
        return 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300';
      case 'Waiting':
      case 'Compiling':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300';
    }
  };

  const handleSubmit = () => {
    // Здесь будет логика отправки решения
    console.log('Submitting solution:', { code, language: selectedLanguage });
  };

  const handleSubmissionClick = (submissionId: string) => {
    window.location.href = `/submission/${submissionId}`;
  };

  // Определяем расширение для CodeMirror по языку
  const getLanguageExtension = () => {
    switch (selectedLanguage) {
      case 'cpp': return cpp();
      case 'python': return python();
      case 'java': return java();
      case 'javascript': return javascript();
      default: return cpp();
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-secondary-900 dark:text-white">Задача #{id}</h1>
          <p className="text-secondary-600 mt-1 dark:text-secondary-400">Сумма чисел</p>
        </div>
      </div>

      {/* Problem Description */}
      <div className="card">
        <div className="card-header">
          <h2 className="card-title dark:text-white">Условие задачи</h2>
        </div>
        <div className="card-content prose prose-sm max-w-none dark:prose-invert">
          <p className="dark:text-secondary-300">
            Дано число N. Найдите сумму всех чисел от 1 до N включительно.
          </p>
          <h3 className="dark:text-white">Входные данные</h3>
          <p className="dark:text-secondary-300">Одно целое число N (1 ≤ N ≤ 10^9).</p>
          <h3 className="dark:text-white">Выходные данные</h3>
          <p className="dark:text-secondary-300">Одно число — сумма всех чисел от 1 до N.</p>
          <h3 className="dark:text-white">Примеры</h3>
          <div className="bg-secondary-50 rounded-lg p-4 space-y-3 dark:bg-secondary-800">
            <div>
              <strong className="dark:text-white">Входные данные:</strong>
              <pre className="bg-white p-2 rounded mt-1 dark:bg-secondary-900 dark:text-secondary-100">5</pre>
            </div>
            <div>
              <strong className="dark:text-white">Выходные данные:</strong>
              <pre className="bg-white p-2 rounded mt-1 dark:bg-secondary-900 dark:text-secondary-100">15</pre>
            </div>
          </div>
        </div>
      </div>

      {/* Solution Section */}
      <div className="card">
        <div className="card-header">
          <div className="flex items-center justify-between">
            <h2 className="card-title dark:text-white">Решение</h2>
            <div className="flex items-center space-x-2">
              <select
                value={selectedLanguage}
                onChange={(e) => setSelectedLanguage(e.target.value)}
                className="input w-32"
              >
                {languages.map(lang => (
                  <option key={lang.value} value={lang.value}>
                    {lang.icon} {lang.label}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </div>
        <div className="card-content">
          <CodeMirror
            value={code}
            height="380px"
            theme={oneDark}
            extensions={[getLanguageExtension()]}
            onChange={value => setCode(value)}
            className="rounded-lg border border-secondary-700 dark:border-secondary-600"
            basicSetup={{ lineNumbers: true, highlightActiveLine: true }}
          />
        </div>
        <div className="card-footer">
          <div className="flex items-center justify-between w-full">
            <button 
              onClick={() => setCode(`#include <iostream>\n#include <vector>\nusing namespace std;\n\nint main() {\n    int n;\n    cin >> n;\n    \n    long long sum = 0;\n    for (int i = 1; i <= n; i++) {\n        sum += i;\n    }\n    \n    cout << sum << endl;\n    return 0;\n}`)}
              className="btn-outline btn-sm"
            >
              <RotateCcw className="w-4 h-4 mr-2" />
              Сбросить
            </button>
            <button 
              onClick={handleSubmit}
              className="btn-primary btn-lg"
            >
              <Play className="w-4 h-4 mr-2" />
              Отправить
            </button>
          </div>
        </div>
      </div>

      {/* Submissions Table */}
      <div className="card">
        <div className="card-header">
          <h2 className="card-title dark:text-white">Посылки</h2>
        </div>
        <div className="card-content">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-secondary-200 dark:border-secondary-700">
                  <th className="text-left py-3 px-4 font-medium text-secondary-900 dark:text-white">ID</th>
                  <th className="text-left py-3 px-4 font-medium text-secondary-900 dark:text-white">Компилятор</th>
                  <th className="text-left py-3 px-4 font-medium text-secondary-900 dark:text-white">Вердикт</th>
                  <th className="text-left py-3 px-4 font-medium text-secondary-900 dark:text-white">Тест</th>
                  <th className="text-left py-3 px-4 font-medium text-secondary-900 dark:text-white">Время</th>
                  <th className="text-left py-3 px-4 font-medium text-secondary-900 dark:text-white">Память</th>
                  <th className="text-left py-3 px-4 font-medium text-secondary-900 dark:text-white">Код</th>
                </tr>
              </thead>
              <tbody>
                {submissions.map((submission, index) => (
                  <tr 
                    key={submission.id} 
                    className={`border-b border-secondary-100 dark:border-secondary-800 hover:bg-secondary-50 dark:hover:bg-secondary-800 cursor-pointer ${
                      index % 2 === 0 ? 'bg-white dark:bg-secondary-900' : 'bg-secondary-50 dark:bg-secondary-800'
                    }`}
                    onClick={() => handleSubmissionClick(submission.id)}
                  >
                    <td className="py-3 px-4 text-secondary-900 dark:text-white font-mono">
                      {submission.id}
                    </td>
                    <td className="py-3 px-4 text-secondary-700 dark:text-secondary-300">
                      {submission.language}
                    </td>
                    <td className="py-3 px-4">
                      <div className="flex items-center space-x-2">
                        {getVerdictIcon(submission.verdict)}
                        <span className={`px-2 py-1 rounded-full text-xs font-medium ${getVerdictColor(submission.verdict)}`}>
                          {submission.verdict}
                        </span>
                      </div>
                    </td>
                    <td className="py-3 px-4 text-secondary-700 dark:text-secondary-300">
                      {submission.test}
                    </td>
                    <td className="py-3 px-4 text-secondary-700 dark:text-secondary-300">
                      {submission.executionTime} мс
                    </td>
                    <td className="py-3 px-4 text-secondary-700 dark:text-secondary-300">
                      {submission.memoryUsed} КБ
                    </td>
                    <td className="py-3 px-4">
                      <button 
                        className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
                        onClick={(e) => {
                          e.stopPropagation();
                          handleSubmissionClick(submission.id);
                        }}
                      >
                        <FileText className="w-4 h-4" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
};

export default TaskPage;
