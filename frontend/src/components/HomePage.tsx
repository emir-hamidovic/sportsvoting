import { useNavigate } from "react-router-dom";
import Polls from "./Polls";
import { Box } from "@mui/material";
import { Button } from "react-bootstrap";
import useAuth from "../hooks/use-auth";
import '../css/HomePage.css';

const HomePage = () => {

	const createPollStyles = {
    background: 'lightblue', // Set your desired background color
    width: '20%', // Set your desired width
    padding: '10px', // Set any additional styling
    borderRadius: '8px', // Set border-radius for rounded corners // Center the box horizontally
    textAlign: 'center', // Center text within the box
	height: '15%',
	color: 'white',
	fontWeight: 'bold',
	fontSize: '16px',
	marginLeft: '2rem',
	marginTop: '4.5rem',
  };
	const navigate = useNavigate();
	const {auth} = useAuth();
	return (
		<div className="centered-polls">
			<Polls />
			{auth.user ? <Box sx={createPollStyles} mt={1}>
			<div className="createPollText">
				<p>Interested in making a new poll? Create a new poll and contribute to the community!</p>
				<Button variant="contained" color="primary" onClick={() => navigate("/create-poll")}>
				Create New Poll
				</Button>
				</div>
			</Box> : ''}
		</div>
	);
};

export default HomePage;