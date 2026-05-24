export const SubmissionStatus = {
  UNSPECIFIED: 'SUBMISSION_STATUS_UNSPECIFIED',
  PENDING: 'SUBMISSION_STATUS_PENDING',
  RUNNING: 'SUBMISSION_STATUS_RUNNING',
  ACCEPTED: 'SUBMISSION_STATUS_ACCEPTED',
  FAILED: 'SUBMISSION_STATUS_FAILED',
  WRONG_ANSWER: 'SUBMISSION_STATUS_WRONG_ANSWER',
  RUNTIME_ERROR: 'SUBMISSION_STATUS_RUNTIME_ERROR',
  LIMIT_EXCEEDED: 'SUBMISSION_STATUS_LIMIT_EXCEEDED',
  COMPILE_ERROR: 'SUBMISSION_STATUS_COMPILE_ERROR',
} as const;

export type SubmissionStatus = typeof SubmissionStatus[keyof typeof SubmissionStatus];

export const SubmissionType = {
  UNSPECIFIED: 'SUBMISSION_TYPE_UNSPECIFIED',
  TEST: 'SUBMISSION_TYPE_TEST',
  OFFICIAL: 'SUBMISSION_TYPE_OFFICIAL',
} as const;

export type SubmissionType = typeof SubmissionType[keyof typeof SubmissionType];

export interface ProblemSummary {
  id: string;
  title: string;
}

export interface ListProblemsResponse {
  problems: ProblemSummary[];
}

export interface Problem {
  id: string;
  title: string;
  content: string;
  input_format: string;
  output_format: string;
  test_cases?: TestCase[];
  time_limit: number;
  memory_limit: number;
}

export interface TestCaseResult {
  test_case_id: string;
  status: SubmissionStatus;
  actual_output: string;
}

export interface TestCase {
  id?: string;
  input: string;
  expected_output: string;
  is_sample?: boolean;
}

export interface UpsertProblemRequest {
  id?: string;
  title: string;
  content: string;
  input_format: string;
  output_format: string;
  time_limit: number;
  memory_limit: number;
}

export interface UpsertProblemResponse {
  id: string;
}

export interface UpsertTestCasesRequest {
  problem_id: string;
  test_cases: TestCase[];
}

export interface UpsertTestCasesResponse {
  success: boolean;
}

export interface SubmissionResultDetailResponse {
  test_case_results: TestCaseResult[];
  status: SubmissionStatus;
}
