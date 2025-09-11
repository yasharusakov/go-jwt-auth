import axios from 'axios'
import {AuthResponse} from '../types/responses/AuthResponse'

export const API_URL = import.meta.env.VITE_API_URL

const $api = axios.create({
    baseURL: API_URL,
    withCredentials: true
})

$api.interceptors.request.use((config) => {
    config.headers.Authorization = `Bearer ${localStorage.getItem('token')}`
    return config
})

$api.interceptors.response.use((config) => {
    return config
}, async (error) => {
    const originalRequest = error.config
    if (error.response.status == 401 && error.config && !error.config._isRetry) {
        originalRequest._isRetry = true
        try {
            const response = await axios.get<AuthResponse>(`${API_URL}/refresh`)
            localStorage.setItem('token', response.data.access_token)
            return $api.request(originalRequest)
        } catch (e) {
            console.log('NOT AUTHORIZED')
        }
    }
    throw error
})

export default $api