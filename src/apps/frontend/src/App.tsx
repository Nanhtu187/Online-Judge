import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import ProblemListPage from './pages/ProblemListPage';
import ProblemPage from './pages/ProblemPage';
import UpsertProblemPage from './pages/UpsertProblemPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<ProblemListPage />} />
        <Route path="/problems" element={<ProblemListPage />} />
        <Route path="/problems/new" element={<UpsertProblemPage />} />
        <Route path="/problems/:id" element={<ProblemPage />} />
        <Route path="/problems/:id/edit" element={<UpsertProblemPage />} />
      </Routes>
    </Router>
  );
}

export default App;
