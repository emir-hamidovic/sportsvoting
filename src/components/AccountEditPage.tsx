import React, { useEffect, useState } from 'react';
import {
  Avatar,
  Box,
  Button,
  Container,
  Divider,
  Grid,
  Paper,
  Snackbar,
  TextField,
  Typography,
} from '@mui/material';
import { useParams } from 'react-router-dom';
import axios from 'axios';
import useAuth from '../hooks/use-auth';

const USER_REGEX = /^[A-z][A-z0-9-_]{3,23}$/;
const PWD_REGEX = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%]).{8,24}$/;
const EMAIL_REGEX = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

const AccountEditPage: React.FC = () => {
  const { auth } = useAuth();
  const { userId } = useParams();
  const [userInfo, setUserInfo] = useState({
    email: '',
    username: '',
    profile_pic: '',
    oldPassword: '',
    newPassword: '',
  });

  const [initialEmail, setInitialEmail] = useState('');

  useEffect(() => {
    axios.get(`http://localhost:8080/api/get-user/${userId}`)
      .then(response => {
        const userData = response.data;
        setUserInfo({
          email: userData.email,
          username: userData.username,
          profile_pic: userData.profile_pic,
          oldPassword: '', // Leave empty for security reasons
          newPassword: '', // Leave empty for security reasons
        });
        setInitialEmail(userData.email);
      })
      .catch(error => {
        console.error(error);
      });
  }, [userId]);
  const [successMessage, setSuccessMessage] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const handleCloseSnackbar = () => {
    setSuccessMessage('');
    setErrorMessage('');
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const v1 = USER_REGEX.test(userInfo.username);
    const v2 = PWD_REGEX.test(userInfo.newPassword);
    const v3 = EMAIL_REGEX.test(userInfo.email);

    if ((!v1 && userInfo.username !== "") || (!v2 && userInfo.newPassword !== "") || (!v3 && userInfo.email !== "")) {
        setErrorMessage("Invalid Entry");
        return;
    }

    try {
      if(auth.user !== userInfo.username) {
        await axios.post('http://localhost:8080/api/update-username', {olduser: auth.user, username: userInfo.username});
        auth.user = userInfo.username;
      }

      if (initialEmail !== userInfo.email) {
        await axios.post('http://localhost:8080/api/update-email', {email: userInfo.email, username: auth.user});
      }

      if (userInfo.newPassword !== "") {
        await axios.post('http://localhost:8080/api/update-password', {
            oldPassword: userInfo.oldPassword,
            newPassword: userInfo.newPassword,
            username: auth.user
        });
      }

      setSuccessMessage('Account information updated successfully!');
    } catch (error) {
      setErrorMessage('An error occurred. Please try again later.');
    }
    
    setUserInfo({ ...userInfo, oldPassword: '', newPassword: '' });
  };

  return (
    <Container maxWidth="lg">
      <Paper elevation={3} sx={{ padding: 3 }}>
        <Grid container spacing={3}>
          <Grid item md={4} sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
            <Avatar src={userInfo.profile_pic} sx={{ width: 120, height: 120, marginBottom: 2 }}/>
            <Button variant="outlined" fullWidth>
              Upload New Picture
            </Button>
          </Grid>
          <Grid item md={8}>
            <Typography variant="h4" gutterBottom>
              Account Information
            </Typography>
            <Divider />

            <form onSubmit={handleSubmit}>
              <Box>
                <TextField
                  fullWidth
                  label="Email"
                  variant="outlined"
                  value={userInfo.email}
                  onChange={(e) => setUserInfo({ ...userInfo, email: e.target.value })}
                  InputProps={{ style: { marginTop: '25px' } }}
                />
              </Box>
              <Box>
                <TextField
                  fullWidth
                  label="Username"
                  variant="outlined"
                  value={userInfo.username}
                  onChange={(e) => setUserInfo({ ...userInfo, username: e.target.value })}
                  InputProps={{ style: { marginTop: '25px' } }}
                />
              </Box>

              <Divider sx={{ mt: 3, mb: 2 }} />
              <Typography variant="h5" gutterBottom>
                Change Password
              </Typography>
              <Box>
                <TextField
                  fullWidth
                  type="password"
                  label="Old Password"
                  variant="outlined"
                  value={userInfo.oldPassword}
                  onChange={(e) => setUserInfo({ ...userInfo, oldPassword: e.target.value })}
                  InputProps={{ style: { marginTop: '25px' } }}
                />
              </Box>
              <Box>
                <TextField
                  fullWidth
                  type="password"
                  label="New Password"
                  variant="outlined"
                  value={userInfo.newPassword}
                  onChange={(e) => setUserInfo({ ...userInfo, newPassword: e.target.value })}
                  InputProps={{ style: { marginTop: '25px' } }}
                />
              </Box>
              <Box mt={3}>
                <Button type="submit" variant="contained" color="primary">
                  Save Changes
                </Button>
              </Box>
            </form>
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
          </Grid>
        </Grid>
      </Paper>
    </Container>
  );
};

export default AccountEditPage;
