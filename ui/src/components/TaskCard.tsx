import React from "react";
import { Link } from "react-router-dom";
import { ExternalLink, Clock, Target } from "lucide-react";
import clsx from "clsx";

export interface TaskCardProps {
  id: string;
  name: string;
  source?: string;
  difficulty?: string;
  description?: string;
  tags?: string[];
}

const TaskCard: React.FC<TaskCardProps> = ({ 
  id, 
  name, 
  source, 
  difficulty, 
  description,
  tags = []
}) => {
  const getDifficultyColor = (diff: string) => {
    switch (diff.toLowerCase()) {
      case 'easy':
        return 'bg-success-100 text-success-700 border-success-200 dark:bg-success-900 dark:text-success-300 dark:border-success-700';
      case 'medium':
        return 'bg-warning-100 text-warning-700 border-warning-200 dark:bg-warning-900 dark:text-warning-300 dark:border-warning-700';
      case 'hard':
        return 'bg-danger-100 text-danger-700 border-danger-200 dark:bg-danger-900 dark:text-danger-300 dark:border-danger-700';
      default:
        return 'bg-secondary-100 text-secondary-700 border-secondary-200 dark:bg-secondary-700 dark:text-secondary-300 dark:border-secondary-600';
    }
  };

  const getSourceIcon = (src: string) => {
    switch (src.toLowerCase()) {
      case 'codeforces':
        return '🔵';
      case 'timus':
        return '🟡';
      case 'atcoder':
        return '🟢';
      default:
        return '📝';
    }
  };

  return (
    <div className="card hover:shadow-md transition-all duration-200 hover:-translate-y-1 group mt-4">
      <div className="card-content">
        <div className="flex items-start justify-between mb-3">
          <div className="flex-1">
            <Link 
              to={`/task/${id}`}
              className="group-hover:text-primary-600 transition-colors"
            >
              <h3 className="text-lg font-semibold text-secondary-900 mb-1 group-hover:text-primary-600 transition-colors dark:text-white">
                {name}
              </h3>
            </Link>
            
            {description && (
              <p className="text-sm text-secondary-600 mb-3 line-clamp-2 dark:text-secondary-400">
                {description}
              </p>
            )}
          </div>
          
          <Link 
            to={`/task/${id}`}
            className="flex items-center justify-center w-8 h-8 rounded-full bg-secondary-100 text-secondary-600 hover:bg-primary-100 hover:text-primary-600 transition-colors ml-2 dark:bg-secondary-700 dark:text-secondary-400 dark:hover:bg-primary-900 dark:hover:text-primary-300"
          >
            <ExternalLink className="w-4 h-4" />
          </Link>
        </div>
        
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center space-x-3">
            {source && (
              <div className="flex items-center space-x-1 text-sm text-secondary-600 dark:text-secondary-400">
                <span>{getSourceIcon(source)}</span>
                <span>{source}</span>
              </div>
            )}
            
            {difficulty && (
              <div className={clsx(
                "px-2 py-1 rounded-full text-xs font-medium border",
                getDifficultyColor(difficulty)
              )}>
                {difficulty}
              </div>
            )}
          </div>
          
          <div className="flex items-center space-x-2 text-sm text-secondary-500 dark:text-secondary-400">
            <Clock className="w-4 h-4" />
            <span>5 мин</span>
          </div>
        </div>
        
        {tags.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {tags.map((tag, index) => (
              <span
                key={index}
                className="px-2 py-1 bg-secondary-100 text-secondary-700 text-xs rounded-md dark:bg-secondary-700 dark:text-secondary-300"
              >
                {tag}
              </span>
            ))}
          </div>
        )}
        
        <div className="flex items-center justify-between mt-4 pt-3 border-t border-secondary-100 dark:border-secondary-700">
          <div className="flex items-center space-x-4 text-sm text-secondary-500 dark:text-secondary-400">
            <div className="flex items-center space-x-1">
              <Target className="w-4 h-4" />
              <span>0 решений</span>
            </div>
          </div>
          
          <Link
            to={`/task/${id}`}
            className="btn-primary btn-sm"
          >
            Решить
          </Link>
        </div>
      </div>
    </div>
  );
};

export default TaskCard;
