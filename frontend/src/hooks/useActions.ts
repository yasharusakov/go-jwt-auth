import {useAppDispatch} from './useAppDispath'
import {bindActionCreators} from 'redux'
import * as AuthCreators from '../redux/slices/authSlice'

export const useActions = () => {
    const dispatch = useAppDispatch()
    return bindActionCreators(
        {...AuthCreators},
        dispatch
    )
}