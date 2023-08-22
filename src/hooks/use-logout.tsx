import useAuth from "./use-auth";
import axios from "axios";

export const useLogout = () => {
    const { setAuth } = useAuth();

    const logout = async () => {
        setAuth({user: "", pwd: "", accessToken: ""});
        try {
            const response = await axios('http://localhost:8080/logout', {
                withCredentials: true
            });
        } catch (err) {
            console.error(err);
        }
    }

    return logout;
}