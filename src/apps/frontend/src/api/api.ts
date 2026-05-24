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
  SubmissionType
} from '../types';

const API_BASE_URL = 'http://localhost:8081';

const client = axios.create({
  baseURL: API_BASE_URL,
});

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
