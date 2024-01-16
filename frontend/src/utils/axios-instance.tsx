import axios from 'axios';

const isDevelopmentEnv = process.env.NODE_ENV === 'development';
let baseURL: string = isDevelopmentEnv ? 'http://localhost:8080/api' : "/api";

const axiosInstance = axios.create({
	baseURL: baseURL,
});

export default axiosInstance;
