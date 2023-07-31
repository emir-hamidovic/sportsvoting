import { Link } from 'react-router-dom';

function Header() {
  return (
    <header className="header">
      <Link to="/"> <h1 className="underline text-3xl">Sport Voting</h1></Link>
      <div className="header-buttons">
        <Link to="/login">Login</Link>
        <Link to="/signup">Sign up</Link>
      </div>
    </header>
  );
}

export default Header;
