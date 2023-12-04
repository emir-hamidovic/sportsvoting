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
import Unauthorized from './components/Unauthorized';
import AdminUserCreationForm from './components/AdminUserCreationForm';
import PollCreationPage from './components/PollCreationPage';
import MyVotesPage from './components/MyVotesPage';
import MyPollsPage from './components/MyPollsPage';
import EditPollPage from './components/EditPollPage';

function App() {
  return (
    <BrowserRouter>
    <AuthProvider>
      <Header />
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Register />} />
        <Route path="/unauthorized" element={<Unauthorized />} />

        <Route element={<PersistLogin />}>
          <Route path="/results/:pollId" element={<Results />} />
          <Route element={<RequireAuth allowedRoles={["user", "admin"]} />}>
            <Route path="/poll/:pollId" element={<TableData endpoint='/polls/players/get'/>} />
          </Route>

          <Route element={<RequireAuth allowedRoles={["user", "admin"]} />}>
            <Route path="/create-poll" element={<PollCreationPage />} />
          </Route>

          <Route element={<RequireAuth allowedRoles={["user", "admin"]} />}>
            <Route path="/edit-user/:userId" element={<AccountEditPage />} />
          </Route>

          <Route element={<RequireAuth allowedRoles={["user", "admin"]} />}>
            <Route path="/edit-poll/:pollId" element={<EditPollPage />} />
          </Route>

          <Route element={<RequireAuth allowedRoles={["user", "admin"]} />}>
            <Route path="/my-votes/:userId" element={<MyVotesPage />} />
          </Route>

          <Route element={<RequireAuth allowedRoles={["user", "admin"]} />}>
            <Route path="/my-polls/:userId" element={<MyPollsPage />} />
          </Route>
          
          <Route element={<RequireAuth allowedRoles={["admin"]}/>}>
            <Route path="/admin/edit-user/:userId" element={<AccountEditPage />} />
          </Route>

          <Route element={<RequireAuth allowedRoles={["admin"]}/>}>
            <Route path="/admin/users" element={<UserListPage />} />
          </Route>

          <Route element={<RequireAuth allowedRoles={["admin"]}/>}>
            <Route path="/admin/create-user" element={<AdminUserCreationForm />} />
          </Route>
          <Route path="/" element={<HomePage />} />
        </Route>
      </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;