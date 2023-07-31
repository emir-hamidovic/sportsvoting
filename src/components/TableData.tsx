import { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import { CustomersTable } from './CustomersTable';
import { useSelection } from '../hooks/use-selection';
import { APIResponse, FlattenedAPIResponse, flattenObject, useCustomerIds, useCustomers } from '../utils/api-response';

interface TableDataProps {
  endpoint: string;
}

export default function TableData ({ endpoint }: TableDataProps) {
  const [data, setData] = useState<FlattenedAPIResponse[]>([]);

  const [page, setPage] = useState<number>(0);
  const [rowsPerPage, setRowsPerPage] = useState<number>(5);
  const customers = useCustomers(data, page, rowsPerPage);
  const customersIds = useCustomerIds(customers);
  const customersSelection = useSelection(customersIds);
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
      const response = await axios.get<APIResponse[]>(endpoint);
      const transformedData: FlattenedAPIResponse[] = response.data.map((item) => ({
        playerid: item.playerid,
        name: item.name,
        g: item.stats.g  ?? 0,
        mpg: item.stats.mpg ?? '',
        ppg: item.stats.ppg ?? '',
        apg: item.stats.apg ?? '',
        rpg: item.stats.rpg ?? '',
        spg: item.stats.spg ?? '',
        bpg: item.stats.bpg ?? '',
        topg: item.stats.topg ?? '',
        fgpct: item.stats.fgpct ?? '',
        threefgpct: item.stats.threefgpct ?? '',
        ftpct: item.stats.ftpct ?? '',
        per: item.advstats.per ?? '',
        ows: item.advstats.ows ?? '',
        dws: item.advstats.dws ?? '',
        ws: item.advstats.ws ?? '',
        obpm: item.advstats.obpm ?? '',
        dbpm: item.advstats.dbpm ?? '',
        bpm: item.advstats.bpm ?? '',
        vorp: item.advstats.vorp ?? '',
        offrtg: item.advstats.offrtg ?? '',
        defrtg: item.advstats.defrtg ?? ''
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
    fetchData(endpoint);
  }, [fetchData, endpoint]);

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

