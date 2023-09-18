import './App.css';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import HomePage from './components/HomePage';
import TableData from './components/TableData';
import Header from './components/Header';
import Results from './components/Results';
import Register from './components/Register';
import Login from './components/Login';
import AccountEditPage from './components/AccountEditPage';
import UserListPage from './components/UserList';
import { AuthProvider } from './context/AuthProvider';
import RequireAuth from './components/RequireAuth';
import PersistLogin from './components/PersistLogin';

function App() {
  return (
    <BrowserRouter>
    <AuthProvider>
      <Header />
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Register />} />

        <Route element={<PersistLogin />}>
          <Route path="/mvp/:pollId" element={<TableData endpoint='http://localhost:8080/mvp'/>} />
          <Route path="/sixthman/:pollId" element={<TableData endpoint='http://localhost:8080/sixthman' />} />
          <Route path="/roy/:pollId" element={<TableData endpoint='http://localhost:8080/roy' />} />
          <Route path="/dpoy/:pollId" element={<TableData endpoint='http://localhost:8080/dpoy' />} />
          <Route path="/results/:pollId" element={<Results />} />
          
          <Route element={<RequireAuth />}>
            <Route path="/edit-user/:userId" element={<AccountEditPage />} />
          </Route>
          
          <Route element={<RequireAuth />}>
            <Route path="/admin/edit-user/:userId" element={<AccountEditPage />} />
          </Route>

          <Route element={<RequireAuth />}>
            <Route path="/admin/users" element={<UserListPage />} />
          </Route>   
          <Route path="/" element={<HomePage />} />
        </Route>
      </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;