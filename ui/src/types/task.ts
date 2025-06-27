export interface Task {
  id: string;
  name: string;
  description: string;
  inputFormat: string;
  outputFormat: string;
  examples: TestExample[];
  source?: string;
  difficulty: 'easy' | 'medium' | 'hard';
  tags: string[];
  timeLimit: number;
  memoryLimit: number;
}

export interface TestExample {
  id: number;
  input: string;
  output: string;
}

export interface Submission {
  id: string;
  language: string;
  verdict: 'Accepted' | 'Wrong Answer' | 'Time Limit Exceeded' | 'Memory Limit Exceeded' | 'Compilation Error' | 'Runtime Error' | 'Waiting' | 'Compiling';
  test: string;
  executionTime: number;
  memoryUsed: number;
  submittedAt: string;
  code?: string;
}

export interface Language {
  value: string;
  label: string;
  icon: string;
} 