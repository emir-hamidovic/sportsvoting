import React, { useCallback, useEffect, useState } from 'react';
import axiosInstance from '../utils/axios-instance';
import useAuth from '../hooks/use-auth';
import { getInitials } from './PlayersTable';
import { Avatar } from '@mui/material';
import '../MyVotesPage.css';

interface MyVotesResponse {
	player_id: string;
	player_name: string;
	poll_name: string;
	poll_image: string;
}

const MyVotesPage: React.FC = () => {
	const { auth } = useAuth();
	const [userVotes, setUserVotes] = useState<MyVotesResponse[]>([]);

	const fetchUserVotes = useCallback(async () => {
		try {
			const response = await axiosInstance.get<MyVotesResponse[]>(`/votes/users/get/${auth.id}`);
			setUserVotes(response.data);
		} catch (error) {
			console.error('Error fetching user votes:', error);
		}
	}, [auth.id]);

	useEffect(() => {
		fetchUserVotes();
	}, [fetchUserVotes]);

	return (
		<div className="my-votes">
			<h1 className="text-3xl">Your Votes</h1>
			<ul>
			{userVotes.map((vote) => (
				<li key={vote.player_id} className="vote-item">
					<Avatar src={`../../${vote.player_id}.jpg`} alt={getInitials(vote.player_name)} />

					<div className="player-info">
						<p>{`Player Name: ${vote.player_name}`}</p>
						<p>{`Poll Name: ${vote.poll_name}`}</p>
					</div>

					<Avatar src={`../${vote.poll_image}`} />
				</li>
				))}
			</ul>
		</div>
	);
};

export default MyVotesPage;
