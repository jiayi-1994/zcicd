import request from '@/utils/request'

export interface LoginParams {
  username: string
  password: string
}

export interface LoginResult {
  access_token: string
  refresh_token: string
  expires_in: number
  token_type: string
}

export const authApi = {
  login: (data: LoginParams) => request.post<any, { code: number; data: LoginResult }>('/auth/login', data),
  register: (data: { username: string; email: string; password: string }) =>
    request.post('/auth/register', data),
  getProfile: () => request.get('/auth/profile'),
  refreshToken: (refresh_token: string) => request.post('/auth/refresh', { refresh_token }),
}
