import axios, { type AxiosInstance, type AxiosError } from 'axios';
import type {
  Message,
  Contact,
  Template,
  SendMessageRequest,
  PaginationParams,
  PaginatedResponse,
  APIError,
} from '../types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';
const API_KEY = import.meta.env.VITE_API_KEY || '';

class APIClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${API_KEY}`,
      },
    });

    // Response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError<APIError>) => {
        if (error.response?.data) {
          throw error.response.data;
        }
        throw {
          code: 'NETWORK_ERROR',
          message: error.message || 'Network error occurred',
        };
      }
    );
  }

  // Messages
  async sendMessage(request: SendMessageRequest): Promise<Message> {
    const { data } = await this.client.post<Message>('/messages', request);
    return data;
  }

  async getMessage(id: string): Promise<Message> {
    const { data } = await this.client.get<Message>(`/messages/${id}`);
    return data;
  }

  async listMessages(params?: PaginationParams & { phone?: string; status?: string }): Promise<PaginatedResponse<Message>> {
    const { data } = await this.client.get<PaginatedResponse<Message>>('/messages', { params });
    return data;
  }

  async searchMessages(query: string, params?: PaginationParams): Promise<PaginatedResponse<Message>> {
    const { data } = await this.client.get<PaginatedResponse<Message>>('/messages/search', {
      params: { q: query, ...params },
    });
    return data;
  }

  // Contacts
  async getContact(id: string): Promise<Contact> {
    const { data } = await this.client.get<Contact>(`/contacts/${id}`);
    return data;
  }

  async listContacts(params?: PaginationParams): Promise<PaginatedResponse<Contact>> {
    const { data } = await this.client.get<PaginatedResponse<Contact>>('/contacts', { params });
    return data;
  }

  async searchContacts(query: string, params?: PaginationParams): Promise<PaginatedResponse<Contact>> {
    const { data } = await this.client.get<PaginatedResponse<Contact>>('/contacts/search', {
      params: { q: query, ...params },
    });
    return data;
  }

  async updateContact(id: string, updates: Partial<Contact>): Promise<Contact> {
    const { data } = await this.client.patch<Contact>(`/contacts/${id}`, updates);
    return data;
  }

  // Templates
  async createTemplate(template: Omit<Template, 'id' | 'created_at' | 'updated_at'>): Promise<Template> {
    const { data} = await this.client.post<Template>('/templates', template);
    return data;
  }

  async getTemplate(id: string): Promise<Template> {
    const { data } = await this.client.get<Template>(`/templates/${id}`);
    return data;
  }

  async listTemplates(params?: PaginationParams): Promise<PaginatedResponse<Template>> {
    const { data } = await this.client.get<PaginatedResponse<Template>>('/templates', { params });
    return data;
  }

  async updateTemplate(id: string, updates: Partial<Template>): Promise<Template> {
    const { data } = await this.client.patch<Template>(`/templates/${id}`, updates);
    return data;
  }

  async deleteTemplate(id: string): Promise<void> {
    await this.client.delete(`/templates/${id}`);
  }
}

export const apiClient = new APIClient();
