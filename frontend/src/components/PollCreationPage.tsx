import React, { useCallback, useEffect, useState } from 'react';
import {
	Button,
	Container,
	Grid,
	Paper,
	TextField,
	Typography,
	FormControl,
	Select,
	MenuItem,
	SelectChangeEvent,
} from '@mui/material';
import axiosInstance from '../utils/axios-instance';
import { useNavigate } from 'react-router-dom';
import useAuth from '../hooks/use-auth';


const PollCreationPage: React.FC = () => {
	const [name, setName] = useState<string>('');
	const [description, setDescription] = useState<string>('');
	const [season, setSeason] = useState<string>('');
	const [selectedStats, setSelectedStats] = useState<string>('');
	const [statsOptions] = useState<string[]>(["All stats", "Defensive", "Sixth man", "Rookie", "GOAT stats"]);
	const [seasonOptions, setSeasonOptions] = useState<string[]>([]);
	const [fetchedSeasonOptions, setFetchedSeasonOptions] = useState<string[]>([]);
	const [selectedFile, setSelectedFile] = useState<File | null>(null);
	const [isSeasonDisabled, setIsSeasonDisabled] = useState<boolean>(true);
	const navigate = useNavigate();
	const { auth } = useAuth();

	const fetchData = useCallback(async () => {
		try {
			const response = await axiosInstance.get('/seasons/get');
			setSeasonOptions(response.data);
			setFetchedSeasonOptions(response.data);
		} catch (error) {
			console.error('Error fetching data:', error);
		}
	}, []);

	useEffect(() => {
		fetchData();
	}, [fetchData]);

	const handleNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		setName(event.target.value);
	};

	const handleDescriptionChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		setDescription(event.target.value);
	};

	const handleSeasonChange = (
		event: SelectChangeEvent<string>
	) => {
		setSeason(event.target.value as string);
	};

	const handleStatsChange = (
		event: SelectChangeEvent<string>
	) => {
		const selectedStatsType = event.target.value as string;
		setSelectedStats(selectedStatsType);

		setIsSeasonDisabled(selectedStatsType === '');
		if (selectedStatsType === 'GOAT stats') {
			setSeasonOptions(['All', 'Playoffs', 'Career']);
		} else {
			setSeasonOptions(fetchedSeasonOptions);
		}

		setSeason('');

	};

	const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
		if (event.target.files && event.target.files.length > 0) {
			setSelectedFile(event.target.files[0]);
		}
	};

	const handleSubmit = async (event: React.FormEvent) => {
		event.preventDefault();

		if (!name || selectedStats.length === 0 || !season || !selectedFile) {
			alert('Name, season, at least one stat, and a photo must be provided.');
			return;
		}

		const data = new FormData();
		data.append('name', name);
		data.append('description', description);
		data.append('season', season);
		data.append('selectedStats', selectedStats);
		data.append('photo', selectedFile);
		data.append('userid', auth.id.toString());
		try {
			await axiosInstance.post('/polls/create', data);
			alert('Poll created successfully');
			navigate(-1);
		} catch (error) {
			console.error('Error creating poll:', error);
			alert('An error occurred while creating the poll.');
		}
	};

	return (
		<Container maxWidth="lg" className="poll-creation-container">
			<Paper elevation={3} sx={{ padding: 3 }}>
				<Typography variant="h4" gutterBottom>
					Create a New Poll
				</Typography>
				<form onSubmit={handleSubmit}>
					<Grid container spacing={3}>
						<Grid item md={12}>
							<div className="label">Name:</div>
							<TextField
								fullWidth
								variant="outlined"
								value={name}
								onChange={handleNameChange}
							/>
						</Grid>
						<Grid item md={12}>
							<div className="label">Description:</div>
							<TextField
								fullWidth
								variant="outlined"
								value={description}
								onChange={handleDescriptionChange}
							/>
						</Grid>
						<Grid item md={12}>
							<div className="label">Select stats type to display:</div>
							<FormControl fullWidth variant="standard">
								<Select
									value={selectedStats}
									onChange={handleStatsChange}
									label="Select Stats type"
								>
									{statsOptions.map((stat) => (
										<MenuItem key={stat} value={stat}>
											{stat}
										</MenuItem>
									))}
								</Select>
							</FormControl>
						</Grid>
						<Grid item md={12}>
							<div className="label">Season:</div>
							<FormControl fullWidth variant="standard">
								<Select
									value={season}
									onChange={handleSeasonChange}
									label="Season"
									disabled={isSeasonDisabled}
								>
								 {seasonOptions.map((stat) => (
								<MenuItem key={stat} value={stat}>
									{stat}
								</MenuItem>
							))}
								</Select>
							</FormControl>
						</Grid>
						<Grid item md={12}>
							<div className="label">Upload new photo:</div>
							<input
								type="file"
								accept="image/*"
								onChange={handleFileChange}
							/>
						</Grid>
						<Grid item md={12}>
							<Button
								type="submit"
								variant="contained"
								color="primary"
								className="submit-button"
							>
								Create Poll
							</Button>
						</Grid>
					</Grid>
				</form>
			</Paper>
		</Container>
	);
};

export default PollCreationPage;
