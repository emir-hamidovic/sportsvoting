import { useCallback, useEffect, useState } from 'react';
import { PlayersTable } from './PlayersTable';
import { useSelection } from '../hooks/use-selection';
import { APIResponse, usePlayerIds } from '../utils/api-response';
import { useParams } from 'react-router-dom';
import axiosInstance from '../utils/axios-instance';

interface TableDataProps {
	endpoint: string;
}

const convertToStringArray = (obj: any, parentKey = ''): string[] => {
	return Object.entries(obj).flatMap(([key, value]) => {
		const newKey = parentKey ? `${parentKey}.${key}` : key;
	
		if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
			return convertToStringArray(value, newKey);
		} else {
			return newKey;
		}
	});
};

export default function TableData ({ endpoint }: TableDataProps) {
	const [data, setData] = useState<APIResponse[]>([]);
	const [page, setPage] = useState<number>(0);
	const { pollId } = useParams();
	const [rowsPerPage, setRowsPerPage] = useState<number>(25);
	const players = data;
	const playerIds = usePlayerIds(players);
	const playersSelection = useSelection(playerIds);
	const [resKeys, setResKeys] = useState<string[]>([]);

	const handlePageChange = useCallback(
		(event: React.MouseEvent<HTMLButtonElement> | null, value: React.SetStateAction<number>) => {
			setPage(value);
		},
		[]
	);

	const handleRowsPerPageChange = useCallback(
		(event: React.ChangeEvent<HTMLInputElement>) => {
			setRowsPerPage(Number(event.target.value));
		},
		[]
	);

	const fetchData = useCallback(async (endpoint: string) => {
		try {
			const response = await axiosInstance.get<APIResponse[]>(endpoint + "/" + pollId);
			let playerInfo = response.data;

			if (playerInfo.length > 0) {
				console.log(playerInfo[0]);
				setResKeys(convertToStringArray(playerInfo[0]));
			}

			setData(playerInfo);
		} catch (error) {
			console.error('Error fetching data:', error);
		}
	}, [pollId]);

	useEffect(() => {
		fetchData(endpoint);
	}, [fetchData, endpoint]);

	return (
		<div>
			<PlayersTable
				count={data.length}
				items={players}
				onDeselectOne={playersSelection.handleDeselectOne}
				onPageChange={handlePageChange}
				onRowsPerPageChange={handleRowsPerPageChange}
				onSelectOne={playersSelection.handleSelectOne}
				page={page}
				rowsPerPage={rowsPerPage}
				selected={playersSelection.selected}
				columns={resKeys} />
		</div>
	);
};

