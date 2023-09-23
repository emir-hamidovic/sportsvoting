import useAuth from './use-auth';
import axios from 'axios';

const useRefreshToken = () => {
    const { setAuth } = useAuth();

    const refresh = async () => {
        const response = await axios.get('http://localhost:8080/refresh', {
            withCredentials: true
        });

        setAuth({user: response.data.user, pwd: "", accessToken: response.data.access_token, id: response.data.id, roles: response.data.roles});
        return response.data;
    }
    return refresh;
};

export default useRefreshToken;