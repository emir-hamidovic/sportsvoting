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
        <Route path="/mvp" element={<TableData endpoint='http://localhost:8080/mvp'/>} />
        <Route path="/sixthman" element={<TableData endpoint='http://localhost:8080/sixthman' />} />
        <Route path="/roy" element={<TableData endpoint='http://localhost:8080/roy' />} />
        <Route path="/dpoy" element={<TableData endpoint='http://localhost:8080/dpoy' />} />
        <Route path="/" element={<HomePage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;