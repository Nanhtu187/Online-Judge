import axios from 'axios';
import type { 
  Problem, 
  ProblemSummary, 
  SubmissionResultDetailResponse,
  TestCase,
  UpsertProblemRequest,
  UpsertProblemResponse,
  UpsertTestCasesRequest,
  UpsertTestCasesResponse,
  SubmissionType,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
  User,
  Submission
} from '../types';

const API_BASE_URL = 'http://localhost:8081';
const IAM_BASE_URL = 'http://localhost:8080';

const client = axios.create({
  baseURL: API_BASE_URL,
});

// Auth interceptor
client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

const iamClient = axios.create({
  baseURL: IAM_BASE_URL,
});

iamClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const login = async (data: LoginRequest): Promise<string> => {
  const response = await iamClient.post<LoginResponse>('/v1/iam/login', data);
  return response.data.token;
};

export const register = async (data: RegisterRequest): Promise<string> => {
  const response = await iamClient.post<RegisterResponse>('/v1/iam/register', data);
  return response.data.user_id;
};

export const getUserInfo = async (): Promise<User> => {
  const response = await iamClient.get<{ user: User }>('/v1/iam/user');
  return response.data.user;
};

export const listProblems = async (page = 1, pageSize = 20): Promise<ProblemSummary[]> => {
  const response = await client.get('/v1/problems', {
    params: { page, page_size: pageSize },
  });
  return response.data.problems || [];
};

export const getProblem = async (id: string): Promise<Problem> => {
  const response = await client.get(`/v1/problems/${id}`);
  return response.data.problem;
};

export const upsertProblem = async (data: UpsertProblemRequest): Promise<string> => {
  const response = await client.post<UpsertProblemResponse>('/v1/problems', data);
  return response.data.id;
};

export const listTestCases = async (problemId: string): Promise<TestCase[]> => {
  const response = await client.get(`/v1/problems/${problemId}/test-cases`);
  return response.data.test_cases || [];
};

export const upsertTestCases = async (problemId: string, testCases: TestCase[]): Promise<boolean> => {
  const request: UpsertTestCasesRequest = {
    problem_id: problemId,
    test_cases: testCases,
  };
  const response = await client.post<UpsertTestCasesResponse>(`/v1/problems/${problemId}/test-cases`, request);
  return response.data.success;
};

export const listSubmissions = async (params: { problem_id?: string, user_id?: string, page?: number, page_size?: number }): Promise<Submission[]> => {
  const response = await client.get('/v1/submissions', {
    params,
  });
  return response.data.submissions || [];
};

export const submitCode = async (problemId: string, codeContent: string, language: string, type: SubmissionType): Promise<string> => {
  const response = await client.post('/v1/submissions', {
    problem_id: problemId,
    code_content: codeContent,
    language,
    submission_type: type,
  });
  return response.data.submission_id;
};

export const getSubmissionResultDetail = async (submissionId: string): Promise<SubmissionResultDetailResponse> => {
  const response = await client.get(`/v1/submissions/${submissionId}/results/detail`);
  return response.data;
};
