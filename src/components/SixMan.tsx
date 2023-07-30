import { useEffect, useMemo, useState, useCallback } from 'react';
import axios from 'axios';
import StatsTable from './StatsTable';

export default function SixMan () {
  const [data, setData] = useState([]);

  const columns = useMemo(
    () => [
      {
        Header: "Name",
        accessor: "name",
      },
      {
        Header: "Games",
        accessor: "stats.g",
      },
      {
        Header: "Minutes",
        accessor: "stats.mpg",
      },
      {
        Header: "Points",
        accessor: "stats.ppg",
      },
      {
        Header: "Rebounds",
        accessor: "stats.rpg",
      },
      {
        Header: "Assists",
        accessor: "stats.apg",
      },
      {
        Header: "Steals",
        accessor: "stats.spg",
      },
      {
        Header: "Blocks",
        accessor: "stats.bpg",
      },
      {
        Header: "Turnovers",
        accessor: "stats.topg",
      },
      {
        Header: "FG%",
        accessor: "stats.fgpct",
      },
      {
        Header: "3FG%",
        accessor: "stats.threefgpct",
      },
      {
        Header: "FT%",
        accessor: "stats.ftpct",
      },
      {
        Header: "PER",
        accessor: "advstats.per",
      },
      {
        Header: "OWS",
        accessor: "advstats.ows",
      },
      {
        Header: "DWS",
        accessor: "advstats.dws",
      },
      {
        Header: "WS",
        accessor: "advstats.ws",
      },
      {
        Header: "OBPM",
        accessor: "advstats.obpm",
      },
      {
        Header: "DBPM",
        accessor: "advstats.dbpm",
      },
      {
        Header: "BPM",
        accessor: "advstats.bpm",
      },
      {
        Header: "VORP",
        accessor: "advstats.vorp",
      },
      {
        Header: "ORtg",
        accessor: "advstats.offrtg",
      },
      {
        Header: "DRtg",
        accessor: "advstats.defrtg",
      },
    ],
    []
  );

  const fetchData = useCallback(async () => {
    try {
      const response = await axios.get('http://localhost:8080/sixthman');
      setData(response.data);
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  }, []);

  useEffect(() => {
    // Fetch data from Go server
    fetchData();
  }, [fetchData]); 

  return (
    <div>
      <StatsTable columns={columns} data={data} />
    </div>
  );
};

