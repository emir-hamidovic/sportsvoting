import { createContext, useState } from "react";

interface AuthData {
    user: string;
    pwd: string;
    accessToken: string;
}

interface AuthContextType {
    auth: AuthData;
    setAuth: React.Dispatch<React.SetStateAction<AuthData>>;
    persist: boolean;
    setPersist: React.Dispatch<React.SetStateAction<boolean>>;
}

// Create the AuthContext with the specified type
const AuthContext = createContext<AuthContextType>({
    auth: {user: "", pwd: "", accessToken: ""},
    setAuth: () => {},
    persist: false,
    setPersist: () => {},
});

export const AuthProvider = (children: React.ReactNode) => {
    const [auth, setAuth] = useState<AuthData>({user: "", pwd: "", accessToken: ""});
    const [persist, setPersist] = useState<boolean>(JSON.parse(localStorage.getItem("persist") || "") || false);

    return (
        <AuthContext.Provider value={{ auth, setAuth, persist, setPersist }}>
            {children}
        </AuthContext.Provider>
    )
}

export default AuthContext;