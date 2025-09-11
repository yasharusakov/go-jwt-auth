import {useEffect, useState} from 'react'
import {useActions} from '../../hooks/useActions'
import Auth from '../auth'
import UserService from '../../services/UserService'
import {IUser} from '../../types/IUser'
import {useAppSelector} from '../../hooks/useAppSelector'

const App = () => {
    const {isLoading, isAuth, user} = useAppSelector(state => state.auth)
    const {checkAuth, logout} = useActions()
    const [users, setUsers] = useState<IUser[]>([])

    useEffect(() => {
        if (localStorage.getItem('token')) {
            checkAuth()
        }
    }, [])

    async function getUsers() {
        try {
            const response = await UserService.fetchUsers()
            setUsers(response.data)
        } catch (e) {
            console.log(e)
        }
    }

    if (isLoading) {
        return <div>Loading...</div>
    }

    if (!isAuth) {
        return (
            <div>
                <Auth/>
                <button onClick={getUsers}>Get users</button>
            </div>
        )
    }

    return (
        <div>
            <h1>{isAuth ? `User is authorized ${user?.email}` : 'PLEASE AUTHORIZE'}</h1>
            <button onClick={() => logout()}>Logout</button>
            <div>
                <button onClick={getUsers}>Get users</button>
            </div>
            {users.map(user =>
                <div key={user?.email}>{user?.email}</div>
            )}
        </div>
    )
}

export default App