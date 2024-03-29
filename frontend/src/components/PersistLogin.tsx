import { Outlet } from "react-router-dom";
import { useState, useEffect } from "react";
import useRefreshToken from '../hooks/use-refresh-token';
import useAuth from '../hooks/use-auth';

const PersistLogin = () => {
	const [isLoading, setIsLoading] = useState(true);
	const refresh = useRefreshToken();
	const { auth, persist } = useAuth();

	useEffect(() => {
		let isMounted = true;

		const verifyRefreshToken = async () => {
			try {
				await refresh();
			}
			catch (err) {
				console.error(err);
			}
			finally {
				isMounted && setIsLoading(false);
			}
		}

		!auth?.accessToken && persist ? verifyRefreshToken() : setIsLoading(false);

		return () => {
			isMounted = false;
		};
	}, [])

	useEffect(() => {
		console.log(`isLoading: ${isLoading}`)
	}, [isLoading])

	return (
		<>
			{!persist
				? <Outlet />
				: isLoading
					? <p>Loading...</p>
					: <Outlet />
			}
		</>
	)
}

export default PersistLogin