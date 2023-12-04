import axios from 'axios';

const isDockerContainer = process.env.DOCKER_CONTAINER === 'true';
const isDevelopment = isDockerContainer ? false : true;

const baseURL: string = isDevelopment ? 'http://localhost:8080/api' : "/api";
const axiosInstance = axios.create({
	baseURL: baseURL,
});

export default axiosInstance;
