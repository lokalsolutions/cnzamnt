import type { Artwork, CreateArtworkInput, User } from './types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '';
const DEV_USER_ID = import.meta.env.VITE_DEV_USER_ID ?? localStorage.getItem('cnzamnt.devUserId') ?? '1';

function headers() {
  return {
    'Content-Type': 'application/json',
    'X-Dev-User-Id': DEV_USER_ID,
  };
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers: {
      ...headers(),
      ...options.headers,
    },
  });

  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || `Request failed: ${response.status}`);
  }

  return response.json() as Promise<T>;
}

export function getMe() {
  return request<User>('/api/users/me');
}

export function getArtworks() {
  return request<Artwork[]>('/api/artworks');
}

export function createArtwork(input: CreateArtworkInput) {
  return request<Artwork>('/api/artworks', {
    method: 'POST',
    body: JSON.stringify(input),
  });
}
