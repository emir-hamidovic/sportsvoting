import { useCallback, useEffect, useState } from 'react';
import { PlayersTable } from './PlayersTable';
import { useSelection } from '../hooks/use-selection';
import { APIResponse, FlattenedAPIResponse, flattenObject, usePlayerIds, usePlayers } from '../utils/api-response';
import { useParams } from 'react-router-dom';
import axiosInstance from '../utils/axios-instance';

interface TableDataProps {
	endpoint: string;
}

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
			// const transformedData: FlattenedAPIResponse[] = response.data.map((item) => ({
			// 	playerid: item.playerid,
			// 	name: item.name,
			// 	g: item.stats.g  ?? 0,
			// 	mpg: item.stats.mpg ?? '',
			// 	ppg: item.stats.ppg ?? '',
			// 	apg: item.stats.apg ?? '',
			// 	rpg: item.stats.rpg ?? '',
			// 	spg: item.stats.spg ?? '',
			// 	bpg: item.stats.bpg ?? '',
			// 	topg: item.stats.topg ?? '',
			// 	fgpct: item.stats.fgpct ?? '',
			// 	threefgpct: item.stats.threefgpct ?? '',
			// 	ftpct: item.stats.ftpct ?? '',
			// 	per: item.advstats.per ?? '',
			// 	ows: item.advstats.ows ?? '',
			// 	dws: item.advstats.dws ?? '',
			// 	ws: item.advstats.ws ?? '',
			// 	obpm: item.advstats.obpm ?? '',
			// 	dbpm: item.advstats.dbpm ?? '',
			// 	bpm: item.advstats.bpm ?? '',
			// 	vorp: item.advstats.vorp ?? '',
			// 	offrtg: item.advstats.offrtg ?? '',
			// 	defrtg: item.advstats.defrtg ?? '',
			// }));

			if (playerInfo.length > 0) {
				const firstObjectKeys = flattenObject(playerInfo[0]);
				setResKeys(firstObjectKeys);
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

