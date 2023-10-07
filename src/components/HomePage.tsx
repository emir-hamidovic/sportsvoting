import { useNavigate } from "react-router-dom";
import Polls from "./Polls";
import { Box } from "@mui/material";
import { Button } from "react-bootstrap";
import useAuth from "../hooks/use-auth";

const HomePage = () => {
  const navigate = useNavigate();
  const {auth} = useAuth();
  return (
    <div className="App">
        <Polls />
        {auth.user ? <Box textAlign="center" mt={1}>
            <Button variant="contained" color="primary" onClick={() => navigate("/create-quiz")}>
            Create New Quiz
            </Button>
        </Box> : ''}
    </div>
  );
};

export default HomePage;