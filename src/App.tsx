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
import QuizCreationPage from './components/QuizCreationPage';

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
          <Route path="/quiz/:pollId" element={<TableData endpoint='http://localhost:8080/quiz'/>} />
          <Route path="/results/:pollId" element={<Results />} />
          <Route path="/create-quiz" element={<QuizCreationPage />} />

          <Route element={<RequireAuth allowedRoles={["user", "admin"]} />}>
            <Route path="/edit-user/:userId" element={<AccountEditPage />} />
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