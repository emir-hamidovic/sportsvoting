import useAuth from "./use-auth";
import axios from "axios";

const useLogout = () => {
    const { setAuth } = useAuth();

    const logout = async () => {
        setAuth({user: "", pwd: "", accessToken: ""});
        try {
            const response = await axios('/logout', {
                withCredentials: true
            });
        } catch (err) {
            console.error(err);
        }
    }

    return logout;
}

export default useLogout