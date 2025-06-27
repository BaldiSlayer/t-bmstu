import React from "react";
import { useParams, Link } from "react-router-dom";
import { ArrowLeft, CheckCircle, XCircle, Clock, AlertTriangle, FileText, Calendar, Code, Zap, HardDrive } from "lucide-react";
import { Submission } from "../types/task";

const SubmissionPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();

  // Моковые данные для демонстрации
  const submission: Submission = {
    id: id || "1",
    language: "C++",
    verdict: "Accepted",
    test: "1-10",
    executionTime: 15,
    memoryUsed: 1024,
    submittedAt: "2024-01-15 14:30:00",
    code: `#include <iostream>
#include <vector>
using namespace std;

int main() {
    int n;
    cin >> n;
    
    long long sum = 0;
    for (int i = 1; i <= n; i++) {
        sum += i;
    }
    
    cout << sum << endl;
    return 0;
}`
  };

  const getVerdictIcon = (verdict: Submission['verdict']) => {
    switch (verdict) {
      case 'Accepted':
        return <CheckCircle className="w-6 h-6 text-green-500" />;
      case 'Wrong Answer':
      case 'Time Limit Exceeded':
      case 'Memory Limit Exceeded':
      case 'Runtime Error':
        return <XCircle className="w-6 h-6 text-red-500" />;
      case 'Compilation Error':
        return <AlertTriangle className="w-6 h-6 text-orange-500" />;
      case 'Waiting':
      case 'Compiling':
        return <Clock className="w-6 h-6 text-blue-500" />;
      default:
        return <AlertTriangle className="w-6 h-6 text-gray-500" />;
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

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <Link 
            to="/task/1" 
            className="flex items-center text-secondary-600 hover:text-secondary-900 dark:text-secondary-400 dark:hover:text-white transition-colors"
          >
            <ArrowLeft className="w-5 h-5 mr-2" />
            Назад к задаче
          </Link>
          <div>
            <h1 className="text-3xl font-bold text-secondary-900 dark:text-white">Посылка #{id}</h1>
            <p className="text-secondary-600 mt-1 dark:text-secondary-400">Детальная информация</p>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Submission Info */}
        <div className="lg:col-span-1 space-y-6">
          {/* Status Card */}
          <div className="card">
            <div className="card-header">
              <h2 className="card-title dark:text-white">Статус</h2>
            </div>
            <div className="card-content">
              <div className="flex items-center space-x-3">
                {getVerdictIcon(submission.verdict)}
                <span className={`px-3 py-1 rounded-full text-sm font-medium ${getVerdictColor(submission.verdict)}`}>
                  {submission.verdict}
                </span>
              </div>
            </div>
          </div>

          {/* Details Card */}
          <div className="card">
            <div className="card-header">
              <h2 className="card-title dark:text-white">Детали</h2>
            </div>
            <div className="card-content space-y-4">
              <div className="flex items-center justify-between">
                <span className="text-secondary-600 dark:text-secondary-400">Язык:</span>
                <span className="font-medium dark:text-white">{submission.language}</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-secondary-600 dark:text-secondary-400">Тест:</span>
                <span className="font-medium dark:text-white">{submission.test}</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-secondary-600 dark:text-secondary-400">Время:</span>
                <div className="flex items-center space-x-1">
                  <Zap className="w-4 h-4 text-yellow-500" />
                  <span className="font-medium dark:text-white">{submission.executionTime} мс</span>
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-secondary-600 dark:text-secondary-400">Память:</span>
                <div className="flex items-center space-x-1">
                  <HardDrive className="w-4 h-4 text-blue-500" />
                  <span className="font-medium dark:text-white">{submission.memoryUsed} КБ</span>
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-secondary-600 dark:text-secondary-400">Дата:</span>
                <div className="flex items-center space-x-1">
                  <Calendar className="w-4 h-4 text-green-500" />
                  <span className="font-medium dark:text-white">{submission.submittedAt}</span>
                </div>
              </div>
            </div>
          </div>

          {/* Actions Card */}
          <div className="card">
            <div className="card-header">
              <h2 className="card-title dark:text-white">Действия</h2>
            </div>
            <div className="card-content space-y-3">
              <button className="w-full btn-outline btn-md">
                <FileText className="w-4 h-4 mr-2" />
                Скачать код
              </button>
              <button className="w-full btn-outline btn-md">
                <Code className="w-4 h-4 mr-2" />
                Поделиться
              </button>
            </div>
          </div>
        </div>

        {/* Code Section */}
        <div className="lg:col-span-2">
          <div className="card">
            <div className="card-header">
              <div className="flex items-center justify-between">
                <h2 className="card-title dark:text-white">Код решения</h2>
                <div className="flex items-center space-x-2">
                  <span className="text-sm text-secondary-600 dark:text-secondary-400">
                    {submission.language}
                  </span>
                </div>
              </div>
            </div>
            <div className="card-content">
              <div className="bg-secondary-900 rounded-lg p-4 overflow-x-auto">
                <pre className="text-secondary-100 font-mono text-sm whitespace-pre-wrap">
                  {submission.code}
                </pre>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SubmissionPage; 