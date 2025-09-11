import {useState} from 'react'
import {useActions} from '../../hooks/useActions.js'

const Auth = () => {
    const {login, registration} = useActions()

    const [email, setEmail] = useState<string>('')
    const [password, setPassword] = useState<string>('')

    return (
        <div>
            <input
                onChange={e => setEmail(e.target.value)}
                value={email}
                type="text"
                placeholder="Email"
            />
            <input
                onChange={e => setPassword(e.target.value)}
                value={password}
                type="password"
                placeholder="Password"
            />
            <button onClick={() => login({email, password})}>
                Login
            </button>
            <button onClick={() => registration({email, password})}>
                Registration
            </button>
        </div>
    )
}

export default Auth