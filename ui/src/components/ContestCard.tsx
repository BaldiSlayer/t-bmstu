import React from "react";
import { Link } from "react-router-dom";
import { Calendar, Clock, Users, Trophy, Edit, Eye } from "lucide-react";
import { Contest } from "../types/contest";
import { useAuth } from "../contexts/AuthContext";
import clsx from "clsx";

interface ContestCardProps {
  contest: Contest;
}

const ContestCard: React.FC<ContestCardProps> = ({ contest }) => {
  const { user, canManageContest } = useAuth();

  const getStatusColor = (status: Contest['status']) => {
    switch (status) {
      case 'upcoming':
        return 'bg-blue-100 text-blue-700 border-blue-200 dark:bg-blue-900 dark:text-blue-300 dark:border-blue-700';
      case 'active':
        return 'bg-green-100 text-green-700 border-green-200 dark:bg-green-900 dark:text-green-300 dark:border-green-700';
      case 'finished':
        return 'bg-gray-100 text-gray-700 border-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:border-gray-600';
      default:
        return 'bg-secondary-100 text-secondary-700 border-secondary-200 dark:bg-secondary-700 dark:text-secondary-300 dark:border-secondary-600';
    }
  };

  const getStatusText = (status: Contest['status']) => {
    switch (status) {
      case 'upcoming':
        return 'Скоро';
      case 'active':
        return 'Активен';
      case 'finished':
        return 'Завершен';
      default:
        return status;
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const formatDuration = (minutes: number) => {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return `${hours}ч ${mins}мин`;
  };

  return (
    <div className="card hover:shadow-md transition-all duration-200 hover:-translate-y-1 group">
      <div className="card-content">
        <div className="flex items-start justify-between mb-3">
          <div className="flex-1">
            <Link 
              to={`/contest/${contest.id}`}
              className="group-hover:text-primary-600 transition-colors"
            >
              <h3 className="text-lg font-semibold text-secondary-900 mb-1 group-hover:text-primary-600 transition-colors dark:text-white">
                {contest.name}
              </h3>
            </Link>
            
            <p className="text-sm text-secondary-600 mb-3 line-clamp-2 dark:text-secondary-400">
              {contest.description}
            </p>
          </div>
          
          <div className="flex items-center space-x-2 ml-4">
            <div className={clsx(
              "px-2 py-1 rounded-full text-xs font-medium border",
              getStatusColor(contest.status)
            )}>
              {getStatusText(contest.status)}
            </div>
            
            {canManageContest(contest.createdBy) && (
              <Link
                to={`/contest/${contest.id}/edit`}
                className="flex items-center justify-center w-8 h-8 rounded-full bg-secondary-100 text-secondary-600 hover:bg-primary-100 hover:text-primary-600 transition-colors dark:bg-secondary-700 dark:text-secondary-400 dark:hover:bg-primary-900 dark:hover:text-primary-300"
              >
                <Edit className="w-4 h-4" />
              </Link>
            )}
          </div>
        </div>
        
        <div className="space-y-2 mb-4">
          <div className="flex items-center space-x-2 text-sm text-secondary-600 dark:text-secondary-400">
            <Calendar className="w-4 h-4" />
            <span>Начало: {formatDate(contest.startTime)}</span>
          </div>
          
          <div className="flex items-center space-x-2 text-sm text-secondary-600 dark:text-secondary-400">
            <Clock className="w-4 h-4" />
            <span>Длительность: {formatDuration(contest.duration)}</span>
          </div>
          
          <div className="flex items-center space-x-2 text-sm text-secondary-600 dark:text-secondary-400">
            <Trophy className="w-4 h-4" />
            <span>Задач: {contest.tasks.length}</span>
          </div>
          
          <div className="flex items-center space-x-2 text-sm text-secondary-600 dark:text-secondary-400">
            <Users className="w-4 h-4" />
            <span>Участников: {contest.participants.length}</span>
          </div>
        </div>
        
        <div className="flex items-center justify-between pt-3 border-t border-secondary-100 dark:border-secondary-700">
          <div className="flex items-center space-x-2 text-sm text-secondary-500 dark:text-secondary-400">
            <span>Создан: {formatDate(contest.createdAt)}</span>
          </div>
          
          <div className="flex items-center space-x-2">
            <Link
              to={`/contest/${contest.id}`}
              className="btn-outline btn-sm"
            >
              <Eye className="w-4 h-4 mr-1" />
              Просмотр
            </Link>
            
            {contest.status === 'upcoming' && (
              <Link
                to={`/contest/${contest.id}/register`}
                className="btn-primary btn-sm"
              >
                Участвовать
              </Link>
            )}
            
            {contest.status === 'active' && (
              <Link
                to={`/contest/${contest.id}`}
                className="btn-primary btn-sm"
              >
                Войти
              </Link>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default ContestCard; 