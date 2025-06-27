import React from "react";
import { Link, useLocation } from "react-router-dom";
import { Code, Trophy, User, LogOut, Calendar } from "lucide-react";
import clsx from "clsx";
import { useAuth } from "../contexts/AuthContext";
import ThemeToggle from "./ThemeToggle";

const Navbar: React.FC = () => {
  const location = useLocation();
  const { user, logout, canCreateContest } = useAuth();
  
  const isActive = (path: string) => {
    return location.pathname === path;
  };

  const handleLogout = () => {
    logout();
  };

  return (
    <nav className="bg-white border-b border-secondary-200 shadow-sm dark:bg-secondary-900 dark:border-secondary-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex items-center">
            <Link to="/tasks" className="flex items-center space-x-2">
              <div className="flex items-center justify-center w-8 h-8 bg-gradient-to-br from-primary-500 to-primary-600 rounded-lg">
                <Code className="w-5 h-5 text-white" />
              </div>
              <span className="text-xl font-bold text-secondary-900 dark:text-white">T-BMSTU</span>
            </Link>
          </div>
          
          <div className="flex items-center space-x-1">
            <Link
              to="/tasks"
              className={clsx(
                "flex items-center px-3 py-2 rounded-md text-sm font-medium transition-colors",
                isActive("/tasks")
                  ? "bg-primary-100 text-primary-700 dark:bg-primary-900 dark:text-primary-300"
                  : "text-secondary-600 hover:text-secondary-900 hover:bg-secondary-50 dark:text-secondary-300 dark:hover:text-white dark:hover:bg-secondary-800"
              )}
            >
              <Trophy className="w-4 h-4 mr-2" />
              Задачи
            </Link>
            
            <Link
              to="/contests"
              className={clsx(
                "flex items-center px-3 py-2 rounded-md text-sm font-medium transition-colors",
                isActive("/contests")
                  ? "bg-primary-100 text-primary-700 dark:bg-primary-900 dark:text-primary-300"
                  : "text-secondary-600 hover:text-secondary-900 hover:bg-secondary-50 dark:text-secondary-300 dark:hover:text-white dark:hover:bg-secondary-800"
              )}
            >
              <Calendar className="w-4 h-4 mr-2" />
              Контесты
            </Link>
            
            <ThemeToggle />
            
            {user ? (
              <div className="flex items-center space-x-2 ml-4">
                <div className="flex items-center space-x-2 px-3 py-2 rounded-md text-sm text-secondary-600 dark:text-secondary-300">
                  <User className="w-4 h-4" />
                  <span className="dark:text-white">{user.name}</span>
                  <span className="px-2 py-1 bg-secondary-100 text-secondary-700 text-xs rounded-full dark:bg-secondary-700 dark:text-secondary-300">
                    {user.role === 'admin' ? 'Админ' : user.role === 'teacher' ? 'Учитель' : 'Пользователь'}
                  </span>
                </div>
                <button 
                  onClick={handleLogout}
                  className="flex items-center px-3 py-2 rounded-md text-sm font-medium text-secondary-600 hover:text-secondary-900 hover:bg-secondary-50 transition-colors dark:text-secondary-300 dark:hover:text-white dark:hover:bg-secondary-800"
                >
                  <LogOut className="w-4 h-4 mr-2" />
                  Выйти
                </button>
              </div>
            ) : (
              <Link
                to="/login"
                className={clsx(
                  "flex items-center px-3 py-2 rounded-md text-sm font-medium transition-colors",
                  isActive("/login")
                    ? "bg-primary-100 text-primary-700 dark:bg-primary-900 dark:text-primary-300"
                    : "text-secondary-600 hover:text-secondary-900 hover:bg-secondary-50 dark:text-secondary-300 dark:hover:text-white dark:hover:bg-secondary-800"
                )}
              >
                <User className="w-4 h-4 mr-2" />
                Войти
              </Link>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar; 