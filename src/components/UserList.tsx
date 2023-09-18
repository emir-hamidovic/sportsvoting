import React, { useEffect, useState } from 'react';
import {
  Avatar,
  Box,
  Container,
  Button,
  IconButton,
  List,
  ListItem,
  ListItemAvatar,
  ListItemSecondaryAction,
  ListItemText,
  Paper,
  Typography,
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import AdminIcon from '@mui/icons-material/SupervisorAccount';
import axios from 'axios';
import { Link } from 'react-router-dom';

type User = {
  id: number;
  username: string;
  email: string;
  password: string;
  refresh_token: string;
  profile_pic: string;
  is_admin: boolean;
};

const UserListPage: React.FC = () => {
  const [users, setUsers] = useState<User[]>([]);

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = () => {
    axios.get<User[]>('http://localhost:8080/users/get')
      .then(response => {
        setUsers(response.data);
      })
      .catch(error => {
        console.error('Error fetching users:', error);
      });
  };


  const deleteUser = (id: number) => {
    axios.delete(`http://localhost:8080/users/delete/${id}`)
      .then(() => {
        fetchUsers();
      })
      .catch(error => {
        console.error('Error deleting user:', error);
      });
  };

  const changeAdmin = (id: number) => {
    axios.post('http://localhost:8080/api/update-admin', id)
      .then(() => {
        fetchUsers();
      })
      .catch(error => {
        console.error('Error changing admin:', error);
      });
  };


  return (
    <Container maxWidth="md">
      <Typography variant="h4" gutterBottom>
        User List
      </Typography>
      <Paper elevation={3} sx={{ padding: 2 }}>
        <List>
          {users.map((user) => (
            <ListItem key={user.id} sx={{
                transition: 'box-shadow 0.3s ease-in-out',
                '&:hover': {
                  boxShadow: '0px 5px 10px rgba(0, 0, 0, 0.2)',
                },
              }}>
              <ListItemAvatar>
                <Avatar src={user.profile_pic} />
              </ListItemAvatar>
              <ListItemText primary={user.username} secondary={user.email} />
              <ListItemSecondaryAction>
                <Box>
                  <Link to={`/admin/edit-user/${user.id}`}>
                    <IconButton edge="end" aria-label="edit">
                      <EditIcon />
                    </IconButton>
                  </Link>
                  <IconButton edge="end" aria-label="admin" onClick={() => changeAdmin(user.id)}>
                    <AdminIcon color={user.is_admin ? 'primary' : 'inherit'} />
                  </IconButton>
                  <IconButton edge="end" aria-label="delete" onClick={() => deleteUser(user.id)}>
                    <DeleteIcon />
                  </IconButton>
                </Box>
              </ListItemSecondaryAction>
            </ListItem>
          ))}
        </List>
        <Box textAlign="center" mt={1}>
            <Button variant="contained" color="primary">
            Create New User
            </Button>
        </Box>
      </Paper>
    </Container>
  );
};

export default UserListPage;