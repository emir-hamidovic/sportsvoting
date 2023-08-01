import {Avatar, Box, Card, Checkbox, Stack, Table, TableBody, TableCell, TableHead, TablePagination, TableRow, Typography, Button} from '@mui/material';
import { FlattenedAPIResponse } from '../utils/api-response';
import axios from 'axios';
import { useParams } from 'react-router-dom';


interface CustomersTableProps {
  count: number;
  items: FlattenedAPIResponse[];
  onDeselectOne?: (customer: string) => void;
  onPageChange?: (event: React.MouseEvent<HTMLButtonElement> | null, page: React.SetStateAction<number>) => void;
  onRowsPerPageChange?: (event: React.ChangeEvent<HTMLInputElement>) => void;
  onSelectOne?: (customer: string) => void;
  page: number;
  rowsPerPage: number;
  selected: string[];
  columns: string[]
}

const getInitials = (name = '') => name
  .replace(/\s+/, ' ')
  .split(' ')
  .slice(0, 2)
  .map((v) => v && v[0].toUpperCase())
  .join('');

const columnNameMapping: { [key: string]: string } = {
    g: 'Games',
    mpg: 'Minutes',
    ppg: 'Points',
    apg: 'Assists',
    rpg: 'Rebounds',
    spg: 'Steals',
    bpg: 'Blocks',
    topg: 'Turnovers',
    fgpct: 'FG%',
    threefgpct: '3FG%',
    ftpct: 'FT%',
    per: 'PER',
    ows: 'OWS',
    dws: 'DWS',
    ws: 'WS',
    obpm: 'OBPM',
    dbpm: 'DBPM',
    bpm: 'BPM',
    vorp: 'VORP',
    offrtg: 'ORtg',
    defrtg: 'DRtg',
};

export const CustomersTable = (props: CustomersTableProps) => {
  const {
    count = 0,
    items = [],
    onDeselectOne,
    onPageChange = () => {},
    onRowsPerPageChange,
    onSelectOne,
    page = 0,
    rowsPerPage = 0,
    selected = [],
    columns
  } = props;

  const tableFields = columns.slice(2);
  const { pollId } = useParams();
  const id = pollId ? parseInt(pollId, 10) : undefined;

  const handleVote = () => {
    const selectedCustomerIds = selected.length === 1 ? selected[0] : '';

    const voteEndpoint = 'http://localhost:8080/playervotes/';
    const payload = { playerid: selectedCustomerIds, pollid: Number(id) };
    axios.post(voteEndpoint, payload, {
      headers: {
        'Content-Type': 'application/json',
      },
    })
      .then((response) => {
        console.log(response.data);
      })
      .catch((error) => {
        console.error(error);
      });
  };

  return (
    <Card>
        <Box sx={{ minWidth: 800 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell padding="checkbox"></TableCell>
                {columns.map((column) => (
                <TableCell key={column}>
                    {columnNameMapping[column]}
                </TableCell>
                ))}
              </TableRow>
            </TableHead>
            <TableBody>
              {items.map((customer: FlattenedAPIResponse) => {
                const isSelected = selected.includes(customer.playerid);

                return (
                  <TableRow hover key={customer.playerid} selected={isSelected}>
                    <TableCell padding="checkbox">
                      <Checkbox checked={isSelected}
                        onChange={(event) => {
                          if (event.target.checked) {
                            onSelectOne?.(customer.playerid);
                          } else {
                            onDeselectOne?.(customer.playerid);
                          }
                        }} />
                    </TableCell>
                    <TableCell>
                      <Stack alignItems="center" direction="row" spacing={2}>
                        <Avatar> {getInitials(customer.name)} </Avatar>
                        <Typography variant="subtitle2"> {customer.name} </Typography>
                      </Stack>
                    </TableCell>
                    <TableCell></TableCell>
                    {tableFields.map((column) => (
                      <TableCell key={column}>
                        {customer[column]}
                      </TableCell>
                    ))}
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </Box>
      <TablePagination
        component="div"
        count={count}
        onPageChange={onPageChange}
        onRowsPerPageChange={onRowsPerPageChange}
        page={page}
        rowsPerPage={rowsPerPage}
        rowsPerPageOptions={[5, 10, 25]}
      />
      <Box sx={{ display: 'flex', justifyContent: 'flex-end', marginTop: 2 }}>
        <Button
          variant="contained"
          color="primary"
          disabled={selected.length !== 1}
          onClick={handleVote}
        >
          Vote
        </Button>
      </Box>
    </Card>
  );
};