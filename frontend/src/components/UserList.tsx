import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
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
import { Link } from 'react-router-dom';
import useAuth from '../hooks/use-auth';
import axiosInstance from '../utils/axios-instance';

type User = {
  id: number;
  username: string;
  email: string;
  password: string;
  refresh_token: string;
  profile_pic: string;
};

const UserListPage: React.FC = () => {
  const [users, setUsers] = useState<User[]>([]);
  const { auth, setAuth } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = () => {
    axiosInstance.get<User[]>('/users/get')
      .then(response => {
        setUsers(response.data);
      })
      .catch(error => {
        console.error('Error fetching users:', error);
      });
  };


  const deleteUser = (id: number) => {
    axiosInstance.delete(`/users/delete/${id}`)
      .then(() => {
        fetchUsers();
      })
      .catch(error => {
        console.error('Error deleting user:', error);
      });
  };

  const changeAdmin = async (id: number) => {
    try{
      const response = await axiosInstance.post('/update-admin', {id: id} );
      if (response.data) {
        setAuth(prev => {
          return {
              ...prev,
              roles: response.data.split(",")
          }
      });
        fetchUsers();
      }
    } catch (error) {
      console.error('Error changing admin:', error);
    }
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
                    <AdminIcon color={auth.roles.includes("admin") ? 'primary' : 'secondary'} />
                  </IconButton>
                  {user.id !== auth.id && (
                  <IconButton edge="end" aria-label="delete" onClick={() => deleteUser(user.id)}>
                    <DeleteIcon />
                  </IconButton>
                  )}
                </Box>
              </ListItemSecondaryAction>
            </ListItem>
          ))}
        </List>
        <Box textAlign="center" mt={1}>
            <Button variant="contained" color="primary" onClick={() => navigate("/admin/create-user")}>
            Create New User
            </Button>
        </Box>
      </Paper>
    </Container>
  );
};

export default UserListPage;
