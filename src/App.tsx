import './App.css';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import HomePage from './components/HomePage';
import TableData from './components/TableData';
import Header from './components/Header';

function App() {
  return (
    <BrowserRouter>
      <Header />
      <Routes>
        <Route path="/mvp/:pollId" element={<TableData endpoint='http://localhost:8080/mvp'/>} />
        <Route path="/sixthman/:pollId" element={<TableData endpoint='http://localhost:8080/sixthman' />} />
        <Route path="/roy/:pollId" element={<TableData endpoint='http://localhost:8080/roy' />} />
        <Route path="/dpoy/:pollId" element={<TableData endpoint='http://localhost:8080/dpoy' />} />
        <Route path="/" element={<HomePage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;