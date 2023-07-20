import { useTable, ColumnInterface, useSortBy } from 'react-table';

export interface Column extends ColumnInterface {
    Header: string;
    accessor: string;
}

const emojiStyle = {
    height: '1em', // Adjust the height to your desired size
    width: '1em', // Adjust the width to your desired size
    marginRight: '0.5em',
};

const UpArrowEmoji = () => (
    <img
      draggable="false"
      role="img"
      className="emoji"
      alt="🔼"
      src="https://s.w.org/images/core/emoji/14.0.0/svg/1f53c.svg"
      style={emojiStyle}
    />
);
  
const DownArrowEmoji = () => (
    <img
      draggable="false"
      role="img"
      className="emoji"
      alt="🔽"
      src="https://s.w.org/images/core/emoji/14.0.0/svg/1f53d.svg"
      style={emojiStyle}
    />
);

const StatsTable = ({columns, data} : {columns : Column[], data: any}) => {
    const {
        getTableProps, // table props from react-table
        getTableBodyProps, // table body props from react-table
        headerGroups, // headerGroups, if your table has groupings
        rows, // rows for the table based on the data passed
        prepareRow // Prepare the row (this function needs to be called for each row before getting the row props)
      } = useTable({
        columns,
        data
      }, useSortBy);

      return (
        <div>
          <div className="table-container">
          <table {...getTableProps()}>
          <thead>
          {headerGroups.map((headerGroup) => (
                    <tr {...headerGroup.getHeaderGroupProps()}>
                        {headerGroup.headers.map((column: any) => (
                            <th {...column.getHeaderProps(column.getSortByToggleProps())}>{column.render('Header')}
                            <span>
                                {column.isSorted ? (column.isSortedDesc ? <DownArrowEmoji /> : <UpArrowEmoji />) : ''}
                            </span>
                            </th>
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

export default StatsTable;
