export interface Contest {
  id: string;
  name: string;
  description: string;
  startTime: string;
  endTime: string;
  duration: number; // в минутах
  tasks: ContestTask[];
  participants: ContestParticipant[];
  status: 'upcoming' | 'active' | 'finished';
  isPublic: boolean;
  createdBy: string;
  createdAt: string;
}

export interface ContestTask {
  id: string;
  name: string;
  description: string;
  difficulty: string;
  source?: string;
  tags: string[];
  timeLimit: number; // в секундах
  memoryLimit: number; // в МБ
  points: number;
}

export interface ContestParticipant {
  userId: string;
  name: string;
  email: string;
  joinedAt: string;
  score: number;
  solvedTasks: string[];
}

export interface ContestSubmission {
  id: string;
  taskId: string;
  userId: string;
  language: string;
  code: string;
  status: 'pending' | 'running' | 'accepted' | 'wrong_answer' | 'time_limit' | 'memory_limit' | 'compilation_error';
  score: number;
  time: number; // в мс
  memory: number; // в КБ
  submittedAt: string;
  testResults: TestResult[];
}

export interface TestResult {
  testId: number;
  status: 'passed' | 'failed' | 'time_limit' | 'memory_limit';
  time: number;
  memory: number;
  input?: string;
  expectedOutput?: string;
  actualOutput?: string;
} 