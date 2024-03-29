import { Link, useNavigate } from 'react-router-dom';
import useAuth from '../hooks/use-auth';
import { useState } from 'react';
import { Avatar, IconButton, Menu, MenuItem } from '@mui/material';
import React from 'react';
import {useLogout} from '../hooks/use-logout';
import SearchBar from './SearchBar'
const Logo = require('../images/logo-no-background.png')


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
		navigate(`/edit-user/${auth?.id}`)
	}

	const myVotes = () => {
		handleClose();
		navigate(`/my-votes/${auth?.id}`)
	}

	const myPolls = () => {
		handleClose();
		navigate(`/my-polls/${auth?.id}`)
	}

	const userList = () => {
		handleClose();
		navigate(`/admin/users`)
	}

	const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
		setAnchorEl(event.currentTarget);
	};

	const handleClose = () => {
		setAnchorEl(null);
	};

	return (
		<header className="header">
			<div className='header-logo-styling'>
			 <img src={Logo} className="header-logo" />
			 <Link to="/">
 				 <h1 className="text-3xl" style={{ color: 'white' }}>HoopsVote</h1>
			</Link>
			</div>
			<div className="header-buttons">
				{auth?.accessToken ? (
					<div>
						<IconButton aria-controls="profile-menu" aria-haspopup="true" onClick={handleClick}>
							<Avatar src={`${window.location.origin}/${auth.user}-${auth.id}.jpg`} alt="default-user.jpg" />
						</IconButton>
						<Menu id="profile-menu" anchorEl={anchorEl} keepMounted open={Boolean(anchorEl)} onClose={handleClose}>
						<MenuItem onClick={editAcc}>Account</MenuItem>
						<MenuItem onClick={myVotes}>My votes</MenuItem>
						<MenuItem onClick={myPolls}>My polls</MenuItem>
						{ auth.roles.includes("admin") ?
							<MenuItem onClick={userList}>Users</MenuItem> : ''
						}
						<MenuItem onClick={signOut}>Sign out</MenuItem>
					</Menu>
				</div>
				) : (
					<React.Fragment>
						<div className="search-bar">
							<SearchBar></SearchBar>
						</div>
						<Link to="/login">Login</Link>
						<Link to="/signup">Sign up</Link>
					</React.Fragment>
					
				)}
			</div>
		</header>
	);
}

export default Header;
