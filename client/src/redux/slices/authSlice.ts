import {createAsyncThunk, createSlice, Draft, PayloadAction} from '@reduxjs/toolkit'
import axios from 'axios'

import {IUser} from '../../types/IUser'
import AuthService from '../../services/AuthService'
import {AuthResponse} from '../../types/responses/AuthResponse'
import {API_URL} from '../../http'

interface IAuthState {
    user: IUser
    isAuth: boolean
    isLoading: boolean
}

const initialState: IAuthState = {
    user: {} as IUser,
    isAuth: false,
    isLoading: false
}

interface IUserProps {
    email: string
    password: string
}

export const login = createAsyncThunk(
    'api/login',
    async (user: IUserProps) => {
        const response = await AuthService.login(user.email, user.password)
        return {
            access_token: response.data.access_token,
            user: response.data.user
        }
    }
)

export const registration = createAsyncThunk(
    'api/register',
    async (user: IUserProps) => {
        const response = await AuthService.registration(user.email, user.password)
        return {
            access_token: response.data.access_token,
            user: response.data.user
        }
    }
)

export const logout = createAsyncThunk(
    'api/logout',
    async () => {
        await AuthService.logout()
    }
)

export const checkAuth = createAsyncThunk(
    'api/refresh',
    async () => {
        const response = await axios.get<AuthResponse>(`${API_URL}/refresh`, {withCredentials: true})
        return {
            access_token: response.data.access_token,
            user: response.data.user
        }
    }
)

const authSlice = createSlice({
    name: 'auth',
    initialState,
    reducers: {
        setAuth(state: Draft<IAuthState>, action: PayloadAction<boolean>) {
            state.isAuth = action.payload
        },
        setUser(state: Draft<IAuthState>, action: PayloadAction<IUser>) {
            state.user = action.payload
        },
        setLoading(state: Draft<IAuthState>, action: PayloadAction<boolean>) {
            state.isLoading = action.payload
        }
    },
    extraReducers: builder => {
        builder.addCase(login.pending, (state) => {
            state.isLoading = true
        })
        builder.addCase(login.fulfilled, (state, action) => {
            localStorage.setItem('token', action.payload.access_token)
            state.user = action.payload.user
            state.isAuth = true
            state.isLoading = false
        })
        builder.addCase(login.rejected, (state) => {
            state.isLoading = false
        })
        builder.addCase(registration.pending, (state) => {
            state.isLoading = true
        })
        builder.addCase(registration.fulfilled, (state, action) => {
            localStorage.setItem('token', action.payload.access_token)
            state.user = action.payload.user
            state.isAuth = true
            state.isLoading = false
        })
        builder.addCase(registration.rejected, (state) => {
            state.isLoading = false
        })
        builder.addCase(logout.pending, (state) => {
            state.isLoading = true
        })
        builder.addCase(logout.fulfilled, (state) => {
            localStorage.removeItem('token')
            state.user = {} as IUser
            state.isAuth = false
            state.isLoading = false
        })
        builder.addCase(logout.rejected, (state) => {
            state.isLoading = false
        })
        builder.addCase(checkAuth.pending, (state) => {
            state.isLoading = true
        })
        builder.addCase(checkAuth.fulfilled, (state, action) => {
            localStorage.setItem('token', action.payload.access_token)
            state.user = action.payload.user
            state.isAuth = true
            state.isLoading = false
        })
        builder.addCase(checkAuth.rejected, (state) => {
            state.isLoading = false
        })
    }
})

export default authSlice.reducer

export const {
    setAuth,
    setUser,
    setLoading
} = authSlice.actions