import { useLocation, Navigate, Outlet } from "react-router-dom";
import useAuth from "../hooks/use-auth";

interface RequireAuthProps {
    allowedRoles: string[];
}

function splitString(input: any): string[] {
    if (typeof input === 'string') {
      const array = input.split(',');
      return array;
    } else if (Array.isArray(input)) {
        return input;
    } else {
        return [];
    }
  }

const RequireAuth = ({ allowedRoles }: RequireAuthProps) => {
    const { auth } = useAuth();
    const location = useLocation();
    const roles = splitString(auth?.roles);
    return (
        roles.find(role => allowedRoles?.includes(role))
            ? <Outlet />
            : auth?.user ? <Navigate to="/unauthorized" state={{ from: location }} replace />
                : <Navigate to="/login" state={{ from: location }} replace />
    );
}

export default RequireAuth;