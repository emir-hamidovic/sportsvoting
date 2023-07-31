import { useMemo } from "react";


export interface APIResponse {
    playerid: string,
    name: string
    stats: Stats,
    advstats: AdvStats
}
  
export type FlattenedAPIResponse = {playerid: string, name: string} & Stats & AdvStats & {[key: string]: string | number};

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
}

export function applyPagination(documents: FlattenedAPIResponse[], page: number, rowsPerPage: number): FlattenedAPIResponse[] {
    return documents.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage);
}

export const useCustomers = (data: FlattenedAPIResponse[], page: number, rowsPerPage: number): FlattenedAPIResponse[] => {
    return useMemo(
        () => {
        return applyPagination(data, page, rowsPerPage);
        },
        [data, page, rowsPerPage]
    );
};

export const useCustomerIds = (customers: FlattenedAPIResponse[]): string[] => {
return useMemo(
    () => {
    return customers.map((customer) => customer.playerid);
    },
    [customers]
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