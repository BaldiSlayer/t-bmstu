import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import "./index.css";

// Инициализация темы перед рендерингом
const initializeTheme = () => {
  try {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'dark') {
      document.documentElement.classList.add('dark');
    } else if (savedTheme === 'light') {
      document.documentElement.classList.remove('dark');
    } else {
      // Если нет сохраненной темы, проверяем системные настройки
      if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        document.documentElement.classList.add('dark');
      }
    }
  } catch (e) {
    // Игнорируем ошибки localStorage
  }
};

// Инициализируем тему
initializeTheme();

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
