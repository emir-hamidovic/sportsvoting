import { useRef, useState, useEffect } from 'react';
import useAuth from '../hooks/use-auth';
import { Link, useNavigate, useLocation } from 'react-router-dom';

import axios from 'axios';

const Login = () => {
    const { setAuth, persist, setPersist } = useAuth();

    const navigate = useNavigate();
    const location = useLocation();
    const from = location.state?.from?.pathname || "/";

    const userRef: React.RefObject<HTMLInputElement> = useRef<HTMLInputElement>(null);
    const errRef: React.RefObject<HTMLInputElement> = useRef<HTMLInputElement>(null);

    const [user, setUser] = useState('');
    const [pwd, setPwd] = useState('');
    const [errMsg, setErrMsg] = useState('');

    useEffect(() => {
        if (userRef.current) {
            userRef.current.focus();
        }
    }, [])

    useEffect(() => {
        setErrMsg('');
    }, [user, pwd])

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        try {
            const basicAuthToken = btoa(`${user}:${pwd}`);
            const response = await axios.post("http://localhost:8080/login",
                "",
                {
                    headers: { 'Content-Type': 'application/json', Authorization: `Basic ${basicAuthToken}` },
                    withCredentials: true
                }
            );
            console.log(JSON.stringify(response?.data));
            const accessToken = response?.data?.accessToken;
            setAuth({ user, pwd, accessToken });
            setUser('');
            setPwd('');
            navigate(from, { replace: true });
        } catch (err: any) {
            if (!err?.response) {
                setErrMsg('No Server Response');
            } else if (err.response?.status === 400) {
                setErrMsg('Missing Username or Password');
            } else if (err.response?.status === 401) {
                setErrMsg('Unauthorized');
            } else {
                setErrMsg('Login Failed');
            }
            
            if(errRef.current) {
                errRef.current.focus();
            }
        }
    }

    const togglePersist = () => {
        setPersist(prev => !prev);
    }

    useEffect(() => {
        localStorage.setItem("persist", String(persist));
    }, [persist])

    return (

        <section>
            <p ref={errRef} className={errMsg ? "errmsg" : "offscreen"} aria-live="assertive">{errMsg}</p>
            <h1>Sign In</h1>
            <form onSubmit={handleSubmit}>
                <label htmlFor="username">Username:</label>
                <input
                    type="text"
                    id="username"
                    ref={userRef}
                    autoComplete="off"
                    onChange={(e) => setUser(e.target.value)}
                    value={user}
                    required
                />

                <label htmlFor="password">Password:</label>
                <input
                    type="password"
                    id="password"
                    onChange={(e) => setPwd(e.target.value)}
                    value={pwd}
                    required
                />
                <button>Sign In</button>
                <div className="persistCheck">
                    <input type="checkbox" id="persist" onChange={togglePersist} checked={persist}  />
                    <label htmlFor="persist">Trust This Device</label>
                </div>
            </form>
            <p>
                Need an Account?<br />
                <span className="line">
                    <Link to="/signup">Sign Up</Link>
                </span>
            </p>
        </section>

    )
}

export default Login