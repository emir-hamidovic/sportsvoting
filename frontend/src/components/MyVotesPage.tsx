import React, { useCallback, useEffect, useState } from 'react';
import axiosInstance from '../utils/axios-instance';
import useAuth from '../hooks/use-auth';
import { getInitials } from './PlayersTable';
import { Avatar } from '@mui/material';
import '../MyVotesPage.css';
import { Link } from 'react-router-dom';

interface MyVotesResponse {
	poll_id: number;
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
					<div className="link-wrapper">
						<Link to={`/poll/${vote.poll_id}`} className='link-no-style'>
						<Avatar src={`../../${vote.player_id}.jpg`} alt={getInitials(vote.player_name)} />

						<div className="player-info">
							<h1>{`Player Name: ${vote.player_name}`}</h1>
							<h1>{`Poll Name: ${vote.poll_name}`}</h1>
						</div>

						<Avatar src={`../${vote.poll_image}`} />
						</Link>
					</div>
				</li>
				))}
			</ul>
		</div>
	);
};

export default MyVotesPage;
