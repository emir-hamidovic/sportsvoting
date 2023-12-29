import axios from 'axios';

const isDevelopmentEnv = process.env.IS_DEVELOPMENT === 'true';

const baseURL: string = isDevelopmentEnv ? 'http://localhost:8080/api' : "/api";

// const baseURL: string = "/api";
//const baseURL: string = 'http://localhost:8080/api';
const axiosInstance = axios.create({
	baseURL: baseURL,
});

export default axiosInstance;
