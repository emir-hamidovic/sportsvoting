import useAuth from './use-auth';
import axios from 'axios';

const useRefreshToken = () => {
    const { setAuth } = useAuth();

    const refresh = async () => {
        const response = await axios.get('http://localhost:8080/refresh', {
            withCredentials: true
        });

        setAuth(prev => {
            return {
                ...prev,
                accessToken: response.data
            }
        });
        return response.data;
    }
    return refresh;
};

export default useRefreshToken;