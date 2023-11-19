import axios from 'axios';

const isDevelopment: boolean = process.env.IS_DEVELOPMENT === 'true';

const baseURL: string = isDevelopment ? 'http://localhost:8080/api' : "/api";
const axiosInstance = axios.create({
  baseURL: baseURL,
});

export default axiosInstance;
