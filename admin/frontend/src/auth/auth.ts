import { readable, Readable } from "svelte/store";

// types
type HeaderObject = { [key: string]: string };
type AuthHeaders = {
  Authorization: string;
};

// util
function encodeBase64(data: string): string {
  const buffer = Buffer.from(data);
  return buffer.toString("base64");
}

// authentication state handling
let setAuthenticated: Function;
export const authenticated: Readable<boolean> = readable(false, (set) => {
  setAuthenticated = set;
  return () => {};
});

// auth headers
let authHeaders: AuthHeaders = { Authorization: "" };

// authentication logic
export async function authenticate(username: string, password: string) {
  const TMP_AUTH_HEADERS: AuthHeaders = {
    Authorization: `Basic ${btoa(unescape(encodeURIComponent(username + ":" + password)))}`,
  };

  const res = await fetch("/auth", {
    headers: TMP_AUTH_HEADERS,
  });

  if (res.ok) {
    authHeaders = TMP_AUTH_HEADERS;
    setAuthenticated(true);
    return true;
  }

  return false;
}

// utility for other components
export function ApplyAuth(headers: HeaderObject): HeaderObject {
  return { ...headers, ...authHeaders };
}
