import { Avatar, Box, Card, Checkbox, Stack, Table, TableBody, TableCell, TableHead, TablePagination, TableRow, Typography, Button, Snackbar } from '@mui/material';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import { FlattenedAPIResponse, APIResponse, usePlayers } from '../utils/api-response';
import { useNavigate, useParams } from 'react-router-dom';
import useAuth from '../hooks/use-auth';
import { useState } from 'react';
import axiosInstance from '../utils/axios-instance';

interface PlayersTableProps {
	count: number;
	items: APIResponse[];
	onDeselectOne?: (player: string) => void;
	onPageChange?: (event: React.MouseEvent<HTMLButtonElement> | null, page: React.SetStateAction<number>) => void;
	onRowsPerPageChange?: (event: React.ChangeEvent<HTMLInputElement>) => void;
	onSelectOne?: (player: string) => void;
	page: number;
	rowsPerPage: number;
	selected: string[];
	columns: string[];
}

export const getInitials = (name = '') => name
	.replace(/\s+/, ' ')
	.split(' ')
	.slice(0, 2)
	.map((v) => v && v[0].toUpperCase())
	.join('');

const columnNameMapping: { [key: string]: string } = {
	g: 'Games',
	mpg: 'Minutes',
	ppg: 'Points',
	apg: 'Assists',
	rpg: 'Rebounds',
	spg: 'Steals',
	bpg: 'Blocks',
	topg: 'Turnovers',
	fgpct: 'FG%',
	threefgpct: '3FG%',
	ftpct: 'FT%',
	per: 'PER',
	ows: 'OWS',
	dws: 'DWS',
	ws: 'WS',
	obpm: 'OBPM',
	dbpm: 'DBPM',
	bpm: 'BPM',
	vorp: 'VORP',
	offrtg: 'ORtg',
	defrtg: 'DRtg',
};

export const PlayersTable = (props: PlayersTableProps) => {
	const {
		count = 0,
		items = [],
		onDeselectOne,
		onPageChange = () => { },
		onRowsPerPageChange,
		onSelectOne,
		page = 0,
		rowsPerPage = 0,
		selected = [],
		columns
	} = props;

	const tableFields = columns.slice(2);
	const { pollId } = useParams();
	const id = pollId ? parseInt(pollId, 10) : undefined;
	const { auth } = useAuth();
	const navigate = useNavigate();
	const [successMessage, setSuccessMessage] = useState('');
	const [errorMessage, setErrorMessage] = useState('');
	const [sortField, setSortField] = useState<string | null>(null);
	const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

	const handleCloseSnackbar = () => {
		setSuccessMessage('');
		setErrorMessage('');
	};

	const handleVote = () => {
		const selectedPlayerIds = selected.length === 1 ? selected[0] : '';

		const voteEndpoint = '/votes/players';
		const payload = { playerid: selectedPlayerIds, pollid: Number(id), userid: auth.id };
		axiosInstance
			.post(voteEndpoint, payload, {
				headers: {
					'Content-Type': 'application/json',
				},
			})
			.then(() => {
				setSuccessMessage('Vote updated successfully!');
			})
			.catch((error) => {
				setErrorMessage('An error occurred. Please try again later.');
				console.error(error);
			});
	};

	const handleSort = (field: string) => {
		if (sortField === field) {
			setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
		} else {
			setSortField(field);
			setSortOrder('desc');
		}
	};

	const getColumnValue = (obj: APIResponse, column: string): any => {
		if (column.includes('.')) {
		  return column.split('.').reduce((o, key) => o?.[key], obj);
		} else {
		  return obj.stats?.[column] || obj.advstats?.[column] || obj[column];
		}
	  };

	let sortedItems = [...items].sort((a, b) => {
		if (sortField) {
			const aValue = getColumnValue(a, sortField);
    		const bValue = getColumnValue(b, sortField);
			if (typeof aValue === 'string' && typeof bValue === 'string') {
				const comparison = aValue.localeCompare(bValue);
				return sortOrder === 'asc' ? comparison : -comparison;
			}
	
			if (typeof aValue === 'number' && typeof bValue === 'number') {
				const comparison = aValue - bValue;
				return sortOrder === 'asc' ? comparison : -comparison;
			}
		}
	
		return 0;
	});

	sortedItems = usePlayers(sortedItems, page, rowsPerPage);

	return (
		<Card>
			<Box sx={{ minWidth: 800 }}>
				<Table>
					<TableHead>
						<TableRow>
						<TableCell padding="checkbox"></TableCell>
							{columns.map((column) => (
								<TableCell key={column} onClick={() => handleSort(column)}>
									<div
										style={{
											cursor: 'pointer',
											display: 'flex',
											alignItems: 'center',
											gap: '4px',
										}}
									>
										{columnNameMapping[column]}
										{sortField === column && (
											<ArrowDropDownIcon
												style={{
													fontSize: 'inherit',
													transform: `rotate(${sortOrder === 'desc' ? '180deg' : '0deg'})`,
												}}
											/>
										)}
									</div>
								</TableCell>
							))}
						</TableRow>
					</TableHead>
					<TableBody>
						{sortedItems.map((player: APIResponse) => {
							const isSelected = selected.includes(player.playerid);

							return (
								<TableRow hover key={player.playerid} selected={isSelected}>
									<TableCell padding="checkbox">
										<Checkbox checked={isSelected}
											onChange={(event) => {
												if (event.target.checked) {
													onSelectOne?.(player.playerid);
												} else {
													onDeselectOne?.(player.playerid);
												}
											}} />
									</TableCell>
									<TableCell>
										<Stack alignItems="center" direction="row" spacing={2}>
											<Avatar src={`../${player.playerid}.jpg`} alt={getInitials(player.name)} />
											<Typography variant="subtitle2"> {player.name} </Typography>
										</Stack>
									</TableCell>
									<TableCell></TableCell>
									{tableFields.map((column) => (
										<TableCell key={column}>
											{column.includes('.') // Check if the column has nested structure
												? column.split('.').reduce((obj, key) => obj?.[key], player)
												: player.stats?.[column] || player.advstats?.[column] || player[column]}
										</TableCell>
									))}
								</TableRow>
							);
						})}
					</TableBody>
				</Table>
			</Box>
			<TablePagination
				component="div"
				count={count}
				onPageChange={onPageChange}
				onRowsPerPageChange={onRowsPerPageChange}
				page={page}
				rowsPerPage={rowsPerPage}
				rowsPerPageOptions={[25, 50, 100]}
			/>
			<Box sx={{ display: 'flex', justifyContent: 'flex-end', marginTop: 2 }}>
				<Button
					variant="contained"
					color="primary"
					onClick={() => navigate(`/results/${id}`)}
				>
					Check results
				</Button>
				<Button
					variant="contained"
					color="primary"
					disabled={selected.length !== 1}
					onClick={handleVote}
					sx={{ marginLeft: 2 }}
				>
					Vote
				</Button>
			</Box>
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
		</Card>
	);
};
