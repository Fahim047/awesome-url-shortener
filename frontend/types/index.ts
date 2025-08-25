export interface ShortenResponse {
  short_key: string;
  short_url: string;
}

export interface AnalyticsResponse {
  short_key: string;
  click_count: number;
  long_url: string;
  created_at: string;
  expire_at?: string;
  custom_alias?: string;
}
