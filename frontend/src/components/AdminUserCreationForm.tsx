import React, { useState } from 'react';
import { Button, TextField, Container, Paper, Typography, Snackbar } from '@mui/material';
import axiosInstance from '../utils/axios-instance';

const AdminUserCreationForm: React.FC = () => {
	const [formData, setFormData] = useState({
		username: '',
		email: '',
		password: '',
	});

	const USER_REGEX = /^[A-z][A-z0-9-_]{3,23}$/;
	const PWD_REGEX = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%]).{8,24}$/;
	const EMAIL_REGEX = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

	const [successMessage, setSuccessMessage] = useState('');
	const [errorMessage, setErrorMessage] = useState('');

	const handleCloseSnackbar = () => {
		setSuccessMessage('');
		setErrorMessage('');
	};

	const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		const { name, value } = e.target;
		setFormData({
			...formData,
			[name]: value,
		});
	};

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();

		const v1 = USER_REGEX.test(formData.username);
		const v2 = PWD_REGEX.test(formData.password);
		const v3 = EMAIL_REGEX.test(formData.email);

		if (!v1 || !v2 || !v3) {
			setErrorMessage("Invalid Entry");
			return;
		}

		try {
			const response = await axiosInstance.post('/users/admin/create', formData);

			if (response.status === 200) {
				setSuccessMessage('User creation successful');
				setFormData({
					username: '',
					email: '',
					password: '',
				});
			} else {
				setErrorMessage('User creation failed with status code:' + response.status);
			}
		} catch (error) {
			setErrorMessage('Error sending POST request:' + error);
		}
	};

	return (
		<Container maxWidth="sm">
			<Paper elevation={3} style={{ padding: '20px' }}>
				<Typography variant="h5">Create New User (Admin)</Typography>
				<form onSubmit={handleSubmit}>
					<TextField
						fullWidth
						label="Username"
						name="username"
						value={formData.username}
						onChange={handleChange}
						margin="normal"
					/>
					<TextField
						fullWidth
						label="Email"
						name="email"
						type="email"
						value={formData.email}
						onChange={handleChange}
						margin="normal"
					/>
					<TextField
						fullWidth
						label="Password"
						name="password"
						type="password"
						value={formData.password}
						onChange={handleChange}
						margin="normal"
					/>
					<Button
						type="submit"
						variant="contained"
						color="primary"
						fullWidth
						style={{ marginTop: '20px' }}
					>
						Create User
					</Button>
				</form>
			</Paper>
			<Snackbar
				open={Boolean(successMessage)}
				autoHideDuration={3000}
				onClose={handleCloseSnackbar}
				message={successMessage}
				sx={{ bottom: 100 }}
			/>

			<Snackbar
				open={Boolean(errorMessage)}
				autoHideDuration={3000}
				onClose={handleCloseSnackbar}
				message={errorMessage}
				sx={{ bottom: 100 }}
			/>
		</Container>
	);
};

export default AdminUserCreationForm;
