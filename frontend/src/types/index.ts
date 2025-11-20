export interface Message {
  id: string;
  whatsapp_message_id: string;
  from_number: string;
  to_number: string;
  direction: 'inbound' | 'outbound';
  message_type: 'text' | 'image' | 'document' | 'audio' | 'video' | 'template';
  content: string;
  media_url?: string;
  media_mime_type?: string;
  status: 'sent' | 'delivered' | 'read' | 'failed' | 'received';
  timestamp: string;
  created_at: string;
  updated_at: string;
}

export interface Contact {
  id: string;
  phone_number: string;
  name: string;
  last_message_at?: string;
  message_count: number;
  unread_count: number;
  created_at: string;
  updated_at: string;
}

export interface Template {
  id: string;
  name: string;
  language: string;
  category: string;
  content: string;
  status: 'approved' | 'pending' | 'rejected';
  created_at: string;
  updated_at: string;
}

export interface SendMessageRequest {
  to: string;
  message_type: 'text' | 'image' | 'document' | 'audio' | 'video' | 'template';
  content?: string;
  media_url?: string;
  template_name?: string;
  template_language?: string;
  template_params?: string[];
}

export interface PaginationParams {
  page?: number;
  limit?: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

export interface APIError {
  code: string;
  message: string;
  details?: Record<string, unknown>;
}
