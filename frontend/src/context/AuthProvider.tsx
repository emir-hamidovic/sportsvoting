import { createContext, useState } from "react";

interface AuthData {
    user: string;
    pwd: string;
    accessToken: string;
    id: number;
    roles: string[];
}

interface AuthContextType {
    auth: AuthData;
    setAuth: React.Dispatch<React.SetStateAction<AuthData>>;
    persist: boolean;
    setPersist: React.Dispatch<React.SetStateAction<boolean>>;
}

const AuthContext = createContext<AuthContextType>({
    auth: {user: "", pwd: "", accessToken: "", id: 0, roles: []},
    setAuth: () => {},
    persist: false,
    setPersist: () => {},
});

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [auth, setAuth] = useState<AuthData>({user: "", pwd: "", accessToken: "", id: 0, roles: []});
    const storedPersist = localStorage.getItem("persist");
    const initialPersist = storedPersist ? JSON.parse(storedPersist) : false;
    const [persist, setPersist] = useState<boolean>(initialPersist);

    return (
        <AuthContext.Provider value={{ auth, setAuth, persist, setPersist }}>
            {children}
        </AuthContext.Provider>
    )
}

export default AuthContext;