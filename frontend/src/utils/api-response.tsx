import { useMemo } from "react";


export interface APIResponse {
	playerid: string;
	name: string;
	stats?: Stats;
	advstats?: AdvStats;
	playoffstats?: Stats;
	playoffadvstats?: AdvStats;
	totalstats?: TotalStats;
	totalplayoffstats?: TotalStats;
	accolades?: Accolades;
	[key: string]: any;
}

interface Stats {
	g?: number;
	mpg?: string;
	ppg?: string;
	apg?: string;
	rpg?: string;
	spg?: string;
	bpg?: string;
	topg?: string;
	fgpct?: string;
	threefgpct?: string;
	ftpct?: string;
	[key: string]: any;
}

interface AdvStats {
	per?: string;
	ows?: string;
	dws?: string;
	ws?: string;
	obpm?: string;
	dbpm?: string;
	bpm?: string;
	vorp?: string;
	offrtg?: string;
	defrtg?: string;
	[key: string]: any;
}

interface TotalStats {
	total_points?: number;
	total_rebounds?: number;
	total_assists?: number;
	total_steals?: number;
	total_blocks?: number;
	[key: string]: any;
}

interface Accolades{
	allstar?: number;
	allnba?: number;
	alldefense?: number;
	championships?: number;
	dpoy?: number;
	sixman?: number;
	roy?: number;
	fmvp?: number;
	mvp?: number;
	[key: string]: any;
}

export function applyPagination(documents: APIResponse[], page: number, rowsPerPage: number): APIResponse[] {
	return documents.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage);
}

export const usePlayers = (data: APIResponse[], page: number, rowsPerPage: number): APIResponse[] => {
	return useMemo(
		() => {
			return applyPagination(data, page, rowsPerPage);
		},
		[data, page, rowsPerPage]
	);
};

export const usePlayerIds = (players: APIResponse[]): string[] => {
	return useMemo(
		() => {
			return players.map((player) => player.playerid);
		},
		[players]
	);
};

export function flattenObject(obj: object, parentKey = ''): string[] {
	let keys: string[] = [];

	for (const [key, value] of Object.entries(obj)) {
		const currentKey = parentKey ? `${parentKey}.${key}` : key;

		if (typeof value === 'object' && !Array.isArray(value) && value !== null) {
			keys = keys.concat(flattenObject(value, currentKey));
		} else {
			// Extract only the last part of the key (without parent prefix)
			const lastDotIndex = currentKey.lastIndexOf('.');
			keys.push(lastDotIndex !== -1 ? currentKey.slice(lastDotIndex + 1) : currentKey);
		}
	}

	return keys;
}