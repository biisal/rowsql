import { client } from "@/client/client.gen";

const TOKEN_STORAGE_KEY = "rowsql-session-token";
const TOKEN_QUERY_PARAM = "token";

function readTokenFromUrl(): string | null {
  const url = new URL(window.location.href);
  const token = url.searchParams.get(TOKEN_QUERY_PARAM);
  if (!token) return null;

  url.searchParams.delete(TOKEN_QUERY_PARAM);
  window.history.replaceState(null, "", url.pathname + url.search + url.hash);
  return token;
}

function loadSessionToken(): string | null {
  const tokenFromUrl = readTokenFromUrl();
  if (tokenFromUrl) {
    sessionStorage.setItem(TOKEN_STORAGE_KEY, tokenFromUrl);
    return tokenFromUrl;
  }
  return sessionStorage.getItem(TOKEN_STORAGE_KEY);
}

export function installSessionToken() {
  const token = loadSessionToken();
  if (!token) return;
  client.setConfig({ headers: { Authorization: `Bearer ${token}` } });
}
