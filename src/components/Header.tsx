import { Link, useNavigate } from 'react-router-dom';
import useAuth from '../hooks/use-auth';
import { useState } from 'react';
import { Avatar, IconButton, Menu, MenuItem } from '@mui/material';
import React from 'react';
import {useLogout} from '../hooks/use-logout';

function Header() {
  const { auth } = useAuth();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const navigate = useNavigate();
  const logout = useLogout();

  const signOut = async () => {
    await logout();
    handleClose();
    navigate('/login');
  }

  const editAcc = () => {
    handleClose();
    navigate('/account')
  }

  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

auth.accessToken = "AAA";
  console.log(auth);
  return (
    <header className="header">
      <Link to="/"> <h1 className="underline text-3xl">Sport Voting</h1></Link>
      <div className="header-buttons">
        {auth?.accessToken ? (
          <div>
            <IconButton aria-controls="profile-menu" aria-haspopup="true" onClick={handleClick}>
              <Avatar src="https://pbs.twimg.com/media/Fj26barWYAMeoZp?format=jpg&name=4096x4096" alt={auth.user} />
            </IconButton>
            <Menu id="profile-menu" anchorEl={anchorEl} keepMounted open={Boolean(anchorEl)} onClose={handleClose}>
            <MenuItem onClick={editAcc}>Account</MenuItem>
            <MenuItem onClick={signOut}>Sign out</MenuItem>
          </Menu>
        </div>
        ) : (
          <React.Fragment>
            <Link to="/login">Login</Link>
            <Link to="/signup">Sign up</Link>
          </React.Fragment>
          
        )}
      </div>
    </header>
  );
}

export default Header;
