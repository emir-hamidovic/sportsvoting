import { useCallback, useEffect, useMemo, useState } from 'react';
import axios from 'axios';
import { CustomersTable } from './CustomersTable';
import { useSelection } from '../hooks/use-selection';

interface APIResponse {
  playerid: string,
  name: string
  stats: Stats,
  advstats: AdvStats
}

export type FlattenedAPIResponse = {playerid: string, name: string} & Stats & AdvStats & {[key: string]: string | number};

interface Stats {
  g: number;
  mpg: string;
  ppg: string;
  apg: string;
  rpg: string;
  spg: string;
  bpg: string;
  topg: string;
  fgpct: string;
  threefgpct: string;
  ftpct: string;
}

interface AdvStats {
  per: string;
  ows: string;
  dws: string;
  ws: string;
  obpm: string;
  dbpm: string;
  bpm: string;
  vorp: string;
  offrtg: string;
  defrtg: string;
}

function applyPagination(documents: FlattenedAPIResponse[], page: number, rowsPerPage: number): FlattenedAPIResponse[] {
  return documents.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage);
}

const useCustomers = (data: FlattenedAPIResponse[], page: number, rowsPerPage: number): FlattenedAPIResponse[] => {
  return useMemo(
    () => {
      return applyPagination(data, page, rowsPerPage);
    },
    [page, rowsPerPage]
  );
};

const useCustomerIds = (customers: FlattenedAPIResponse[]): string[] => {
  return useMemo(
    () => {
      return customers.map((customer) => customer.playerid);
    },
    [customers]
  );
};

function flattenObject(obj: object, parentKey = ''): string[] {
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

export default function MVP () {
  const [data, setData] = useState<FlattenedAPIResponse[]>([]);

  const [page, setPage] = useState<number>(0);
  const [rowsPerPage, setRowsPerPage] = useState<number>(5);
  const customers = useCustomers(data, page, rowsPerPage);
  const customersIds = useCustomerIds(customers);
  const customersSelection = useSelection(customersIds);
  const [resKeys, setResKeys] = useState<string[]>([]); // Use useState to track resKeys

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

  const fetchData = useCallback(async () => {
    try {
      const response = await axios.get<APIResponse[]>('http://localhost:8080/mvp');
      const transformedData: FlattenedAPIResponse[] = response.data.map((item) => ({
        playerid: item.playerid,
        name: item.name,
        g: item.stats.g,
        mpg: item.stats.mpg,
        ppg: item.stats.ppg,
        apg: item.stats.apg,
        rpg: item.stats.rpg,
        spg: item.stats.spg,
        bpg: item.stats.bpg,
        topg: item.stats.topg,
        fgpct: item.stats.fgpct,
        threefgpct: item.stats.threefgpct,
        ftpct: item.stats.ftpct,
        per: item.advstats.per,
        ows: item.advstats.ows,
        dws: item.advstats.dws,
        ws: item.advstats.ws,
        obpm: item.advstats.obpm,
        dbpm: item.advstats.dbpm,
        bpm: item.advstats.bpm,
        vorp: item.advstats.vorp,
        offrtg: item.advstats.offrtg,
        defrtg: item.advstats.defrtg
      }));

      if (response.data.length > 0) {
        const firstObjectKeys = flattenObject(response.data[0]);
        setResKeys(firstObjectKeys);
      }

      setData(transformedData);
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return (
    <div>
      <CustomersTable
        count={data.length}
        items={customers}
        onDeselectAll={customersSelection.handleDeselectAll}
        onDeselectOne={customersSelection.handleDeselectOne}
        onPageChange={handlePageChange}
        onRowsPerPageChange={handleRowsPerPageChange}
        onSelectAll={customersSelection.handleSelectAll}
        onSelectOne={customersSelection.handleSelectOne}
        page={page}
        rowsPerPage={rowsPerPage}
        selected={customersSelection.selected}
        columns={resKeys} />
    </div>
  );
};

