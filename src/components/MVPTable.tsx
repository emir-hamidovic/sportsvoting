import { useEffect, useMemo, useState } from 'react';
import axios from 'axios';
import { useTable, ColumnInterface } from 'react-table';

export interface Column extends ColumnInterface {
  Header: string;
  accessor: string;
}

export default function MVPTable () {
  const [data, setData] = useState([]);

  const columns: Column[] = useMemo(
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

  useEffect(() => {
    // Fetch data from Go server
   const fetchData = async () => { await axios.get('http://localhost:8080/mvp')
      .then((response) => {
        setData(response.data);
      })
      .catch((error) => {
        console.error('Error fetching data:', error);
      });
 
    }
    fetchData();

}, []);

  const {
    getTableProps, // table props from react-table
    getTableBodyProps, // table body props from react-table
    headerGroups, // headerGroups, if your table has groupings
    rows, // rows for the table based on the data passed
    prepareRow // Prepare the row (this function needs to be called for each row before getting the row props)
  } = useTable({
    columns ,
    data
  });


  return (
    <div>
      <div className="table-container">

      <table {...getTableProps()}>
      <thead>
      {headerGroups.map(headerGroup => (
          <tr {...headerGroup.getHeaderGroupProps()}>
            {headerGroup.headers.map(column => (
              <th {...column.getHeaderProps()}>{column.render("Header")}</th>
            ))}
          </tr>
        ))}
      </thead>
      <tbody {...getTableBodyProps()}>
        {rows.map((row, i) => {
          prepareRow(row);
          return (
            <tr {...row.getRowProps()}>
              {row.cells.map(cell => {
                return <td {...cell.getCellProps()}>{cell.render("Cell")}</td>;
              })}
            </tr>
          );
        })}
      </tbody>

      </table>

      </div>
    </div>
  );
};

