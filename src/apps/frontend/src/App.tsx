import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import ProblemListPage from './pages/ProblemListPage';
import ProblemPage from './pages/ProblemPage';
import UpsertProblemPage from './pages/UpsertProblemPage';
import SubmissionListPage from './pages/SubmissionListPage';
import PortfolioPage from './pages/PortfolioPage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import ForbiddenPage from './pages/ForbiddenPage';
import { AuthProvider } from './AuthContext';
import ProtectedRoute from './ProtectedRoute';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/forbidden" element={<ForbiddenPage />} />
          <Route path="/" element={<ProblemListPage />} />
          <Route path="/problems" element={<ProblemListPage />} />
          <Route path="/submissions" element={<SubmissionListPage />} />
          <Route 
            path="/portfolio" 
            element={
              <ProtectedRoute>
                <PortfolioPage />
              </ProtectedRoute>
            } 
          />
          
          <Route 
            path="/problems/new" 
            element={
              <ProtectedRoute permission="PROBLEM_CREATE">
                <UpsertProblemPage />
              </ProtectedRoute>
            } 
          />
          
          <Route path="/problems/:id" element={<ProblemPage />} />
          
          <Route 
            path="/problems/:id/edit" 
            element={
              <ProtectedRoute permission="PROBLEM_UPDATE">
                <UpsertProblemPage />
              </ProtectedRoute>
            } 
          />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
