/**
 * 主题切换工具
 */

// 获取当前主题模式
export function getThemeMode() {
    // 优先从本地存储获取
    const savedMode = localStorage.getItem('theme-mode');
    if (savedMode) {
        return savedMode;
    }

    // 检查系统偏好
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        return 'dark';
    }

    return 'light';
}

// 设置主题模式
export function setThemeMode(mode) {
    if (mode === 'dark') {
        document.documentElement.classList.add('dark');
    } else {
        document.documentElement.classList.remove('dark');
    }

    // 保存到本地存储
    localStorage.setItem('theme-mode', mode);
}

// 切换主题模式
export function toggleThemeMode() {
    const currentMode = getThemeMode();
    const newMode = currentMode === 'dark' ? 'light' : 'dark';
    setThemeMode(newMode);
    return newMode;
}

// 初始化主题
export function initTheme() {
    const mode = getThemeMode();
    setThemeMode(mode);
    return mode;
}