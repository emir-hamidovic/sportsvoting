import axiosInstance from '../utils/axios-instance';
import useAuth from './use-auth';

const useRefreshToken = () => {
    const { setAuth } = useAuth();

    const refresh = async () => {
        const response = await axiosInstance.get('/refresh', {
            withCredentials: true
        });

        setAuth({user: response.data.user, pwd: "", accessToken: response.data.access_token, id: response.data.id, roles: response.data.roles});
        return response.data;
    }
    return refresh;
};

export default useRefreshToken;