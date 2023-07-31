import {Avatar, Box, Card, Checkbox, Stack, Table, TableBody, TableCell, TableHead, TablePagination, TableRow, Typography} from '@mui/material';
import { FlattenedAPIResponse } from '../utils/api-response';

interface CustomersTableProps {
  count: number;
  items: FlattenedAPIResponse[];
  onDeselectAll?: () => void;
  onDeselectOne?: (customer: string) => void;
  onPageChange?: (event: React.MouseEvent<HTMLButtonElement> | null, page: React.SetStateAction<number>) => void;
  onRowsPerPageChange?: (event: React.ChangeEvent<HTMLInputElement>) => void;
  onSelectAll?: () => void;
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
    onDeselectAll,
    onDeselectOne,
    onPageChange = () => {},
    onRowsPerPageChange,
    onSelectAll,
    onSelectOne,
    page = 0,
    rowsPerPage = 0,
    selected = [],
    columns
  } = props;

  const selectedSome = (selected.length > 0) && (selected.length < items.length);
  const selectedAll = (items.length > 0) && (selected.length === items.length);
  const tableFields = columns.slice(2);

  return (
    <Card>
        <Box sx={{ minWidth: 800 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell padding="checkbox">
                  <Checkbox
                    checked={selectedAll}
                    indeterminate={selectedSome}
                    onChange={(event) => {
                      if (event.target.checked) {
                        onSelectAll?.();
                      } else {
                        onDeselectAll?.();
                      }
                    }}
                  />
                </TableCell>
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
    </Card>
  );
};