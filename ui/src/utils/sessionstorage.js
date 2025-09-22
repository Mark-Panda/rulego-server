export const SESSIONSTORAGE_KEYS = {
  TOKEN: 'token',
};

export function getSession(key) {
  return JSON.parse(sessionStorage.getItem(key));
}

export function setSession(key, value) {
  if (value === null || value === undefined) {
    sessionStorage.removeItem(key);
  } else {
    sessionStorage.setItem(key, JSON.stringify(value));
  }
}

export function clearSession(key) {
  sessionStorage.removeItem(key);
}

export function clearAllSession() {
  sessionStorage.clear();
}
