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
  Snackbar,
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import RefreshIcon from '@mui/icons-material/Refresh';
import { Link } from 'react-router-dom';
import useAuth from '../hooks/use-auth';
import axiosInstance from '../utils/axios-instance';

export type Poll = {
  id: number;
  name: string;
  description: string;
  image: string;
  season: string;
  userid: number;
};

const MyPollsPage: React.FC = () => {
  const [polls, setPolls] = useState<Poll[]>([]);
  const { auth } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    fetchPolls();
  }, []);

  const fetchPolls = () => {
    axiosInstance.get<Poll[]>(`/polls/users/get/${auth.id}`)
      .then(response => {
        setPolls(response.data);
      })
      .catch(error => {
        console.error('Error fetching polls:', error);
      });
  };

  const deletePoll = (id: number) => {
    axiosInstance.delete(`/polls/delete/${id}`)
      .then(() => {
        fetchPolls();
        setSuccessMessage('Poll deleted successfully!');
    })
      .catch(error => {
        console.error('Error deleting poll:', error);
        setErrorMessage('An error occurred. Please try again later.');
      });
  };

  const resetVotes = (id: number) => {
    axiosInstance.post("/polls/votes/reset", id)
    .then(() => {
      fetchPolls();
      setSuccessMessage('Poll votes resetted successfully!');
    })
    .catch(error => {
      console.error('Error reseting poll votes:', error);
      setErrorMessage('An error occurred. Please try again later.');
    });
  };

  const [successMessage, setSuccessMessage] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const handleCloseSnackbar = () => {
    setSuccessMessage('');
    setErrorMessage('');
  };

  return (
    <Container maxWidth="md">
      <Typography variant="h4" gutterBottom>
        Poll List
      </Typography>
      <Paper elevation={3} sx={{ padding: 2 }}>
        <List>
          {polls.map((poll) => (
            <ListItem key={poll.id} sx={{
                transition: 'box-shadow 0.3s ease-in-out',
                '&:hover': {
                  boxShadow: '0px 5px 10px rgba(0, 0, 0, 0.2)',
                },
              }}>
              <ListItemAvatar>
                <Avatar src={`../../${poll.image}`}/>
              </ListItemAvatar>
              <ListItemText primary={poll.name} secondary={poll.description} />
              <ListItemSecondaryAction>
                <Box>
                  <Link to={`/edit-poll/${poll.id}`}>
                    <IconButton edge="end" aria-label="edit">
                      <EditIcon />
                    </IconButton>
                  </Link>
                  <IconButton edge="end" aria-label="reset" onClick={() => resetVotes(poll.id)}>
                    <RefreshIcon />
                  </IconButton>
                  <IconButton edge="end" aria-label="delete" onClick={() => deletePoll(poll.id)}>
                    <DeleteIcon />
                  </IconButton>
                </Box>
              </ListItemSecondaryAction>
            </ListItem>
          ))}
        </List>
        <Box textAlign="center" mt={1}>
          <Button variant="contained" color="primary" onClick={() => navigate("/create-poll")}>
            Create New Poll
          </Button>
        </Box>
      </Paper>
      <Snackbar
        open={Boolean(successMessage)}
        autoHideDuration={3000}
        onClose={handleCloseSnackbar}
        message={successMessage}
        sx={{ bottom: 100 }}
      />

      <Snackbar
        open={Boolean(errorMessage)}
        autoHideDuration={3000}
        onClose={handleCloseSnackbar}
        message={errorMessage}
        sx={{ bottom: 100 }}
      />
    </Container>
  );
};

export default MyPollsPage;
