import React from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";

import { AuthProvider } from "./contexts/AuthContext";
import { ThemeProvider } from "./contexts/ThemeContext";
import Navbar from "./components/Navbar";
import LoginPage from "./pages/LoginPage";
import TaskListPage from "./pages/TaskListPage";
import TaskPage from "./pages/TaskPage";
import SubmissionPage from "./pages/SubmissionPage";
import ContestListPage from "./pages/ContestListPage";

const App: React.FC = () => {
  return (
    <ThemeProvider>
      <AuthProvider>
        <BrowserRouter>
          <div className="min-h-screen bg-white dark:bg-gray-900 transition-colors duration-200">
            <Navbar />
            <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 bg-gray-50 dark:bg-gray-900 min-h-screen">
              <Routes>
                <Route path="/" element={<Navigate to="/tasks" replace />} />
                <Route path="/login" element={<LoginPage />} />
                <Route path="/tasks" element={<TaskListPage />} />
                <Route path="/task/:id" element={<TaskPage />} />
                <Route path="/submission/:id" element={<SubmissionPage />} />
                <Route path="/contests" element={<ContestListPage />} />
              </Routes>
            </main>
          </div>
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  );
};

export default App;
