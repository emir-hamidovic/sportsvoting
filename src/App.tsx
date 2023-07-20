import './App.css';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import HomePage from './components/HomePage';
import MVPTable from './components/MVPTable';
import Header from './components/Header';
import "./styles/table.css";

function App() {
  return (
    <BrowserRouter>
      <Header />
      <Routes>
        <Route path="/mvp" element={<MVPTable />} />
        <Route path="/" element={<HomePage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;