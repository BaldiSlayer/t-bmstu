import React, { useState } from "react";
import { Search, Filter, Grid, List } from "lucide-react";
import TaskCard from "../components/TaskCard";

const mockTasks = [
  { 
    id: "1", 
    name: "Сумма чисел", 
    source: "Codeforces", 
    difficulty: "Easy",
    description: "Найдите сумму всех чисел от 1 до N включительно.",
    tags: ["математика", "циклы"]
  },
  { 
    id: "2", 
    name: "Поиск максимума", 
    source: "Timus", 
    difficulty: "Medium",
    description: "Найдите максимальный элемент в массиве и его позицию.",
    tags: ["массивы", "алгоритмы"]
  },
  { 
    id: "3", 
    name: "Палиндром", 
    source: "AtCoder", 
    difficulty: "Hard",
    description: "Проверьте, является ли строка палиндромом.",
    tags: ["строки", "двух указатели"]
  },
  { 
    id: "4", 
    name: "Сортировка пузырьком", 
    source: "Codeforces", 
    difficulty: "Easy",
    description: "Реализуйте алгоритм сортировки пузырьком.",
    tags: ["сортировка", "алгоритмы"]
  },
];

const TaskListPage: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedDifficulty, setSelectedDifficulty] = useState<string>("all");
  const [selectedSource, setSelectedSource] = useState<string>("all");
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");

  const filteredTasks = mockTasks.filter(task => {
    const matchesSearch = task.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         task.description?.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesDifficulty = selectedDifficulty === "all" || task.difficulty === selectedDifficulty;
    const matchesSource = selectedSource === "all" || task.source === selectedSource;
    
    return matchesSearch && matchesDifficulty && matchesSource;
  });

  const difficulties = ["all", "Easy", "Medium", "Hard"];
  const sources = ["all", "Codeforces", "Timus", "AtCoder"];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-secondary-900 dark:text-white">Задачи</h1>
          <p className="text-secondary-600 mt-1 dark:text-secondary-400">
            Решайте задачи по олимпиадному программированию
          </p>
        </div>
      </div>

      {/* Filters and Search */}
      <div className="card">
        <div className="card-content">
          <div className="flex flex-col lg:flex-row gap-4">
            {/* Search */}
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-secondary-400" />
                <input
                  type="text"
                  placeholder="Поиск задач..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="input pl-10"
                />
              </div>
            </div>

            {/* Filters */}
            <div className="flex gap-3">
              <div className="relative">
                <select
                  value={selectedDifficulty}
                  onChange={(e) => setSelectedDifficulty(e.target.value)}
                  className="input pr-8 appearance-none cursor-pointer"
                >
                  {difficulties.map(diff => (
                    <option key={diff} value={diff}>
                      {diff === "all" ? "Все сложности" : diff}
                    </option>
                  ))}
                </select>
                <Filter className="absolute right-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-secondary-400 pointer-events-none" />
              </div>

              <div className="relative">
                <select
                  value={selectedSource}
                  onChange={(e) => setSelectedSource(e.target.value)}
                  className="input pr-8 appearance-none cursor-pointer"
                >
                  {sources.map(source => (
                    <option key={source} value={source}>
                      {source === "all" ? "Все источники" : source}
                    </option>
                  ))}
                </select>
                <Filter className="absolute right-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-secondary-400 pointer-events-none" />
              </div>

              {/* View Mode Toggle */}
              <div className="flex border border-secondary-300 rounded-md dark:border-secondary-600">
                <button
                  onClick={() => setViewMode("grid")}
                  className={`p-2 ${viewMode === "grid" ? "bg-primary-100 text-primary-700 dark:bg-primary-900 dark:text-primary-300" : "text-secondary-600 hover:text-secondary-900 dark:text-secondary-400 dark:hover:text-white"}`}
                >
                  <Grid className="w-4 h-4" />
                </button>
                <button
                  onClick={() => setViewMode("list")}
                  className={`p-2 ${viewMode === "list" ? "bg-primary-100 text-primary-700 dark:bg-primary-900 dark:text-primary-300" : "text-secondary-600 hover:text-secondary-900 dark:text-secondary-400 dark:hover:text-white"}`}
                >
                  <List className="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Results Info */}
      <div className="flex items-center justify-between">
        <p className="text-sm text-secondary-600 dark:text-secondary-400">
          Найдено {filteredTasks.length} задач
        </p>
        <div className="flex items-center space-x-2 text-sm text-secondary-600 dark:text-secondary-400">
          <span>Сортировка:</span>
          <select className="bg-transparent border-none text-secondary-900 font-medium cursor-pointer dark:text-white">
            <option>По названию</option>
            <option>По сложности</option>
            <option>По источнику</option>
          </select>
        </div>
      </div>

      {/* Tasks Grid/List */}
      {filteredTasks.length > 0 ? (
        <div className={viewMode === "grid" 
          ? "grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 mt-6" 
          : "space-y-4"
        }>
          {filteredTasks.map((task) => (
            <TaskCard key={task.id} {...task} />
          ))}
        </div>
      ) : (
        <div className="text-center py-12">
          <div className="w-16 h-16 bg-secondary-100 rounded-full flex items-center justify-center mx-auto mb-4 dark:bg-secondary-800">
            <Search className="w-8 h-8 text-secondary-400" />
          </div>
          <h3 className="text-lg font-medium text-secondary-900 mb-2 dark:text-white">
            Задачи не найдены
          </h3>
          <p className="text-secondary-600 dark:text-secondary-400">
            Попробуйте изменить параметры поиска или фильтры
          </p>
        </div>
      )}
    </div>
  );
};

export default TaskListPage;
