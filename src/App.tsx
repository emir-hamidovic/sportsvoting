import './App.css';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import HomePage from './components/HomePage';
import MVP from './components/MVP';
import SixMan from './components/SixMan';
import Roy from './components/Roy';
import DPoy from './components/DPoy';
import Header from './components/Header';
import "./styles/table.css";

function App() {
  return (
    <BrowserRouter>
      <Header />
      <Routes>
        <Route path="/mvp" element={<MVP />} />
        <Route path="/sixthman" element={<SixMan />} />
        <Route path="/roy" element={<Roy />} />
        <Route path="/dpoy" element={<DPoy />} />
        <Route path="/" element={<HomePage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;