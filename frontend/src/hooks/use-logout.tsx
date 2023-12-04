import axiosInstance from "../utils/axios-instance";
import useAuth from "./use-auth";

export const useLogout = () => {
	const { setAuth } = useAuth();

	const logout = async () => {
		setAuth({user: "", pwd: "", accessToken: "", id: 0, roles: []});
		try {
			await axiosInstance('/logout', {
				withCredentials: true
			});
		} catch (err) {
			console.error(err);
		}
	}

	return logout;
}