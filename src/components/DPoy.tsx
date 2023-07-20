import { useEffect, useMemo, useState } from 'react';
import axios from 'axios';
import StatsTable from './StatsTable';

export default function DPoy () {
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
        Header: "Rebounds",
        accessor: "stats.rpg",
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
        Header: "DWS",
        accessor: "advstats.dws",
      },
      {
        Header: "DBPM",
        accessor: "advstats.dbpm",
      },
      {
        Header: "DRtg",
        accessor: "advstats.defrtg",
      },
    ],
    []
  );

  useEffect(() => {
    // Fetch data from Go server
   const fetchData = async () => { await axios.get('http://localhost:8080/dpoy')
      .then((response) => {
        setData(response.data);
      })
      .catch((error) => {
        console.error('Error fetching data:', error);
      });
 
    }
    fetchData();

}, []);

  return (
    <div>
      <StatsTable columns={columns} data={data} />
    </div>
  );
};

