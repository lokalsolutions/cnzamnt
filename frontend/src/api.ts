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

  const contentType = response.headers.get('content-type') ?? '';
  const body = await response.text();

  if (!response.ok) {
    const message = body;
    throw new Error(message || `Request failed: ${response.status}`);
  }

  if (!contentType.includes('application/json')) {
    throw new Error(`Expected JSON from ${path}, but received ${contentType || 'unknown content type'}. Check that the frontend is proxying /api to the Go backend.`);
  }

  return JSON.parse(body) as T;
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
