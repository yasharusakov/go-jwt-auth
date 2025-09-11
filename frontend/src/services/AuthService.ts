import $api from '../http'
import {AxiosResponse} from 'axios'
import {AuthResponse} from '../types/responses/AuthResponse'

export default class AuthService {
    static async login(email: string, password: string): Promise<AxiosResponse<AuthResponse>> {
        return $api.post<AuthResponse>('/auth/login', {email, password})
    }

    static async registration(email: string, password: string): Promise<AxiosResponse<AuthResponse>> {
        return $api.post<AuthResponse>('/auth/register', {email, password})
    }

    static async logout(): Promise<void> {
        return $api.post('/auth/logout')
    }

}