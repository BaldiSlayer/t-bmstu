import React, { useState } from "react";
import { Search, Filter, Grid, List, Plus, Calendar } from "lucide-react";
import ContestCard from "../components/ContestCard";
import { Contest } from "../types/contest";
import { useAuth } from "../contexts/AuthContext";

const mockContests: Contest[] = [
  {
    id: "1",
    name: "Зимний контест 2024",
    description: "Ежегодный зимний контест по алгоритмам и структурам данных",
    startTime: "2024-12-20T10:00:00Z",
    endTime: "2024-12-20T14:00:00Z",
    duration: 240,
    tasks: [
      { id: "1", name: "Сумма чисел", description: "Найдите сумму", difficulty: "Easy", tags: ["математика"], timeLimit: 1000, memoryLimit: 256, points: 100 },
      { id: "2", name: "Поиск максимума", description: "Найдите максимум", difficulty: "Medium", tags: ["алгоритмы"], timeLimit: 2000, memoryLimit: 512, points: 200 },
    ],
    participants: [
      { userId: "1", name: "Иван Иванов", email: "ivan@example.com", joinedAt: "2024-12-19T10:00:00Z", score: 0, solvedTasks: [] },
      { userId: "2", name: "Петр Петров", email: "petr@example.com", joinedAt: "2024-12-19T11:00:00Z", score: 0, solvedTasks: [] },
    ],
    status: "upcoming",
    isPublic: true,
    createdBy: "1",
    createdAt: "2024-12-15T10:00:00Z",
  },
  {
    id: "2",
    name: "Весенний турнир",
    description: "Турнир для студентов первого курса",
    startTime: "2024-03-15T09:00:00Z",
    endTime: "2024-03-15T12:00:00Z",
    duration: 180,
    tasks: [
      { id: "3", name: "Палиндром", description: "Проверьте палиндром", difficulty: "Easy", tags: ["строки"], timeLimit: 1000, memoryLimit: 256, points: 100 },
    ],
    participants: [
      { userId: "3", name: "Анна Сидорова", email: "anna@example.com", joinedAt: "2024-03-14T10:00:00Z", score: 100, solvedTasks: ["3"] },
    ],
    status: "finished",
    isPublic: true,
    createdBy: "1",
    createdAt: "2024-03-10T10:00:00Z",
  },
  {
    id: "3",
    name: "Летний марафон",
    description: "Интенсивный контест для продвинутых программистов",
    startTime: "2024-06-01T14:00:00Z",
    endTime: "2024-06-01T18:00:00Z",
    duration: 240,
    tasks: [
      { id: "4", name: "Сложная задача", description: "Очень сложная задача", difficulty: "Hard", tags: ["динамика"], timeLimit: 3000, memoryLimit: 512, points: 300 },
    ],
    participants: [],
    status: "upcoming",
    isPublic: false,
    createdBy: "1",
    createdAt: "2024-05-25T10:00:00Z",
  },
];

const ContestListPage: React.FC = () => {
  const { canCreateContest } = useAuth();
  const [searchTerm, setSearchTerm] = useState("");
  const [selectedStatus, setSelectedStatus] = useState<string>("all");
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");

  const filteredContests = mockContests.filter(contest => {
    const matchesSearch = contest.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         contest.description.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = selectedStatus === "all" || contest.status === selectedStatus;
    
    return matchesSearch && matchesStatus;
  });

  const statuses = [
    { value: "all", label: "Все статусы" },
    { value: "upcoming", label: "Скоро" },
    { value: "active", label: "Активные" },
    { value: "finished", label: "Завершенные" },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-secondary-900 dark:text-white">Контесты</h1>
          <p className="text-secondary-600 mt-1 dark:text-secondary-400">
            Соревнования по программированию
          </p>
        </div>
        {canCreateContest() && (
          <button className="btn-primary btn-lg">
            <Plus className="w-5 h-5 mr-2" />
            Создать контест
          </button>
        )}
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
                  placeholder="Поиск контестов..."
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
                  value={selectedStatus}
                  onChange={(e) => setSelectedStatus(e.target.value)}
                  className="input pr-8 appearance-none cursor-pointer"
                >
                  {statuses.map(status => (
                    <option key={status.value} value={status.value}>
                      {status.label}
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
          Найдено {filteredContests.length} контестов
        </p>
        <div className="flex items-center space-x-2 text-sm text-secondary-600 dark:text-secondary-400">
          <span>Сортировка:</span>
          <select className="bg-transparent border-none text-secondary-900 font-medium cursor-pointer dark:text-white">
            <option>По дате начала</option>
            <option>По названию</option>
            <option>По количеству участников</option>
          </select>
        </div>
      </div>

      {/* Contests Grid/List */}
      {filteredContests.length > 0 ? (
        <div className={viewMode === "grid" 
          ? "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6" 
          : "space-y-4"
        }>
          {filteredContests.map((contest) => (
            <ContestCard key={contest.id} contest={contest} />
          ))}
        </div>
      ) : (
        <div className="text-center py-12">
          <div className="w-16 h-16 bg-secondary-100 rounded-full flex items-center justify-center mx-auto mb-4 dark:bg-secondary-800">
            <Calendar className="w-8 h-8 text-secondary-400" />
          </div>
          <h3 className="text-lg font-medium text-secondary-900 mb-2 dark:text-white">
            Контесты не найдены
          </h3>
          <p className="text-secondary-600 dark:text-secondary-400">
            Попробуйте изменить параметры поиска или фильтры
          </p>
        </div>
      )}
    </div>
  );
};

export default ContestListPage; 