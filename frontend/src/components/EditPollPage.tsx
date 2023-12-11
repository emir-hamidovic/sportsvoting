import React, { useCallback, useEffect, useRef, useState } from 'react';
import {
  Box,
  Button,
  Container,
  Divider,
  FormControl,
  Grid,
  MenuItem,
  Paper,
  Select,
  SelectChangeEvent,
  Snackbar,
  TextField,
  Typography,
} from '@mui/material';
import { useParams } from 'react-router-dom';
import axiosInstance from '../utils/axios-instance';

interface Poll {
  id: number;
  name: string;
  description: string;
  image: string;
  selected_stats: string;
  season: string;
}

const EditPollPage: React.FC = () => {
  	const { pollId } = useParams();
	const [pollInfo, setPollInfo] = useState<Poll>({
		id: 0,
		name: '',
		description: '',
		image: '',
		selected_stats: '',
		season: '',
	});

	const [initialPollInfo, setInitialPollInfo] = useState<Poll>({
		id: 0,
		name: '',
		description: '',
		image: '',
		selected_stats: '',
		season: '',
	});
	const imageInputRef = useRef<HTMLInputElement | null>(null);
	const [statsOptions] = useState<string[]>(["All stats", "Defensive", "Sixth man", "Rookie", "GOAT stats"]);
	const [seasonOptions, setSeasonOptions] = useState<string[]>([]);
	const [selectedStats, setSelectedStats] = useState<string>("");
	const [selectedSeason, setSelectedSeason] = useState<string>("");
	const [isSeasonDisabled, setIsSeasonDisabled] = useState<boolean>(true);
	const [fetchedSeasonOptions, setFetchedSeasonOptions] = useState<string[]>([]);

	const fetchData = useCallback(async () => {
		try {
			const response = await axiosInstance.get(`/polls/get/${pollId}`);
			const pollData = response.data;
		
			setPollInfo({
				id: parseInt(pollId ?? '0', 10),
				name: pollData.name,
				description: pollData.description,
				image: pollData.image,
				selected_stats: pollData.selected_stats,
				season: pollData.season,
			});
			setInitialPollInfo({
				id: parseInt(pollId ?? '0', 10),
				name: pollData.name,
				description: pollData.description,
				image: pollData.image,
				selected_stats: pollData.selected_stats,
				season: pollData.season,
			});

			setSelectedStats(pollData.selected_stats);
			setSelectedSeason(pollData.season.toString());
			setIsSeasonDisabled(false);
			const seasonsResponse = await axiosInstance.get('/seasons/get');
			setFetchedSeasonOptions(seasonsResponse.data);

			if (pollData.selected_stats == 'GOAT stats') {
				setSeasonOptions(['All', 'Playoffs', 'Career'])
			} else {
				setSeasonOptions(seasonsResponse.data);
			}
		} catch (error) {
			console.error('Error fetching data:', error);
		}
	}, [pollId]);

	useEffect(() => {
		fetchData();
	}, [fetchData]);

	const [successMessage, setSuccessMessage] = useState('');
	const [errorMessage, setErrorMessage] = useState('');

	const handleCloseSnackbar = () => {
		setSuccessMessage('');
		setErrorMessage('');
	};

	const handleStatsChange = (
		event: SelectChangeEvent<string>
	) => {
		setPollInfo({ ...pollInfo, selected_stats: event.target.value })
		const selectedStatsType = event.target.value as string;
		setIsSeasonDisabled(selectedStatsType === '');
		if (selectedStatsType === 'GOAT stats') {
			setSeasonOptions(['All', 'Playoffs', 'Career']);
		} else {
			setSeasonOptions(fetchedSeasonOptions);
		}

		setSelectedStats(selectedStatsType);
	};

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();

		try {
			if (initialPollInfo.name !== pollInfo.name || initialPollInfo.description !== pollInfo.description || initialPollInfo.selected_stats !== pollInfo.selected_stats || initialPollInfo.season !== pollInfo.season) {
			await axiosInstance.post('/polls/update', {
				id: pollInfo.id,
				name: pollInfo.name,
				description: pollInfo.description,
				selected_stats: pollInfo.selected_stats,
				season: pollInfo.season,
				});
			}

			setSuccessMessage('Poll information updated successfully!');
			setInitialPollInfo({ ...pollInfo, image: pollInfo.image, name: pollInfo.name, description: pollInfo.description, selected_stats: pollInfo.selected_stats, season: pollInfo.season });
		} catch (error) {
		setErrorMessage('An error occurred. Please try again later.');
		}
	};

  const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
	const fileInput = e.target;

	if (fileInput && fileInput.files && fileInput.files.length > 0) {
	  const fileToUpload = fileInput.files[0];

	  try {
		const formData = new FormData();
		formData.append('pollImage', fileToUpload, `poll-${pollInfo.id}.jpg`);
		formData.append("pollId", pollInfo.id.toString());

		const response = await axiosInstance.post('/polls/image/update', formData, {
			headers: {
				'Content-Type': 'multipart/form-data',
			},
		});

		if (response.data.success) {
			setPollInfo({ ...pollInfo, image: response.data.fileName });
			setInitialPollInfo({ ...pollInfo, image: response.data.fileName });
			setSuccessMessage('Poll image updated successfully!');
		} else {
			setErrorMessage('Failed to upload poll image');
		}
	  } catch (error) {
		setErrorMessage('An error occurred while uploading the poll image.');
	  }
	}
  };

  return (
	<Container maxWidth="lg">
		<Paper elevation={3} sx={{ padding: 3 }}>
		<Grid container spacing={3}>
			<Grid item md={4} sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
			<img src={`../../${pollInfo.image}`} alt="Poll" style={{ width: '100%', maxHeight: '200px', objectFit: 'cover', marginBottom: 2 }} />
			<input type="file" accept="image/*" onChange={handleImageUpload} style={{ display: 'none' }} ref={imageInputRef} />
			<Button variant="outlined" fullWidth onClick={() => {
				if (imageInputRef.current) {
					imageInputRef.current.click();
				}
			}}>Upload new Poll Image</Button>
			</Grid>
			<Grid item md={8}>
			<Typography variant="h4" gutterBottom>
				Poll Information
			</Typography>
			<Divider />

			<form onSubmit={handleSubmit}>
				<Box>
				<TextField
					fullWidth
					label="Name"
					variant="outlined"
					value={pollInfo.name}
					onChange={(e) => setPollInfo({ ...pollInfo, name: e.target.value })}
					InputProps={{ style: { marginTop: '25px' } }}
				/>
				</Box>
				<Box>
				<TextField
					fullWidth
					label="Description"
					variant="outlined"
					value={pollInfo.description}
					onChange={(e) => setPollInfo({ ...pollInfo, description: e.target.value })}
					InputProps={{ style: { marginTop: '25px' } }}
				/>
				</Box>
				<Box sx={{marginTop: '25px'}}>
				<div className="label">Stats:</div>
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
				</Box>
				<Box sx={{marginTop: '25px'}}>
				<div className="label">Season:</div>
				<FormControl fullWidth variant="standard">
					<Select
					value={selectedSeason}
					onChange={(e) => {setPollInfo({ ...pollInfo, season: e.target.value }); setSelectedSeason(e.target.value as string);}}
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
				</Box>
				<Box mt={3}>
				<Button type="submit" variant="contained" color="primary">
					Save Changes
				</Button>
				</Box>
			</form>
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
			</Grid>
		</Grid>
		</Paper>
	</Container>
  );
};

export default EditPollPage;
